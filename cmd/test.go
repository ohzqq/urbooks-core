package cmd

import (
	"fmt"

	"github.com/ohzqq/urbooks-core/book"
	"github.com/ohzqq/urbooks-core/calibredb"
	"github.com/ohzqq/urbooks-core/urbooks"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		lib := urbooks.Lib(lib)
		fields := book.NewFields()
		fields.ParseDBFieldMeta(lib.RawPref["FieldMeta"], lib.RawPref["DisplayFields"])
		//fmt.Printf("%+v\n", fields.GetFieldIndex("narrators"))
		//fmt.Printf("%+v\n", fields.GetField("narrators").Label)
		for _, f := range fields.Each() {
			fmt.Printf("%+v\n", f.Label)
		}
		req = buildRequest(args)
		somebooks()
		fmt.Printf("%+v\n", calibredb.FieldList())
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

}
