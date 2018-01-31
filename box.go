package packr

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"compress/gzip"

	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
)

var (
	ErrResOutsideBox = errors.New("Can't find a resource outside the box")
)

// NewBox returns a Box that can be used to
// retrieve files from either disk or the embedded
// binary.
func NewBox(path string) Box {
	var cd string
	if !filepath.IsAbs(path) {
		_, filename, _, _ := runtime.Caller(1)
		cd = filepath.Dir(filename)
	}

	// this little hack courtesy of the `-cover` flag!!
	cov := filepath.Join("_test", "_obj_test")
	cd = strings.Replace(cd, string(filepath.Separator)+cov, "", 1)
	if !filepath.IsAbs(cd) && cd != "" {
		cd = filepath.Join(envy.GoPath(), "src", cd)
	}

	return Box{
		Path:       path,
		callingDir: cd,
	}
}

// Box represent a folder on a disk you want to
// have access to in the built Go binary.
type Box struct {
	Path        string
	callingDir  string
	data        map[string][]byte
	directories map[string]bool
}

// String of the file asked for or an empty string.
func (b Box) String(name string) string {
	return string(b.Bytes(name))
}

// MustString returns either the string of the requested
// file or an error if it can not be found.
func (b Box) MustString(name string) (string, error) {
	bb, err := b.MustBytes(name)
	return string(bb), err
}

// Bytes of the file asked for or an empty byte slice.
func (b Box) Bytes(name string) []byte {
	bb, _ := b.MustBytes(name)
	return bb
}

// MustBytes returns either the byte slice of the requested
// file or an error if it can not be found.
func (b Box) MustBytes(name string) ([]byte, error) {
	f, err := b.find(name)
	if err == nil {
		bb := &bytes.Buffer{}
		bb.ReadFrom(f)
		return bb.Bytes(), err
	}
	return nil, err
}

// Has returns true if the resource exists in the box
func (b Box) Has(name string) bool {
	_, err := b.find(name)
	if err != nil {
		return false
	}
	return true
}

func (b Box) decompress(bb []byte) []byte {
	reader, err := gzip.NewReader(bytes.NewReader(bb))
	if err != nil {
		return bb
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return bb
	}
	return data
}

func (b Box) find(name string) (File, error) {
	if b.directories == nil {
		b.indexDirectories()
	}

	cleanName := filepath.ToSlash(filepath.Clean(name))
	// Ensure name is not outside the box
	if strings.HasPrefix(cleanName, "../") {
		return nil, ErrResOutsideBox
	}
	// Absolute name is considered as relative to the box root
	cleanName = strings.TrimPrefix(cleanName, "/")

	// Try to get the resource from the box
	if _, ok := data[b.Path]; ok {
		if bb, ok := data[b.Path][cleanName]; ok {
			bb = b.decompress(bb)
			return newVirtualFile(cleanName, bb), nil
		}
		if filepath.Ext(cleanName) != "" {
			// The Handler created by http.FileSystem checks for those errors and
			// returns http.StatusNotFound instead of http.StatusInternalServerError.
			return nil, os.ErrNotExist
		}
		if _, ok := b.directories[cleanName]; ok {
			return newVirtualDir(cleanName), nil
		}
		return nil, os.ErrNotExist
	}

	// Not found in the box virtual fs, try to get it from the file system
	cleanName = filepath.FromSlash(cleanName)
	p := filepath.Join(b.callingDir, b.Path, cleanName)
	if bb, err := ioutil.ReadFile(p); err == nil {
		return newVirtualFile(cleanName, bb), nil
	}
	return nil, os.ErrNotExist
}

type WalkFunc func(string, File) error

func (b Box) Walk(wf WalkFunc) error {
	if data[b.Path] == nil {
		base, err := filepath.EvalSymlinks(filepath.Join(b.callingDir, b.Path))
		if err != nil {
			return errors.WithStack(err)
		}
		return filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			shortPath := strings.TrimPrefix(path, base)
			if info == nil || info.IsDir() {
				return nil
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			return wf(shortPath, physicalFile{f})
		})
	}
	for n := range data[b.Path] {
		f, err := b.find(n)
		if err != nil {
			return err
		}
		err = wf(n, f)
		if err != nil {
			return err
		}
	}
	return nil
}

// Open returns a File using the http.File interface
func (b Box) Open(name string) (http.File, error) {
	return b.find(name)
}

// List shows "What's in the box?"
func (b Box) List() []string {
	var keys []string

	if b.data == nil {
		b.Walk(func(path string, info File) error {
			finfo, _ := info.FileInfo()
			if !finfo.IsDir() {
				keys = append(keys, finfo.Name())
			}
			return nil
		})
	} else {
		for k := range b.data {
			keys = append(keys, k)
		}
	}
	return keys
}

func (b *Box) indexDirectories() {
	b.directories = map[string]bool{}
	if _, ok := data[b.Path]; ok {
		for name, _ := range data[b.Path] {
			prefix, _ := path.Split(name)
			// Even on Windows the suffix appears to be a /
			prefix = strings.TrimSuffix(prefix, "/")
			b.directories[prefix] = true
		}
	}
}
