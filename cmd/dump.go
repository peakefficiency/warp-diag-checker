package cmd

import (
	"fmt"

	"github.com/peakefficiency/warp-diag-checker/warp"
	"github.com/spf13/cobra"
)

var filename string

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump /path/to/diag.zip",
	Short: "dump zip contents or a specific files contents to stdout",
	Long: `dump can be used to extract data from a warp diag when it is to be passed on to other commands as with the unix pipe operator "|"
	
	wdc dump /path/to/diag.zip -f daemon.log | grep error 

	returns all lines that contain the string 'error' fromt he daemon.log file in the diag.zip file

	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		warp.ZipPath = args[0]
		contents, err := warp.ExtractToMemory(warp.ZipPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		contents.DumpFiles(filename)

	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().StringVarP(&filename, "filename", "f", "", "Specify file to dump")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dumpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
