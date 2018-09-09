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

	p := &Parser{
		Prospects: []*File{f1, f2},
	}
	r.NoError(p.Run())
}

const basicGoTmpl = `package %s

import "github.com/gobuffalo/packr"

func init() {
	packr.New("my box", "./templates")
	packr.NewBox("./old-templates")
}
`
