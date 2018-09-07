package resolver

import (
	"path/filepath"
	"strings"
)

var defaultIgnoredFolders = []string{"vendor", ".git", "node_modules", ".idea"}

func IsProspect(path string, ignore ...string) bool {
	path = strings.ToLower(path)

	if strings.HasSuffix(path, "-packr.go") {
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

	return ext == ".go" || ext == ""
}
