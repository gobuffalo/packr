package store

import (
	"os"
	"path/filepath"
	"strings"

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
	callback := func(path string, info *godirwalk.Dirent) error {
		base := filepath.Base(path)
		if base == ".git" || base == "vendor" || base == "node_modules" {
			return filepath.SkipDir
		}
		if info == nil || info.IsDir() {
			return nil
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
