package cmd

import (
	"fmt"

	"github.com/peakefficiency/warp-diag-checker/warp"
	"github.com/peakefficiency/warp-diag-checker/wdc"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check /path/to/diag.zip ",
	Short: "Check diag for known issues",
	Long:  `The Check command attempts to search for and surface any known issues and report them in markdown`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		warp.ZipPath = args[0]
		contents, err := warp.ExtractToMemory(warp.ZipPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		wdc.GetOrLoadConfig(wdc.WdcConfig) // Make sure the config is loaded

		if !wdc.Offline {
			wdc.CheckForAppUpdate() // Check for application updates
		}

		info := contents.GetInfo(warp.ZipPath)

		warp.NewPrinter().PrintCheckResult(info.VersionCheck())

		warp.NewPrinter().PrintCheckResult(info.SplitTunnelCheck())

		warp.NewPrinter().PrintCheckResult(info.DefaultExcludeCheck())

		contents.LogSearch(info)

		warp.NewPrinter().PrintString(warp.ReportLogSearch(warp.LogSearchOutput))

	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.PersistentFlags().BoolVarP(&wdc.Plain, "plain", "p", false, "Output the report in plain markdown")
	checkCmd.PersistentFlags().BoolVarP(&wdc.Verbose, "verbose", "v", false, "Increase output verbosity")
	checkCmd.PersistentFlags().BoolVarP(&wdc.Debug, "debug", "", false, "Enable debug printing if the output is not as expected")
	checkCmd.PersistentFlags().BoolVarP(&wdc.Offline, "offline", "o", false, "Force the use of the local YAML cache file")
	// rootCmd.PersistentFlags().BoolVarP(&warp.SaveReport, "report", "r", false, "<Save the generated report in the local folder")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
