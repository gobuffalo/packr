package cmd

import (
	"github.com/gobuffalo/packr/v2/jam"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:                "build",
	Short:              "Wraps the go build command with packr",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cargs []string
		for _, a := range args {
			if a == "--legacy" {
				globalOptions.Legacy = true
				continue
			}
			cargs = append(cargs, a)
		}
		if err := jam.Pack(globalOptions.PackOptions); err != nil {
			return errors.WithStack(err)
		}
		return goCmd("build", cargs...)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
