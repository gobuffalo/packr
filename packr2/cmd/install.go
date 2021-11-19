package cmd

import (
	"fmt"

	"github.com/gobuffalo/packr/v2/jam"
	"github.com/spf13/cobra"
)

const dont = `Please don't.
The following commands have been deprecated and should not be used:

* packr2 build
* packr2 install

They are, I'll be kind and say, "problematic" and cause more issues
than than the actually solve. Sorry about that. My bad.

It is recommended you use two commands instead:

$ packr2
$ go build/install
`

var installCmd = &cobra.Command{
	Use:                "install",
	Short:              "Don't. ru",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cargs := parseArgs(args)
		if globalOptions.Verbose {
			fmt.Println(dont)
		}
		if err := jam.Pack(globalOptions.PackOptions); err != nil {
			return err
		}
		return goCmd("install", cargs...)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
