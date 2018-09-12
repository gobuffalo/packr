package a

import "github.com/gobuffalo/packr"

func init() {
	packr.New("a-box", "../c")
}
