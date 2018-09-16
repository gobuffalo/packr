package parser

import (
	"go/build"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr/plog"
	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

type finder struct {
	seen map[string]string
}

// findAllGoFiles *.go files for a given diretory
func (fd *finder) findAllGoFiles(dir string) ([]string, error) {
	plog.Debug(fd, "findAllGoFiles", "dir", dir)
	var names []string

	callback := func(path string, do *godirwalk.Dirent) error {
		ext := filepath.Ext(path)
		if ext != ".go" {
			return nil
		}
		names = append(names, path)
		return nil
	}
	err := godirwalk.Walk(dir, &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback:            callback,
	})

	return names, err
}

func (fd *finder) findAllGoFilesImports(dir string) ([]string, error) {
	var names []string

	ctx := build.Default

	if len(ctx.SrcDirs()) == 0 {
		return names, errors.New("no src directories found")
	}

	pkg, err := ctx.ImportDir(dir, 0)

	if err != nil {
		if !strings.Contains(err.Error(), "cannot find package") {
			if _, ok := errors.Cause(err).(*build.NoGoError); !ok {
				return names, errors.WithStack(err)
			}
		}
	}

	if pkg.Goroot {
		return names, nil
	}
	if len(pkg.GoFiles) <= 0 {
		return names, nil
	}

	plog.Debug(fd, "findAllGoFilesImports", "dir", dir)

	names, _ = fd.findAllGoFiles(dir)
	for _, n := range pkg.GoFiles {
		names = append(names, filepath.Join(pkg.Dir, n))
	}
	for _, imp := range pkg.Imports {
		if _, ok := fd.seen[imp]; ok {
			continue
		}
		fd.seen[imp] = imp
		if len(ctx.SrcDirs()) == 0 {
			continue
		}
		d := ctx.SrcDirs()[len(ctx.SrcDirs())-1]
		ip := filepath.Join(d, imp)
		n, err := fd.findAllGoFilesImports(ip)
		if err != nil {
			return n, nil
		}
		names = append(names, n...)
	}
	return names, nil
}
