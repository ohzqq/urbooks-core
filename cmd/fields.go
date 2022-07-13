package cmd

import (
	"fmt"

	"github.com/ohzqq/urbooks-core/urbooks"
	"github.com/spf13/cobra"
)

// fieldsCmd represents the fields command
var fieldsCmd = &cobra.Command{
	Use:   "fields",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c := urbooks.NewCalibredbCmd().
			SetLib(lib).
			Verbose(verbose).
			List("fields")
		fmt.Println(c.Run())
	},
}

func init() {
	lsCmd.AddCommand(fieldsCmd)
}
