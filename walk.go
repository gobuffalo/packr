package packr

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/gobuffalo/packr/file"
	"github.com/gobuffalo/packr/file/resolver"
	"github.com/pkg/errors"
)

type WalkFunc func(string, file.File) error

// Walk will traverse the box and call the WalkFunc for each file in the box/folder.
func (b Box) Walk(wf WalkFunc) error {
	m := map[string]file.File{}

	cd := filepath.Join(b.ResolutionDir.OsPath(), b.Name.OsPath())
	d := &resolver.Disk{Root: resolver.Ident(cd)}
	for n, f := range d.FileMap() {
		m[resolver.Ident(n).Name()] = f
	}

	res := resolver.BoxResolvers(b.Name)

	for n, r := range res {
		f, err := r.Find(n)
		if err != nil {
			return errors.WithStack(err)
		}
		m[n.Name()] = f
	}

	var keys = make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if err := wf(k, m[k]); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// WalkPrefix will call box.Walk and call the WalkFunc when it finds paths that have a matching prefix
func (b Box) WalkPrefix(prefix string, wf WalkFunc) error {
	ipref := resolver.Ident(prefix).OsPath()
	return b.Walk(func(path string, f File) error {
		ipath := resolver.Ident(path).OsPath()
		if strings.HasPrefix(ipath, ipref) {
			if err := wf(path, f); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})
}
