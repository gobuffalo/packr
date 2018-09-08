package resolver

import (
	"fmt"
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
	fmt.Println("InMemory: Find", name)
	d.moot.RLock()
	f, ok := d.files[name]
	d.moot.RUnlock()
	if ok {
		return f, nil
	}
	return nil, os.ErrNotExist
}

func (d *InMemory) Pack(name Ident, f file.File) error {
	d.moot.Lock()
	d.files[name] = f
	d.moot.Unlock()
	return nil
}

func (d *InMemory) FileMap() map[string]file.File {
	d.moot.RLock()
	m := map[string]file.File{}
	for k, v := range d.files {
		m[k.Name()] = v
	}
	d.moot.RUnlock()
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
