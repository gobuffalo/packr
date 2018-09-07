package store

import "github.com/gobuffalo/packr/file"

type Store interface {
	Pack(name string, f file.File) error
	Close() error
}
