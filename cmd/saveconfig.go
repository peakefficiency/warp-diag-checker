package cmd

import (
	"fmt"
	"os"

	"github.com/peakefficiency/warp-diag-checker/wdc"
	"github.com/spf13/cobra"
)

// saveconfigCmd represents the saveconfig command
var saveconfigCmd = &cobra.Command{
	Use:   "saveconfig",
	Short: "Save 'check command' config yaml file locally for later use offline",
	Long: `Saves configuration used for the 'check' function to the users home dir as wdc-config.yaml
	This will be checked when using the '-o' offline function on the 'check' function.`,
	Run: func(cmd *cobra.Command, args []string) {

		err := wdc.SaveConfig(wdc.WdcConfig)

		if err != nil {

			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(saveconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// saveconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
