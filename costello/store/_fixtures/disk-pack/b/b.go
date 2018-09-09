package b

import "github.com/gobuffalo/packr"

func init() {
	packr.New("b-box", "../c")
	packr.New("cb-box", "../c")
}
