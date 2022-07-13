package cmd

import (
	"fmt"

	"github.com/ohzqq/urbooks-core/urbooks"
	"github.com/spf13/cobra"
)

// libsCmd represents the libs command
var libsCmd = &cobra.Command{
	Use:   "libs",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c := urbooks.NewCalibredbCmd().
			WithLib(lib).
			Verbose(verbose).
			List("libs")
		fmt.Println(c.Run())
	},
}

func init() {
	lsCmd.AddCommand(libsCmd)
}
