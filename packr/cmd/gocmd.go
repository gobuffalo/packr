package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/packr/jam/parser"
	"github.com/pkg/errors"
)

func goCmd(name string, args ...string) error {
	cargs := []string{name}
	cargs = append(cargs, args...)
	if len(args) > 0 {
		fi, err := os.Stat(args[len(args)-1])
		if err != nil {
			return errors.WithStack(err)
		}
		path := fi.Name()
		if !fi.IsDir() {
			path = filepath.Dir(path)
		}
		path, err = filepath.Abs(path)
		if err != nil {
			return errors.WithStack(err)
		}

		p, err := parser.NewFromRoots([]string{path})
		if err != nil {
			return errors.WithStack(err)
		}

		boxes, err := p.Run()
		if err != nil {
			return errors.WithStack(err)
		}
		for _, b := range boxes {
			if b.PackageDir == path {
				pk := fmt.Sprintf("a_%s-packr.go", b.Package)
				filepath.Join(path, pk)
				cargs = append(cargs, pk)
			}
		}
	}
	cp := exec.Command(packr.GoBin(), cargs...)
	cp.Stderr = os.Stderr
	cp.Stdin = os.Stdin
	cp.Stdout = os.Stdout
	return cp.Run()
}
