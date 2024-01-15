package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/peakefficiency/warp-diag-checker/warp"
	"github.com/spf13/cobra"
)

type ParsedDiagInfo warp.ParsedDiag

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info /path/to/diag.zip",
	Short: "Return key details from the supplied warp diag",
	Long:  `Return key details from the supplied warp diag `,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		warp.ZipPath = args[0]
		contents, err := warp.ExtractToMemory(warp.ZipPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		warp.GetOrLoadConfig(warp.WdcConfig) // Make sure the config is loaded

		if !warp.Offline {
			warp.CheckForAppUpdate()             // Check for application updates
		}
		info := contents.GetInfo(warp.ZipPath)

		warp.NewPrinter().PrintString(info.ReportInfo())

	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.PersistentFlags().BoolVarP(&warp.Plain, "plain", "p", false, "Output the report in plain markdown")
	infoCmd.PersistentFlags().BoolVarP(&warp.Offline, "offline", "o", false, "Force the use of the local YAML cache file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (info ParsedDiagInfo) ReportInfo() (string, error) {
	var markdown strings.Builder

	markdown.WriteString("## Warp Diag Information\n")

	markdown.WriteString(fmt.Sprintf("* Name: %s\n", info.DiagName))
	markdown.WriteString(fmt.Sprintf("* Platform: %s\n", info.PlatformType))

	if warp.Plain {
		return markdown.String(), nil
	}

	return glamour.Render(markdown.String(), "dark")
}
