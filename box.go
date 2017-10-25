package packr

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"compress/gzip"

	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
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
	Path       string
	callingDir string
	data       map[string][]byte
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
	p := filepath.Join(b.callingDir, b.Path, name)
	return ioutil.ReadFile(p)
}

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
	name = strings.TrimPrefix(name, "/")
	name = filepath.ToSlash(name)
	if _, ok := data[b.Path]; ok {
		if bb, ok := data[b.Path][name]; ok {
			bb = b.decompress(bb)
			return newVirtualFile(name, bb), nil
		}
		if filepath.Ext(name) != "" {
			return nil, errors.Errorf("could not find virtual file: %s", name)
		}
		return newVirtualDir(name), nil
	}

	p := filepath.Join(b.callingDir, b.Path, name)
	if f, err := os.Open(p); err == nil {
		return physicalFile{f}, nil
	}
	// make one last ditch effort to find the file below the PWD:
	pwd, _ := os.Getwd()
	p = filepath.Join(pwd, b.Path, name)
	if f, err := os.Open(p); err == nil {
		return physicalFile{f}, nil
	}
	return nil, errors.Errorf("could not find %s in box %s", name, b.Path)
}

type WalkFunc func(string, File) error

func (b Box) Walk(wf WalkFunc) error {
	if data[b.Path] == nil {
		base := filepath.Join(b.callingDir, b.Path)
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
