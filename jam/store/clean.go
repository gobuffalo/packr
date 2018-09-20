package store

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr/jam/parser"
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
		if info.IsDir() && base == "packrd" {
			os.Remove(path)
			return filepath.SkipDir
		}

		if strings.Contains(base, "-packr.go") {
			err := os.Remove(path)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	}
	return godirwalk.Walk(root, &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback:            callback,
	})
}
