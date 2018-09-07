package store

import "github.com/gobuffalo/packr/file"

var _ Store = &Disk{}

type Disk struct {
}

func (d *Disk) Pack(name string, f file.File) error {
	panic("not implemented")
}

func (d *Disk) Close() error {
	panic("not implemented")
}
