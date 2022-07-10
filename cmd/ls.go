package cmd

import (
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long:  ``,
}

func init() {
	cdbCmd.AddCommand(lsCmd)
}
