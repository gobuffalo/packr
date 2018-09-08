package packr

import (
	"sync"

	"github.com/gobuffalo/packr/file/resolver"
)

var boxes = map[string]*Box{}
var gil = &sync.RWMutex{}

func findBox(name resolver.Ident) *Box {
	gil.RLock()
	defer gil.RUnlock()
	return boxes[name.Key()]
}

func placeBox(b *Box) *Box {
	gil.Lock()
	defer gil.Unlock()
	boxes[b.Name.Key()] = b
	return b
}
