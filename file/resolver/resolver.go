package resolver

import (
	"os"

	"github.com/gobuffalo/packr/file"
)

type Resolver interface {
	Find(Ident) (file.File, error)
}

func defaultResolver() Resolver {
	pwd, _ := os.Getwd()
	return &Disk{
		Root: Ident(pwd),
	}
}

var DefaultResolver = defaultResolver()
