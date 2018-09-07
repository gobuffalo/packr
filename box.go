package packr

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/gobuffalo/packr/file"
	"github.com/gobuffalo/packr/file/resolver"
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
	Path       string // Path is deprecated and should no longer be used
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
	ipath := resolver.Ident(path)
	m := map[resolver.Ident]file.File{}
	m[ipath] = file.NewFile(path, t)
	res := resolver.NewInMemory(m)
	return resolver.Register(b.Name, ipath, res)
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
	_, err := resolver.Resolve(b.Name, resolver.Ident(name))
	if err != nil {
		return false
	}
	return true
}

// Open returns a File using the http.File interface
func (b Box) Open(name string) (http.File, error) {
	return resolver.Resolve(b.Name, resolver.Ident(name))
}

// List shows "What's in the box?"
func (b Box) List() []string {
	var keys []string

	b.Walk(func(path string, info File) error {
		finfo, _ := info.FileInfo()
		if !finfo.IsDir() {
			keys = append(keys, finfo.Name())
		}
		return nil
	})
	sort.Strings(keys)
	return keys
}
