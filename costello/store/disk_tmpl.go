package store

const diskGlobalTmpl = `package {{.Package}}

import (
	"log"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/packr/file/resolver"
)

func init() {
	const gk = "__packr_global__"
	g := packr.NewBox(gk)
	hgr, err := resolver.NewHexGzip(map[string]string{
	{{- range $k, $v := .GlobalFiles }}
		"{{$k}}": "{{$v}}",
	{{- end }}
	})
	if err != nil {
		log.Fatal(err)
	}
	g.DefaultResolver = hgr

	{{- range $box := .Boxes}}
	func() {
		b := packr.New("{{$box.Name}}", "{{$box.Path}}")
{{ printFiles $box}}
	}()
	{{ end }}
}
`

const diskImportTmpl = `package {{.Package}}

import _ "{{.Import}}"
`
