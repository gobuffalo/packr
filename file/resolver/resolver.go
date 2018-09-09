package resolver

import (
	"os"

	"github.com/gobuffalo/packr/file"
)

type Resolver interface {
	Find(string, string) (file.File, error)
}

func defaultResolver() Resolver {
	pwd, _ := os.Getwd()
	return &Disk{
		Root: pwd,
	}
}

var DefaultResolver = defaultResolver()
