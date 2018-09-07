package store

import "github.com/gobuffalo/packr/file"

type Packable interface {
	Pack(name string, f file.File) error
}

type Store interface {
	Packable
	Close() error
}
