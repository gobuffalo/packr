package store

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr/v2/jam/parser"
	"github.com/pkg/errors"

	"github.com/karrick/godirwalk"
)

func Clean(root string) error {
	if len(root) == 0 {
		pwd, err := os.Getwd()
		if err != nil {
			return errors.WithStack(err)
		}
		root = pwd
	}
	if _, err := os.Stat(root); err != nil {
		return nil
	}
	callback := func(path string, info *godirwalk.Dirent) error {
		if _, err := os.Stat(path); err != nil {
			return nil
		}
		base := filepath.Base(path)
		for _, d := range parser.DefaultIgnoredFolders {
			if base == d {
				return filepath.SkipDir
			}
		}
		if info == nil {
			return nil
		}
		for _, suf := range []string{"-packr.go", "packrd"} {
			if strings.HasSuffix(base, suf) {
				err := os.RemoveAll(path)
				if err != nil {
					return errors.WithStack(err)
				}
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		return nil
	}
	return godirwalk.Walk(root, &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback:            callback,
	})
}
