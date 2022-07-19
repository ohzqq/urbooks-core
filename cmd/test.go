package cmd

import (
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
		println(lib.Name)
		//fields := book.NewFields()
		//fields.ParseDBFieldMeta(lib.RawPref["FieldMeta"], lib.RawPref["DisplayFields"])
		//fmt.Printf("%+v\n", fields.GetFieldIndex("narrators"))
		//fmt.Printf("%+v\n", fields.GetField("narrators").Label)
		//for _, f := range fields.Each() {
		//fmt.Printf("%+v\n", f.Label)
		//}
		req = urbooks.NewRequest(lib.Name).From("narrators").Limit("1")
		somecat(req)
		//fmt.Printf("%+v\n", calibredb.FieldList())
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

}
