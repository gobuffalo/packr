package example

import "github.com/gobuffalo/packr"

var a = packr.NewBox("./foo")

func init() {
	// packr.NewBox("../idontexists")

	b := "./bar"
	packr.NewBox(b) // won't work, no variables allowed, only strings

	foo("/templates", packr.NewBox("./templates"))
	packr.NewBox("./assets")
}

func foo(s string, box packr.Box) {}
