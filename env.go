package packr

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GoPath returns the current GOPATH env var
// or if it's missing, the default.
func GoPath() string {
	cmd := exec.Command("go", "env", "GOPATH")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return filepath.Join(os.Getenv("HOME"), "go")
	}
	return strings.TrimSpace(string(b))
}

// GoBin returns the current GO_BIN env var
// or if it's missing, a default of "go"
func GoBin() string {
	go_bin := os.Getenv("GO_BIN")
	if go_bin == "" {
		return "go"
	}
	return go_bin
}
