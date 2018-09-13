package store

const diskGlobalTmpl = `package {{.Package}}

import (
	"log"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/packr/file/resolver"
)

var _ = func() error {
	const gk = "{{.GK}}"
	log.Println("initializing packr global store", gk)
	g := packr.New(gk, "")
	hgr, err := resolver.NewHexGzip(map[string]string{
	{{- range $k, $v := .GlobalFiles }}
		"{{$k}}": "{{$v}}",
	{{- end }}
	})
	if err != nil {
		return err
	}
	g.DefaultResolver = hgr

	{{- range $box := .Boxes}}
	func() {
		b := packr.New("{{$box.Name}}", "{{$box.Path}}")
{{ printFiles $box}}
	}()
	{{ end }}
	return nil
}()
`

const diskImportTmpl = `package {{.Package}}

import _ "{{.Import}}"
`
