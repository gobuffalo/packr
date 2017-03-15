package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Clean up an *-packr.go files
func Clean(root string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		base := filepath.Base(path)
		if base == ".git" || base == "vendor" || base == "node_modules" {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(base, "-packr.go") {
			fmt.Printf("----> cleaning up %s\n", path)
			os.Remove(path)
		}
		return nil
	})
}
