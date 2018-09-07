package packr

import "github.com/gobuffalo/packr/genny/packr/resolver"

type Options struct {
	Roots          []string
	IgnoredBoxes   []string
	IgnoredFolders []string
	Resolver       resolver.Resolver
	// add your stuff here
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.Resolver == nil {
		opts.Resolver = &resolver.DiskResolver{
			Roots: opts.Roots,
		}
	}
	return nil
}
