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
		fmt.Printf("%+V\n", cdb.MediaMetaToBook().Get("series"))
		strmap := cdb.MediaMetaToBook().StringMap()
		fmt.Printf("%+V\n", strmap)
		newB := urbooks.NewBookMeta(strmap)
		fmt.Printf("%+V\n", newB.StringMapToBook())
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
