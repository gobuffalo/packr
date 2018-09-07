package resolver

import (
	"path"
	"sync"

	"github.com/gobuffalo/packr/file"
)

var gil = &sync.RWMutex{}
var resolutions = map[Ident]Resolver{}

func Register(box Ident, file Ident, res Resolver) {
	gil.Lock()
	defer gil.Unlock()
	resolutions[resKey(box, file)] = res
}

func resKey(box Ident, file Ident) Ident {
	return Ident(path.Join(box.Key(), file.Key()))
}

func Resolve(box Ident, file Ident) (file.File, error) {
	gil.RLock()
	if r, ok := resolutions[resKey(box, file)]; ok {
		defer gil.RUnlock()
		return r.Find(file)
	}
	gil.RUnlock()
	d := &Disk{
		Root: box,
	}
	Register(box, file, d)
	return d.Find(file)
}

func ClearRegistry() {
	gil.Lock()
	defer gil.Unlock()
	resolutions = map[Ident]Resolver{}
}
