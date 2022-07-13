package cmd

import (
	"github.com/ohzqq/urbooks-core/urbooks"
	"github.com/spf13/cobra"
)

var (
	dir  string
	fmts string
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if lib == "" {
			lib = urbooks.DefaultLib().Name
		}
		urbooks.NewCalibredbCmd().
			WithLib(lib).
			Verbose(verbose).
			Export(args[0], dir, fmts)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&dir, "dir", "d", "", "specify cover")
	exportCmd.Flags().StringVarP(&fmts, "formats", "f", "", "specify cover")
}
