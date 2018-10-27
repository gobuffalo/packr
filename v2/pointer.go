package packr

import (
	"io/ioutil"

	"github.com/gobuffalo/packr/v2/file"
	"github.com/gobuffalo/packr/v2/file/resolver"
	"github.com/pkg/errors"
)

type Pointer struct {
	ForwardBox  string
	ForwardPath string
}

var _ resolver.Resolver = Pointer{}

func (p Pointer) Find(box string, path string) (file.File, error) {
	b := findBox(p.ForwardBox)
	f, err := b.Resolve(p.ForwardPath)
	if err != nil {
		return f, errors.WithStack(err)
	}
	f.Seek(0, 0)
	x, err := ioutil.ReadAll(f)
	if err != nil {
		return f, errors.WithStack(err)
	}
	return file.NewFile(path, x)
}
