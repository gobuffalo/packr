package resolver

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

var _ Resolver = &DiskResolver{}

type DiskResolver struct {
	Roots     []string
	prospects []genny.File
	boxes     map[string][]genny.File
}

func (d *DiskResolver) Resolve() error {
	if len(d.Roots) == 0 {
		pwd, err := os.Getwd()
		if err != nil {
			return errors.WithStack(err)
		}
		d.Roots = []string{pwd}
	}
	if d.boxes == nil {
		d.boxes = map[string][]genny.File{}
	}

	for _, root := range d.Roots {
		err := godirwalk.Walk(root, &godirwalk.Options{
			FollowSymbolicLinks: true,
			Callback:            texasRanger(root, d),
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (d *DiskResolver) Prospects() []genny.File {
	return d.prospects
}

func (d *DiskResolver) Boxes() map[string][]genny.File {
	return d.boxes
}

func texasRanger(root string, d *DiskResolver) godirwalk.WalkFunc {
	return func(path string, info *godirwalk.Dirent) error {
		var r io.Reader
		if !info.IsDir() {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.WithStack(err)
			}
			r = bytes.NewReader(b)
		}
		f := genny.NewFile(path, r)

		if !IsProspect(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}

		d.prospects = append(d.prospects, f)
		return nil
	}
}
