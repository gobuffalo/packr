package packr

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/gobuffalo/packr/file"
	"github.com/gobuffalo/packr/file/resolver"
	"github.com/pkg/errors"
)

// NewBox returns a Box that can be used to
// retrieve files from either disk or the embedded
// binary.
func NewBox(path string) *Box {
	return New(path, path)
}

func New(name string, path string) *Box {
	iname := resolver.Ident(name)
	b := findBox(iname)
	if b != nil {
		return b
	}
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
	cd = filepath.Join(cd, path)
	b = &Box{
		Path:          path,
		Name:          iname,
		ResolutionDir: resolver.Ident(cd),
		resolvers:     map[string]resolver.Resolver{},
		moot:          &sync.RWMutex{},
	}
	return placeBox(b)
}

// Box represent a folder on a disk you want to
// have access to in the built Go binary.
type Box struct {
	Path          string // Path is deprecated and should no longer be used
	Name          resolver.Ident
	ResolutionDir resolver.Ident
	resolvers     map[string]resolver.Resolver
	moot          *sync.RWMutex
}

func (b *Box) SetResolver(file string, res resolver.Resolver) {
	b.moot.Lock()
	b.resolvers[resolver.Ident(file).Key()] = res
	b.moot.Unlock()
}

// AddString converts t to a byteslice and delegates to AddBytes to add to b.data
func (b *Box) AddString(path string, t string) error {
	return b.AddBytes(path, []byte(t))
}

// AddBytes sets t in b.data by the given path
func (b *Box) AddBytes(path string, t []byte) error {
	ipath := resolver.Ident(path)
	m := map[resolver.Ident]file.File{}
	m[ipath] = file.NewFile(path, t)
	res := resolver.NewInMemory(m)
	b.SetResolver(path, res)
	return nil
}

// String of the file asked for or an empty string.
func (b *Box) String(name string) string {
	return string(b.Bytes(name))
}

// MustString returns either the string of the requested
// file or an error if it can not be found.
func (b *Box) MustString(name string) (string, error) {
	bb, err := b.MustBytes(name)
	return string(bb), err
}

// Bytes of the file asked for or an empty byte slice.
func (b *Box) Bytes(name string) []byte {
	bb, _ := b.MustBytes(name)
	return bb
}

// MustBytes returns either the byte slice of the requested
// file or an error if it can not be found.
func (b *Box) MustBytes(name string) ([]byte, error) {
	f, err := b.resolve(resolver.Ident(name))
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(f)
}

// Has returns true if the resource exists in the box
func (b *Box) Has(name string) bool {
	_, err := b.MustBytes(name)
	if err != nil {
		return false
	}
	return true
}

// Open returns a File using the http.File interface
func (b *Box) Open(name string) (http.File, error) {
	return b.resolve(resolver.Ident(name))
}

// List shows "What's in the box?"
func (b *Box) List() []string {
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

func (b *Box) resolve(key resolver.Ident) (file.File, error) {
	b.moot.RLock()
	r, ok := b.resolvers[key.Key()]
	b.moot.RUnlock()
	if !ok {
		r = resolver.DefaultResolver
		if r == nil {
			return nil, errors.New("resolver.DefaultResolver is nil")
		}
	}
	fmt.Println(b.Name, key, fmt.Sprintf("using resolver - %T", r))

	f, err := r.Find(key)
	if err != nil {
		z := filepath.Join(b.ResolutionDir.OsPath(), key.OsPath())
		f, err = r.Find(resolver.Ident(z))
		if err != nil {
			return f, errors.WithStack(err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return f, errors.WithStack(err)
		}
		f = file.NewFile(key.Name(), b)
	}
	return f, nil
}
