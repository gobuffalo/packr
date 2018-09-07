package packr

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/packr/file"
	"github.com/gobuffalo/packr/file/resolver"
	"github.com/pkg/errors"
)

func New(path string) Box {
	var cd string
	if !filepath.IsAbs(path) {
		_, filename, _, _ := runtime.Caller(1)
		cd = filepath.Dir(filename)
	}

	// this little hack courtesy of the `-cover` flag!!
	cov := filepath.Join("_test", "_obj_test")
	cd = strings.Replace(cd, string(filepath.Separator)+cov, "", 1)
	if !filepath.IsAbs(cd) && cd != "" {
		cd = filepath.Join(GoPath(), "src", cd)
	}
	b := Box{
		Path:       path,
		Name:       resolver.Ident(path),
		callingDir: resolver.Ident(cd),
	}
	return b
}

// Box represent a folder on a disk you want to
// have access to in the built Go binary.
type Box struct {
	Path       string
	Name       resolver.Ident
	callingDir resolver.Ident
	// data        map[string][]byte
	// directories map[string]bool
}

// AddString converts t to a byteslice and delegates to AddBytes to add to b.data
func (b Box) AddString(path string, t string) error {
	return b.AddBytes(path, []byte(t))
}

// AddBytes sets t in b.data by the given path
func (b Box) AddBytes(path string, t []byte) error {
	p, ok := resolver.DefaultResolver.(resolver.Packable)
	if !ok {
		return errors.Errorf("unable to pack %s - resolver.DefaultResolver is not resolver.Packable", path)
	}
	return p.Pack(resolver.Key(b.Name, resolver.Ident(path)), file.NewFile(path, t))
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
	f, err := resolver.Resolve(b.Name, resolver.Ident(name))
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(f)
}

// Has returns true if the resource exists in the box
func (b Box) Has(name string) bool {
	panic("not implemented")
	// _, err := b.find(name)
	// if err != nil {
	// 	return false
	// }
	// return true
}

// Open returns a File using the http.File interface
func (b Box) Open(name string) (http.File, error) {
	panic("not implemented")
	// return b.find(name)
}

// List shows "What's in the box?"
func (b Box) List() []string {
	panic("not implemented")
	// var keys []string
	//
	// if b.data == nil || len(b.data) == 0 {
	// 	b.Walk(func(path string, info File) error {
	// 		finfo, _ := info.FileInfo()
	// 		if !finfo.IsDir() {
	// 			keys = append(keys, finfo.Name())
	// 		}
	// 		return nil
	// 	})
	// } else {
	// 	for k := range b.data {
	// 		keys = append(keys, k)
	// 	}
	// }
	// return keys
}

// func (b Box) decompress(bb []byte) []byte {
// 	reader, err := gzip.NewReader(bytes.NewReader(bb))
// 	if err != nil {
// 		return bb
// 	}
// 	data, err := ioutil.ReadAll(reader)
// 	if err != nil {
// 		return bb
// 	}
// 	return data
// }

// func (b Box) find(name string) (File, error) {
// 	if bb, ok := b.data[name]; ok {
// 		return newVirtualFile(name, bb), nil
// 	}
// 	if b.directories == nil {
// 		b.indexDirectories()
// 	}
//
// 	cleanName := filepath.ToSlash(filepath.Clean(name))
// 	// Ensure name is not outside the box
// 	if strings.HasPrefix(cleanName, "../") {
// 		return nil, ErrResOutsideBox
// 	}
// 	// Absolute name is considered as relative to the box root
// 	cleanName = strings.TrimPrefix(cleanName, "/")
//
// 	// Try to get the resource from the box
// 	if _, ok := data[b.Path]; ok {
// 		if bb, ok := data[b.Path][cleanName]; ok {
// 			bb = b.decompress(bb)
// 			return newVirtualFile(cleanName, bb), nil
// 		}
// 		if _, ok := b.directories[cleanName]; ok {
// 			return newVirtualDir(cleanName), nil
// 		}
// 		if filepath.Ext(cleanName) != "" {
// 			// The Handler created by http.FileSystem checks for those errors and
// 			// returns http.StatusNotFound instead of http.StatusInternalServerError.
// 			return nil, os.ErrNotExist
// 		}
// 		return nil, os.ErrNotExist
// 	}
//
// 	// Not found in the box virtual fs, try to get it from the file system
// 	cleanName = filepath.FromSlash(cleanName)
// 	p := filepath.Join(b.callingDir, b.Path, cleanName)
// 	return fileFor(p, cleanName)
// }

// func (b *Box) indexDirectories() {
// 	b.directories = map[string]bool{}
// 	if _, ok := data[b.Path]; ok {
// 		for name := range data[b.Path] {
// 			prefix, _ := path.Split(name)
// 			// Even on Windows the suffix appears to be a /
// 			prefix = strings.TrimSuffix(prefix, "/")
// 			b.directories[prefix] = true
// 		}
// 	}
// }

// func fileFor(p string, name string) (File, error) {
// 	fi, err := os.Stat(p)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if fi.IsDir() {
// 		return newVirtualDir(p), nil
// 	}
// 	if bb, err := ioutil.ReadFile(p); err == nil {
// 		return newVirtualFile(name, bb), nil
// 	}
// 	return nil, os.ErrNotExist
// }
