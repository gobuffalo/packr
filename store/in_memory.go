package store

import (
	"sync"

	"github.com/gobuffalo/packr/file"
)

var _ Store = &inMemory{}

type inMemory struct {
	files map[string]file.File
	moot  *sync.RWMutex
}

func (d *inMemory) Pack(name string, f file.File) error {
	d.moot.RLock()
	defer d.moot.RUnlock()
	d.files[name] = f
	return nil
}

var _ file.FileMappable = &inMemory{}

func (d *inMemory) FileMap() map[string]file.File {
	d.moot.RLock()
	defer d.moot.RUnlock()
	m := map[string]file.File{}
	for k, v := range d.files {
		m[k] = v
	}
	return m
}

func (d *inMemory) Close() error {
	return nil
}

func NewInMemory() Store {
	return &inMemory{
		files: map[string]file.File{},
		moot:  &sync.RWMutex{},
	}
}
