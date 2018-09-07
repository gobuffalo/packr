package resolver

import (
	"os"
	"sync"

	"github.com/gobuffalo/packr/file"
)

var _ Resolver = &inMemory{}

type inMemory struct {
	files map[Ident]file.File
	moot  *sync.RWMutex
}

func (d *inMemory) Find(name Ident) (file.File, error) {
	d.moot.RLock()
	defer d.moot.RUnlock()
	f, ok := d.files[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return f, nil
}

func NewInMemory(files map[Ident]file.File) Resolver {
	if files == nil {
		files = map[Ident]file.File{}
	}
	return &inMemory{
		files: files,
		moot:  &sync.RWMutex{},
	}
}
