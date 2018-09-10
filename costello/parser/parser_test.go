package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parser_Run(t *testing.T) {
	r := require.New(t)

	f1 := NewFile("a/a.x", strings.NewReader(fmt.Sprintf(basicGoTmpl, "a")))
	f2 := NewFile("b/b.x", strings.NewReader(fmt.Sprintf(basicGoTmpl, "b")))

	p := New(f1, f2)
	boxes, err := p.Run()
	r.NoError(err)

	r.Len(boxes, 4)
}

func Test_NewFrom_Roots(t *testing.T) {
	r := require.New(t)
	p, err := NewFromRoots([]string{"./_fixtures/new_from_roots"})
	r.NoError(err)

	boxes, err := p.Run()
	r.NoError(err)

	r.Len(boxes, 3)
}

const basicGoTmpl = `package %s

import "github.com/gobuffalo/packr"

func init() {
	packr.New("elvis", "./presley")
	packr.NewBox("./buddy-holly")
}
`
