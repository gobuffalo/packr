package resolver

import (
	"fmt"
	"path"
	"sync"

	"github.com/gobuffalo/packr/file"
)

var gil = &sync.RWMutex{}
var resolutions = map[Ident]Resolver{}

func Register(box Ident, file Ident, res Resolver) {
	gil.Lock()
	defer gil.Unlock()
	resolutions[Key(box, file)] = res
}

func Key(box Ident, file Ident) Ident {
	return Ident(path.Join(box.Key(), file.Key()))
}

func Resolve(box Ident, file Ident) (file.File, error) {
	gil.RLock()
	key := Key(box, file)
	fmt.Println("### resolving key ->", key)
	if r, ok := resolutions[key]; ok {
		fmt.Println(key, "found in resolutions")
		defer gil.RUnlock()
		return r.Find(file)
	}
	fmt.Println(key, "not found in resolutions")

	if r, err := DefaultResolver.Find(key); err == nil {
		fmt.Println(key, "found in DefaultResolver")
		defer gil.RUnlock()
		return r, nil
	}
	fmt.Println(key, "not found in DefaultResolver")

	gil.RUnlock()
	d := &Disk{
		Root: box,
	}
	return d.Find(file)
}

func ClearRegistry() {
	gil.Lock()
	defer gil.Unlock()
	resolutions = map[Ident]Resolver{}
	DefaultResolver = NewInMemory(map[Ident]file.File{})
}
