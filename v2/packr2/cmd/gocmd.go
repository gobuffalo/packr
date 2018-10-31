package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2/jam/parser"
	"github.com/gobuffalo/packr/v2/plog"
	"github.com/pkg/errors"
)

func goCmd(name string, args ...string) error {
	cargs := []string{name}
	cargs = append(cargs, args...)
	if len(args) > 0 {
		err := func() error {
			fi, err := os.Stat(args[len(args)-1])
			if err != nil {
				return errors.WithStack(err)
			}
			path := fi.Name()
			if fi.IsDir() {
				return nil
			}
			path, err = filepath.Abs(filepath.Dir(path))
			if err != nil {
				return errors.WithStack(err)
			}

			p, err := parser.NewFromRoots([]string{path}, nil)
			if err != nil {
				return errors.WithStack(err)
			}

			boxes, err := p.Run()
			if err != nil {
				return errors.WithStack(err)
			}
			for _, b := range boxes {
				if b.PackageDir == path {
					pk := fmt.Sprintf("%s-packr.go", b.Package)
					for _, x := range []string{pk, "a_" + pk} {
						y := x
						if _, err := os.Stat(y); err != nil {
							continue
						}
						cargs = append(cargs, y)
						break
					}
				}
			}
			return nil
		}()
		if err != nil {
			return errors.WithStack(err)
		}
	}
	cp := exec.Command(genny.GoBin(), cargs...)
	plog.Logger.Debug(strings.Join(cp.Args, " "))
	cp.Stderr = os.Stderr
	cp.Stdin = os.Stdin
	cp.Stdout = os.Stdout
	return cp.Run()
}
