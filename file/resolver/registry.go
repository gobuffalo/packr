package resolver

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/gobuffalo/packr/file"
	"github.com/pkg/errors"
)

var gil = &sync.RWMutex{}
var resolutions = map[Ident]map[Ident]Resolver{}

func Register(box Ident, file Ident, res Resolver) error {
	br := BoxResolvers(box)

	gil.Lock()
	br[file] = res
	resolutions[box] = br
	gil.Unlock()
	return nil
}

func BoxResolvers(box Ident) map[Ident]Resolver {
	gil.RLock()
	br, ok := resolutions[box]
	gil.RUnlock()
	if !ok {
		br = map[Ident]Resolver{}
		gil.Lock()
		resolutions[box] = br
		gil.Unlock()
	}

	m := map[Ident]Resolver{}

	for k, v := range br {
		m[k] = v
	}

	return m
}

func Resolve(box Ident, name Ident) (file.File, error) {
	br := BoxResolvers(box)

	r, ok := br[name]
	if !ok {
		if DefaultResolver == nil {
			return nil, errors.New("resolver.DefaultResolver is nil")
		}
		r = DefaultResolver
	}
	fmt.Println(filepath.Join(box.OsPath(), name.OsPath()), fmt.Sprintf("using resolver - %T", r))

	f, err := r.Find(name)
	if err != nil {
		f, err = r.Find(Ident(filepath.Join(box.OsPath(), name.OsPath())))
		if err != nil {
			return f, errors.WithStack(err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return f, errors.WithStack(err)
		}
		f = file.NewFile(name.Name(), b)
	}
	return f, nil
}

func ClearRegistry() {
	gil.Lock()
	defer gil.Unlock()
	resolutions = map[Ident]map[Ident]Resolver{}
	DefaultResolver = defaultResolver()
}
