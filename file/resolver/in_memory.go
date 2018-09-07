package resolver

import (
	"os"
	"sync"

	"github.com/gobuffalo/packr/file"
)

var _ Resolver = &InMemory{}

type InMemory struct {
	files map[Ident]file.File
	moot  *sync.RWMutex
}

func (d *InMemory) Find(name Ident) (file.File, error) {
	d.moot.RLock()
	defer d.moot.RUnlock()
	f, ok := d.files[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return f, nil
}

func (d *InMemory) Pack(name Ident, f file.File) error {
	d.moot.Lock()
	defer d.moot.Unlock()
	d.files[name] = f
	return nil
}

func (d *InMemory) FileMap() map[string]file.File {
	d.moot.RLock()
	defer d.moot.RUnlock()
	m := map[string]file.File{}
	for k, v := range d.files {
		m[k.Name()] = v
	}
	return m
}

func NewInMemory(files map[Ident]file.File) *InMemory {
	if files == nil {
		files = map[Ident]file.File{}
	}
	return &InMemory{
		files: files,
		moot:  &sync.RWMutex{},
	}
}
