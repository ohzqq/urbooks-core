package cmd

import (
	"github.com/ohzqq/urbooks-core/urbooks"
	"github.com/spf13/cobra"
)

var cover string
var iniFile string

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import books to calibre",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if lib == "" {
			lib = urbooks.DefaultLib().Name
		}
		if cover == "" {
			cover = urbooks.FindCover()
		}
		urbooks.NewCalibredbCmd().
			WithLib(lib).
			Verbose(verbose).
			Import(args[0], cover, iniFile)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&cover, "cover", "c", "", "specify cover")
	importCmd.Flags().StringVarP(&iniFile, "ini", "i", "", "use a metadata file instead of embedded")
}
