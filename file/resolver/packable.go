package resolver

import "github.com/gobuffalo/packr/file"

type Packable interface {
	Pack(name Ident, f file.File) error
}
