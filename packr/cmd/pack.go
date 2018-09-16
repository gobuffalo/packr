package cmd

import (
	"fmt"
	"os"

	"github.com/gobuffalo/packr/jam/parser"
	"github.com/gobuffalo/packr/jam/store"
	"github.com/pkg/errors"
)

func pack(args ...string) error {
	if err := clean(args...); err != nil {
		return errors.WithStack(err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		return errors.WithStack(err)
	}
	roots := append(args, pwd)
	for _, r := range roots {
		store.Clean(r)
	}
	p, err := parser.NewFromRoots(roots, &parser.RootsOptions{
		IgnoreImports: globalOptions.IgnoreImports,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	boxes, err := p.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	// reduce boxes - remove ones we don't want
	// MB: current assumption is we want all these
	// boxes, just adding a comment suggesting they're
	// might be a reason to exclude some

	fmt.Printf("Found %d boxes\n", len(boxes))

	// "pack" boxes
	d := store.NewDisk("", "")
	for _, b := range boxes {
		if b.Name == store.DISK_GLOBAL_KEY {
			continue
		}
		fmt.Println("box", b.Name, b.AbsPath)
		if err := d.Pack(b); err != nil {
			return errors.WithStack(err)
		}
	}
	return d.Close()
}
