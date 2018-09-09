package q

import "github.com/gobuffalo/packr"

func init() {
	packr.New("tom", "petty")
	packr.NewBox("./heartbreakers")
}
