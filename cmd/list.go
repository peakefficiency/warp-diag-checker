package cmd

import (
	"fmt"

	"github.com/peakefficiency/warp-diag-checker/warp"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list /path/to/diag.zip",
	Short: "List files within zip file",
	Long: `List files within zip file
	This can be useful to check if a file is a warp diag file`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		warp.ZipPath = args[0]
		zipContent, err := warp.ExtractToMemory(warp.ZipPath)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Files in zip:")
		for filename := range zipContent {
			fmt.Println(filename)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
