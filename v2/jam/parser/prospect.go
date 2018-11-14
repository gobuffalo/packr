package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr/v2/file/resolver"
	"github.com/gobuffalo/packr/v2/plog"
)

var DefaultIgnoredFolders = []string{".", "_", "vendor", "node_modules", "_fixtures"}

func IsProspect(path string, ignore ...string) (status bool) {
	ignore = append(ignore, DefaultIgnoredFolders...)
	// plog.Debug("parser", "IsProspect", "path", path, "ignore", ignore)
	defer func() {
		if status {
			plog.Debug("parser", "IsProspect (TRUE)", "path", path, "status", status)
		}
	}()
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		un := filepath.Base(path)
		for _, pre := range ignore {
			if strings.HasPrefix(un, pre) {
				return false
			}
		}
		return true
	}
	path = strings.ToLower(path)

	if strings.HasSuffix(path, "-packr.go") {
		return false
	}

	if strings.HasSuffix(path, "_test.go") {
		return false
	}

	ext := filepath.Ext(path)

	for i, x := range ignore {
		ignore[i] = strings.TrimSpace(strings.ToLower(x))
	}

	parts := strings.Split(resolver.OsPath(path), string(filepath.Separator))
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
