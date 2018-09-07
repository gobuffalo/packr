package resolver

import "github.com/gobuffalo/genny"

type Resolver interface {
	Resolve() error
	Prospects() []genny.File
	Boxes() map[string][]genny.File
}
