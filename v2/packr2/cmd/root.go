package cmd

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/packr/v2/jam"
	"github.com/gobuffalo/packr/v2/plog"
	"github.com/spf13/cobra"
)

var globalOptions = struct {
	jam.PackOptions
	Verbose bool
	Silent  bool
}{
	PackOptions: jam.PackOptions{},
}

var rootCmd = &cobra.Command{
	Use:   "packr2",
	Short: "Packr is a simple solution for bundling static assets inside of Go binaries.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		for _, a := range args {
			if a == "--legacy" {
				globalOptions.Legacy = true
				continue
			}
			if a == "-v" || a == "--verbose" {
				globalOptions.Verbose = true
				continue
			}
		}

		// if the last argument is a .go file or directory we should
		// find boxes from there, not from the current directory.
		//	packr2 build -v cmd/main.go
		if len(args) > 0 {
			i := len(args) - 1
			dir := args[i]
			if _, err := os.Stat(dir); err == nil {
				if filepath.Ext(dir) == ".go" {
					dir = filepath.Dir(dir)
				}
				os.Chdir(dir)
				args[i] = filepath.Base(args[i])
			}
		}

		if globalOptions.Verbose {
			plog.Logger = logger.New(logger.DebugLevel)
		}
		if globalOptions.Silent {
			plog.Logger = logger.New(logger.FatalLevel)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := globalOptions.PackOptions
		roots := opts.Roots
		roots = append(roots, args...)
		opts.Roots = roots
		return jam.Pack(opts)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&globalOptions.Verbose, "verbose", "v", false, "enables verbose logging")
	rootCmd.PersistentFlags().BoolVar(&globalOptions.Legacy, "legacy", false, "uses the legacy resolution and packing system (assumes first arg || pwd for input path)")
	rootCmd.PersistentFlags().BoolVar(&globalOptions.Silent, "silent", false, "silences all output")
	rootCmd.PersistentFlags().BoolVar(&globalOptions.IgnoreImports, "ignore-imports", false, "when set to true packr won't resolve imports for boxes")
	rootCmd.PersistentFlags().StringVar(&globalOptions.StoreCmd, "store-cmd", "", "sub command to use for packing")
}
