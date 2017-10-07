package example

import (
	"fmt"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
)

var a = packr.NewBox("./foo")

type S struct{}

func (S) f(packr.Box) {}

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

	s := S{}
	s.f(packr.NewBox("./sf"))
}

func foo(s string, box packr.Box) {}
