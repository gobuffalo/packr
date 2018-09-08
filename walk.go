package packr

import (
	"sort"
	"strings"

	"github.com/gobuffalo/packr/file"
	"github.com/gobuffalo/packr/file/resolver"
	"github.com/pkg/errors"
)

type WalkFunc func(string, file.File) error

// Walk will traverse the box and call the WalkFunc for each file in the box/folder.
func (b *Box) Walk(wf WalkFunc) error {
	m := map[string]file.File{}

	cd := b.ResolutionDir.OsPath()
	d := &resolver.Disk{Root: resolver.Ident(cd)}
	for n, f := range d.FileMap() {
		m[resolver.Ident(n).Name()] = f
	}

	b.moot.RLock()
	for n, r := range b.resolvers {
		iname := resolver.Ident(n)
		f, err := r.Find(iname)
		if err != nil {
			return errors.WithStack(err)
		}
		m[n] = f
	}
	b.moot.RUnlock()

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
