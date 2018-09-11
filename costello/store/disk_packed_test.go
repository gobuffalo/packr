package store

import (
	"log"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/packr/file/resolver"
)

func init() {
	const gk = "__packr_global__"
	g := packr.NewBox(gk)
	hgr, err := resolver.NewHexGzip(map[string]string{
		"parents/homer.txt": "H4sIAAAAAAAC//Lw93UNUgjOzC0ozs8DBAAA//+J7BNDDQAAAA==",
		"parents/marge.txt": "H4sIAAAAAAAC//J1DHJ3VQjOzC0ozs8DBAAA//8MQ1dGDQAAAA==",
		"kids/bart.txt":     "H4sIAAAAAAAC/3JyDApRCM7MLSjOzwMEAAD//5Be4fIMAAAA",
		"kids/lisa.txt":     "H4sIAAAAAAAC//LxDHZUCM7MLSjOzwMEAAD//8sAhqgMAAAA",
		"kids/maggie.txt":   "H4sIAAAAAAAC//J1dHf3dFUIzswtKM7PAwQAAP//wr802g4AAAA=",
	})
	if err != nil {
		log.Fatal(err)
	}
	g.DefaultResolver = hgr

	packr.New("simpsons", "").DefaultResolver = hgr
}
