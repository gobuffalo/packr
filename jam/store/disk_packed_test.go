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
		"parents/homer.txt": "1f8b08000000000000fff2f0f7750d5208cecc2d28cecf03040000ffff89ec13430d000000",
		"parents/marge.txt": "1f8b08000000000000fff2750c72775508cecc2d28cecf03040000ffff0c4357460d000000",
		"kids/bart.txt":     "1f8b08000000000000ff72720c0a5108cecc2d28cecf03040000ffff905ee1f20c000000",
		"kids/lisa.txt":     "1f8b08000000000000fff2f10c765408cecc2d28cecf03040000ffffcb0086a80c000000",
		"kids/maggie.txt":   "1f8b08000000000000fff2757477f7745508cecc2d28cecf03040000ffffc2bf34da0e000000",
	})
	if err != nil {
		log.Fatal(err)
	}
	g.DefaultResolver = hgr

	packr.New("simpsons", "").DefaultResolver = hgr
}
