package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Visitor(t *testing.T) {
	r := require.New(t)
	v := NewVisitor(NewFile("example/example.go", strings.NewReader(example)))
	r.NoError(v.Run())

	r.Equal("example", v.Package)
	r.Len(v.Errors, 0)
	r.Len(v.Boxes, 7)
	r.Equal([]string{"./assets", "./bar", "./constant", "./foo", "./sf", "./templates", "./variable"}, v.Boxes)
}

const example = `package example

import (
	"github.com/gobuffalo/packr"
)

var a = packr.NewBox("./foo")

const constString = "./constant"

type S struct{}

func (S) f(packr.Box) {}

func init() {
	// packr.NewBox("../idontexists")

	b := "./variable"
	packr.NewBox(b)

	packr.NewBox(constString)

	// Cannot work from a function
	packr.NewBox(strFromFunc())

	// This variable should not be added
	fromFunc := strFromFunc()
	packr.NewBox(fromFunc)

	foo("/templates", packr.NewBox("./templates"))
	packr.NewBox("./assets")

	packr.NewBox("./bar")

	s := S{}
	s.f(packr.NewBox("./sf"))
}

func strFromFunc() string {
	return "./fromFunc"
}

func foo(s string, box packr.Box) {}
`
