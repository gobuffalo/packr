package resolver

import (
	"os"

	"github.com/gobuffalo/packr/v2/file"
)

type Resolver interface {
	Resolve(string, string) (file.File, error)
}

func defaultResolver() Resolver {
	pwd, _ := os.Getwd()
	return &Disk{
		Root: pwd,
	}
}

var DefaultResolver = defaultResolver()
