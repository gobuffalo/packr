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

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2/file"
	"github.com/gobuffalo/packr/v2/file/resolver"
	"github.com/gobuffalo/packr/v2/plog"
	"github.com/markbates/oncer"
	"github.com/pkg/errors"
)

var _ packd.Box = &Box{}
var _ packd.HTTPBox = Box{}
var _ packd.Addable = &Box{}
var _ packd.Walkable = &Box{}
var _ packd.Finder = Box{}

// NewBox returns a Box that can be used to
// retrieve files from either disk or the embedded
// binary.
func NewBox(path string) Box {
	oncer.Deprecate(0, "packr.NewBox", "Use packr.New instead.")
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
		cd = filepath.Join(envy.GoPath(), "src", cd)
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
	Path            string
	Name            string
	ResolutionDir   string
	resolvers       map[string]resolver.Resolver
	DefaultResolver resolver.Resolver
	moot            *sync.RWMutex
}

func (b *Box) SetResolver(file string, res resolver.Resolver) {
	b.moot.Lock()
	plog.Debug(b, "SetResolver", "file", file, "resolver", fmt.Sprintf("%T", res))
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
	f, err := file.NewFile(path, t)
	if err != nil {
		return errors.WithStack(err)
	}
	m[resolver.Key(path)] = f
	res := resolver.NewInMemory(m)
	b.SetResolver(path, res)
	return nil
}

// String is deprecated. Use FindString instead
func (b Box) String(name string) string {
	oncer.Deprecate(0, "github.com/gobuffalo/packr/v2#Box.String", "Use github.com/gobuffalo/packr/v2#Box.FindString instead.")
	return string(b.Bytes(name))
}

// MustString is deprecated. Use FindString instead
func (b Box) MustString(name string) (string, error) {
	oncer.Deprecate(0, "github.com/gobuffalo/packr/v2#Box.MustString", "Use github.com/gobuffalo/packr/v2#Box.FindString instead.")
	return b.FindString(name)
}

// FindString returns either the string of the requested
// file or an error if it can not be found.
func (b Box) FindString(name string) (string, error) {
	bb, err := b.Find(name)
	return string(bb), err
}

// Bytes is deprecated. Use Find instead
func (b Box) Bytes(name string) []byte {
	bb, _ := b.Find(name)
	oncer.Deprecate(0, "github.com/gobuffalo/packr/v2#Box.Bytes", "Use github.com/gobuffalo/packr/v2#Box.Find instead.")
	return bb
}

// MustBytes is deprecated. Use Find instead.
func (b Box) MustBytes(name string) ([]byte, error) {
	oncer.Deprecate(0, "github.com/gobuffalo/packr/v2#Box.MustBytes", "Use github.com/gobuffalo/packr/v2#Box.Find instead.")
	return b.Find(name)
}

// Find returns either the byte slice of the requested
// file or an error if it can not be found.
func (b Box) Find(name string) ([]byte, error) {
	f, err := b.Resolve(name)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(f)
}

// Has returns true if the resource exists in the box
func (b Box) Has(name string) bool {
	_, err := b.Find(name)
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
		if info == nil {
			return nil
		}
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
	plog.Debug(b, "Resolve", "key", key)

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
		f, err = file.NewFile(key, b)
		if err != nil {
			return f, errors.WithStack(err)
		}
	}
	return f, nil
}
