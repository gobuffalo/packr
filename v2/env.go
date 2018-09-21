package packr

import (
	"github.com/gobuffalo/envy"
)

// GoPath returns the current GOPATH env var
// or if it's missing, the default.
func GoPath() string {
	return envy.GoPath()
}

// GoBin returns the current GO_BIN env var
// or if it's missing, a default of "go"
func GoBin() string {
	return envy.Get("GO_BIN", "go")
}
