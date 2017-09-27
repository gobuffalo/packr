package example

import "github.com/gobuffalo/packr"

func init() {
	// packr.NewBox("../idontexists")
	foo("/templates", packr.NewBox("./templates"))
	packr.NewBox("./assets")
}

func foo(s string, box packr.Box) {}
