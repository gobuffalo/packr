package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:                "build",
	Short:              "Wraps the go build command with packr",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := pack(); err != nil {
			return errors.WithStack(err)
		}
		return goCmd("build", args...)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
