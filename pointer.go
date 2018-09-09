package packr

import (
	"github.com/gobuffalo/packr/file"
)

type Pointer struct {
	ForwardBox  string
	ForwardPath string
}

func (p *Pointer) Find(box string, path string) (file.File, error) {
	b := findBox(p.ForwardBox)
	return b.Resolve(p.ForwardPath)
}
