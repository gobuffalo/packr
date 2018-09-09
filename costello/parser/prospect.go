package parser

import (
	"os"
	"path/filepath"
	"strings"
)

var defaultIgnoredFolders = []string{"vendor", ".git", "node_modules", ".idea"}

func IsProspect(path string, ignore ...string) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.IsDir() {
			un := filepath.Base(path)
			return !strings.HasPrefix(un, "_")
		}
	}
	path = strings.ToLower(path)

	if strings.HasSuffix(path, "-packr.go") {
		return false
	}

	if strings.HasSuffix(path, "_test.go") {
		return false
	}

	ext := filepath.Ext(path)

	if len(ignore) == 0 {
		ignore = append(ignore, defaultIgnoredFolders...)
	}
	for i, x := range ignore {
		ignore[i] = strings.TrimSpace(strings.ToLower(x))
	}

	parts := strings.Split(path, string(filepath.Separator))
	if len(parts) == 0 {
		return false
	}
	for _, i := range ignore {
		for _, p := range parts {
			if p == i {
				return false
			}
		}
	}

	un := filepath.Base(path)
	if len(ext) != 0 {
		un = filepath.Base(filepath.Dir(path))
	}
	if strings.HasPrefix(un, "_") {
		return false
	}

	return ext == ".go"
}
