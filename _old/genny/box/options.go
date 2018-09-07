package box

import (
	"os"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

type Options struct {
	Package string // the name of the package that box is in
	Name    string // name of the box (i.e. "./foo/templates")
	Root    string // the directory with the *-packr.go file is written.
	Files   []genny.File
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if len(opts.Name) == 0 {
		return errors.New("a box must have a name")
	}

	if len(opts.Package) == 0 {
		return errors.New("a box must have a package")
	}

	if len(opts.Root) == 0 {
		pwd, err := os.Getwd()
		if err != nil {
			return errors.WithStack(err)
		}
		opts.Root = pwd
	}

	return nil
}
