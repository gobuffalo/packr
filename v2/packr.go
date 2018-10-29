package packr

import (
	"sync"

	"github.com/gobuffalo/packr/v2/file/resolver"
	"github.com/gobuffalo/packr/v2/plog"
	"github.com/pkg/errors"
)

var boxes = map[string]*Box{}
var gil = &sync.RWMutex{}

func findBox(name string) (*Box, error) {
	plog.Debug("packr", "findBox", "name", name)
	gil.RLock()
	defer gil.RUnlock()
	b, ok := boxes[resolver.Key(name)]
	if !ok {
		return nil, errors.Errorf("could not find box %s", name)
	}
	plog.Debug(b, "findBox", "box", b)
	return b, nil
}

func placeBox(b *Box) *Box {
	gil.Lock()
	plog.Debug(b, "placeBox", "name", b.Name, "path", b.Path, "resolution directory", b.ResolutionDir)
	boxes[resolver.Key(b.Name)] = b
	gil.Unlock()
	return b
}
