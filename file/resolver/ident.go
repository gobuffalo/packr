package resolver

import (
	"path/filepath"
	"runtime"
	"strings"
)

type Ident string

func (i Ident) Key() string {
	s := string(i)
	s = strings.Replace(s, "\\", "/", -1)
	return strings.ToLower(s)
}

func (i Ident) OsPath() string {
	s := string(i)
	if runtime.GOOS == "windows" {
		s = strings.Replace(s, "/", string(filepath.Separator), -1)
	} else {
		s = strings.Replace(s, "\\", string(filepath.Separator), -1)
	}
	return filepath.FromSlash(s)
}

func (i Ident) String() string {
	return i.OsPath()
}

func (i Ident) Name() string {
	return i.OsPath()
}
