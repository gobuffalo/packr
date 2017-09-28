package example

import (
	"fmt"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
)

var a = packr.NewBox("./foo")

func init() {
	// packr.NewBox("../idontexists")

	b := "./baz"
	packr.NewBox(b) // won't work, no variables allowed, only strings

	foo("/templates", packr.NewBox("./templates"))
	packr.NewBox("./assets")

	r := render.New(render.Options{
		TemplatesBox: packr.NewBox("./bar"),
	})
	fmt.Println(r)
}

func foo(s string, box packr.Box) {}
