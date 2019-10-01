package packr

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr/v2/file/resolver"
	"github.com/gobuffalo/packr/v2/plog"
)

func construct(name string, path string) *Box {
	var dr resolver.Resolver
	rd := resolutionDir(path)
	if len(rd) > 0 {
		dr = &resolver.Disk{Root: resolver.OsPath(rd)}
	}

	return &Box{
		Path:            path,
		Name:            name,
		ResolutionDir:   rd,
		DefaultResolver: dr,
		resolvers:       resolversMap{},
		dirs:            dirsMap{},
	}
}

func resolutionDirTestFilename(filename, og string) (string, bool) {
	ng := filepath.Join(filepath.Dir(filename), og)

	// // this little hack courtesy of the `-cover` flag!!
	cov := filepath.Join("_test", "_obj_test")
	ng = strings.Replace(ng, string(filepath.Separator)+cov, "", 1)

	if resolutionDirExists(ng, og) {
		return ng, true
	}

	ng = filepath.Join(envy.GoPath(), "src", ng)
	if resolutionDirExists(ng, og) {
		return ng, true
	}

	return og, false
}

func resolutionDirExists(s, og string) bool {
	_, err := os.Stat(s)
	if err != nil {
		return false
	}
	plog.Debug("packr", "resolutionDir", "original", og, "resolved", s)
	return true
}

func resolutionDir(og string) string {
	// packr.New
	_, filename, _, _ := runtime.Caller(3)
	ng, ok := resolutionDirTestFilename(filename, og)
	if ok {
		return ng
	}

	// packr.NewBox (deprecated)
	_, filename, _, _ = runtime.Caller(4)
	ng, ok = resolutionDirTestFilename(filename, og)
	if ok {
		return ng
	}

	return ""
}
