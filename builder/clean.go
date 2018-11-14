package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Clean up an *-packr.go files
func Clean(root string) {
	root, _ = filepath.EvalSymlinks(root)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		base := filepath.Base(path)
		if base == ".git" || base == "vendor" || base == "node_modules" {
			return filepath.SkipDir
		}
		for _, suf := range []string{"-packr.go", "packrd"} {
			if strings.HasSuffix(base, suf) {
				err := os.RemoveAll(path)
				if err != nil {
					fmt.Println(err)
					return errors.WithStack(err)
				}
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		return nil
	})
}
