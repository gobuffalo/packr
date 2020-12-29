package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gobuffalo/packr/builder"
	"github.com/gobuffalo/packr/v2/jam/parser"
	"github.com/spf13/cobra"
)

var input string
var compress bool
var verbose bool
var includeVendored bool

var rootCmd = &cobra.Command{
	Use:   "packr",
	Short: "compiles static files into Go files",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !verbose {
			for _, a := range args {
				if a == "-v" {
					verbose = true
					break
				}
			}
		}

		if verbose {
			builder.DebugLog = func(s string, a ...interface{}) {
				os.Stdout.WriteString(fmt.Sprintf(s, a...))
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		b := builder.New(context.Background(), input)
		b.Compress = compress
		if !includeVendored {
			parser.DefaultIgnoredFolders = append(parser.DefaultIgnoredFolders, "vendor")
		}
		return b.Run()
	},
}

func init() {
	pwd, _ := os.Getwd()
	rootCmd.Flags().StringVarP(&input, "input", "i", pwd, "path to scan for packr Boxes")
	rootCmd.Flags().BoolVarP(&compress, "compress", "z", false, "compress box contents")
	rootCmd.Flags().BoolVar(&includeVendored, "vendor", false, "look for boxes in vendor directory")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "print verbose logging information")
}

// Execute the commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
