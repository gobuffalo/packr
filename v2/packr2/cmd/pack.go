package cmd

import (
	"io"
	"os"

	"github.com/gobuffalo/packr/v2/jam/parser"
	"github.com/gobuffalo/packr/v2/jam/store"
	"github.com/gobuffalo/packr/v2/plog"
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

	plog.Default.Debugf("found %d boxes", len(boxes))

	var st store.Store = store.NewDisk("", "")

	if globalOptions.Legacy {
		st = store.NewLegacy()
	}

	for _, b := range boxes {
		if b.Name == store.DISK_GLOBAL_KEY {
			continue
		}
		if err := st.Pack(b); err != nil {
			return errors.WithStack(err)
		}
	}
	if cl, ok := st.(io.Closer); ok {
		return cl.Close()
	}
	return nil
}
