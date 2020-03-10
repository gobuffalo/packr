package internal

import (
	"os"
	"runtime"
	"strings"
)

// Mods returns true when go modules supports is enabled
func Mods() bool {
	go111 := os.Getenv("GO111MODULE")

	if !inGoPath() {
		return go111 != "off"
	}

	return go111 == "on"
}

func inGoPath() bool {
	pwd, _ := os.Getwd()
	for _, p := range GoPaths() {
		if strings.HasPrefix(pwd, p) {
			return true
		}
	}
	return false
}

// GoPaths return the defined gopath list.
func GoPaths() []string {
	gp := os.Getenv("GOPATH")
	if runtime.GOOS == "windows" {
		return strings.Split(gp, ";") // Windows uses a different separator
	}
	return strings.Split(gp, ":")
}
