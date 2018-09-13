package packr

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/gobuffalo/packr/file"
	"github.com/gobuffalo/packr/file/resolver"
	"github.com/gobuffalo/packr/plog"
	"github.com/pkg/errors"
)

// NewBox returns a Box that can be used to
// retrieve files from either disk or the embedded
// binary.
func NewBox(path string) Box {
	return *New(path, path)
}

func New(name string, path string) *Box {
	b := findBox(name)
	if b != nil {
		return b
	}
	var cd string
	if !filepath.IsAbs(path) {
		_, filename, _, _ := runtime.Caller(2)
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
		Name:          name,
		ResolutionDir: cd,
		resolvers:     map[string]resolver.Resolver{},
		moot:          &sync.RWMutex{},
	}
	return placeBox(b)
}

// Box represent a folder on a disk you want to
// have access to in the built Go binary.
type Box struct {
	Path            string // Path is deprecated and should no longer be used
	Name            string
	ResolutionDir   string
	resolvers       map[string]resolver.Resolver
	DefaultResolver resolver.Resolver
	moot            *sync.RWMutex
}

func (b *Box) SetResolver(file string, res resolver.Resolver) {
	b.moot.Lock()
	b.resolvers[resolver.Key(file)] = res
	b.moot.Unlock()
}

// AddString converts t to a byteslice and delegates to AddBytes to add to b.data
func (b *Box) AddString(path string, t string) error {
	return b.AddBytes(path, []byte(t))
}

// AddBytes sets t in b.data by the given path
func (b *Box) AddBytes(path string, t []byte) error {
	m := map[string]file.File{}
	m[resolver.Key(path)] = file.NewFile(path, t)
	res := resolver.NewInMemory(m)
	b.SetResolver(path, res)
	return nil
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
	f, err := b.Resolve(name)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(f)
}

// Has returns true if the resource exists in the box
func (b Box) Has(name string) bool {
	_, err := b.MustBytes(name)
	if err != nil {
		return false
	}
	return true
}

// Open returns a File using the http.File interface
func (b Box) Open(name string) (http.File, error) {
	return b.Resolve(name)
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

func (b *Box) Resolve(key string) (file.File, error) {
	b.moot.RLock()
	r, ok := b.resolvers[resolver.Key(key)]
	b.moot.RUnlock()
	if !ok {
		r = b.DefaultResolver
		if r == nil {
			r = resolver.DefaultResolver
			if r == nil {
				return nil, errors.New("resolver.DefaultResolver is nil")
			}
		}
	}
	plog.Debugf("resolving %q, %q: %T", b.Name, key, r)

	f, err := r.Find(b.Name, key)
	if err != nil {
		z := filepath.Join(resolver.OsPath(b.ResolutionDir), resolver.OsPath(key))
		f, err = r.Find(b.Name, z)
		if err != nil {
			return f, errors.WithStack(err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return f, errors.WithStack(err)
		}
		f = file.NewFile(key, b)
	}
	return f, nil
}
