package cmd

import (
	"fmt"

	"github.com/ohzqq/urbooks-core/urbooks"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if lib == "" {
			lib = urbooks.DefaultLib().Name
		}
		cdb := urbooks.NewCalibredbCmd().SetUser(calibreUser).SetLib(lib).Add(args[0])
		fmt.Printf("%+v\n", cdb)
		fmt.Printf("%+V\n", cdb.MediaMetaToBook())
		fmt.Printf("%+V\n", cdb.MediaMetaToBook().StringMap())
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
