package cmd

import (
	"github.com/spf13/cobra"
)

// cdbCmd represents the cdb command
var cdbCmd = &cobra.Command{
	Use:   "cdb",
	Short: "A brief description of your command",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(cdbCmd)
}
