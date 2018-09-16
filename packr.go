package packr

import (
	"sync"

	"github.com/gobuffalo/packr/file/resolver"
	"github.com/gobuffalo/packr/plog"
)

var boxes = map[string]*Box{}
var gil = &sync.RWMutex{}

func findBox(name string) *Box {
	gil.RLock()
	defer gil.RUnlock()
	return boxes[resolver.Key(name)]
}

func placeBox(b *Box) *Box {
	gil.Lock()
	plog.Debug(b, "placeBox", "name", b.Name, "path", b.Path, "resolution directory", b.ResolutionDir)
	boxes[resolver.Key(b.Name)] = b
	gil.Unlock()
	return b
}
