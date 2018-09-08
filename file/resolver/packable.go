package resolver

import "github.com/gobuffalo/packr/file"

type Packable interface {
	Pack(name string, f file.File) error
}
