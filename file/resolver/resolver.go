package resolver

import "github.com/gobuffalo/packr/file"

type Resolver interface {
	Find(Ident) (file.File, error)
}
