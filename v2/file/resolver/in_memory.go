package resolver

import (
	"os"
	"sync"

	"github.com/gobuffalo/packr/v2/file"
	"github.com/gobuffalo/packr/v2/plog"
)

var _ Resolver = &InMemory{}

type InMemory struct {
	files map[string]file.File
	moot  *sync.RWMutex
}

func (d *InMemory) Find(box string, name string) (file.File, error) {
	plog.Debug(d, "Find", "box", box, "name", name)
	d.moot.RLock()
	f, ok := d.files[name]
	d.moot.RUnlock()
	if ok {
		return f, nil
	}
	return nil, os.ErrNotExist
}

func (d *InMemory) Pack(name string, f file.File) error {
	plog.Debug(d, "Pack", "name", name)
	d.moot.Lock()
	d.files[name] = f
	d.moot.Unlock()
	return nil
}

func (d *InMemory) FileMap() map[string]file.File {
	d.moot.RLock()
	m := map[string]file.File{}
	for k, v := range d.files {
		m[Key(k)] = v
	}
	d.moot.RUnlock()
	return m
}

func NewInMemory(files map[string]file.File) *InMemory {
	if files == nil {
		files = map[string]file.File{}
	}
	return &InMemory{
		files: files,
		moot:  &sync.RWMutex{},
	}
}
