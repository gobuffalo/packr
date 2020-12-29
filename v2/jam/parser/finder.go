package parser

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobuffalo/packr/v2/plog"
	"github.com/karrick/godirwalk"
	"github.com/markbates/errx"
	"github.com/markbates/oncer"
)

type finder struct {
	id time.Time
}

func (fd *finder) key(m, dir string) string {
	return fmt.Sprintf("%s-*parser.finder#%s-%s", fd.id, m, dir)
}

// findAllGoFiles *.go files for a given diretory
func (fd *finder) findAllGoFiles(dir string) ([]string, error) {
	var err error
	var names []string
	oncer.Do(fd.key("findAllGoFiles", dir), func() {
		plog.Debug(fd, "findAllGoFiles", "dir", dir)

		callback := func(path string, do *godirwalk.Dirent) error {
			if filepath.Base(path) == "vendor" {
				return filepath.SkipDir
			}
			ext := filepath.Ext(path)
			if ext != ".go" {
				return nil
			}
			//check if path is a dir
			fi, err := os.Stat(path)
			if err != nil {
				return nil
			}

			if fi.IsDir() {
				return nil
			}

			names = append(names, path)
			return nil
		}
		err = godirwalk.Walk(dir, &godirwalk.Options{
			FollowSymbolicLinks: true,
			Callback:            callback,
		})
	})

	return names, err
}

func (fd *finder) findAllGoFilesImports(dir string) ([]string, error) {
	return fd.findAllGoFilesImportsIn(".", dir)
}

func (fd *finder) findAllGoFilesImportsIn(path, dir string) ([]string, error) {
	var err error
	var names []string
	oncer.Do(fd.key("findAllGoFilesImports", filepath.Join(dir, path)), func() {
		ctx := build.Default

		if len(ctx.SrcDirs()) == 0 {
			err = fmt.Errorf("no src directories found")
			return
		}

		pkg, err := ctx.Import(path, dir, 0)
		if strings.Contains(pkg.ImportPath, "github.com/gobuffalo/packr") {
			return
		}

		if err != nil {
			if !strings.Contains(err.Error(), "cannot find package") {
				if _, ok := errx.Unwrap(err).(*build.NoGoError); !ok {
					err = err
					return
				}
			}
		}

		if pkg.Goroot {
			return
		}
		if len(pkg.GoFiles) <= 0 {
			return
		}

		plog.Debug(fd, "findAllGoFilesImports", "dir", dir)

		names, _ = fd.findAllGoFiles(pkg.Dir)
		for _, n := range pkg.GoFiles {
			names = append(names, filepath.Join(pkg.Dir, n))
		}
		for _, imp := range pkg.Imports {
			n, err := fd.findAllGoFilesImportsIn(imp, dir)
			if err != nil && len(n) != 0 {
				names = n
				return
			}
			names = append(names, n...)
		}
	})
	return names, err
}
