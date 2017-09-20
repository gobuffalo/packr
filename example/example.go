package example

import "github.com/gobuffalo/packr"

func init() {
	// packr.NewBox("../idontexists")
	packr.NewBox("./assets")
}
