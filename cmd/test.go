package cmd

import (
	"fmt"

	"github.com/ohzqq/urbooks-core/book"
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
		req = urbooks.NewRequest(lib.Name).From("books").Limit("1")
		//somecat(req)
		j := makeReq(req)
		fmt.Println(string(j))
		parsed := book.ParseBooks(j)[0]
		fmt.Printf("%+V\n", parsed.ConvertTo("markdown").Print())
		//fmt.Printf("%+V\n", parsed.GetField("narrators").String())
		//fmt.Printf("%+V\n", parsed.GetField("cover").URL())
		//fmt.Printf("%+v\n", calibredb.FieldList())
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

}
