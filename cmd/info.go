package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/peakefficiency/warp-diag-checker/warp"
	"github.com/peakefficiency/warp-diag-checker/wdc"
	"github.com/spf13/cobra"
)

type ParsedDiagInfo struct {
	warp.ParsedDiag
}

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
		wdc.GetOrLoadConfig(wdc.WdcConfig) // Make sure the config is loaded

		if !wdc.Offline {
			wdc.CheckForAppUpdate() // Check for application updates
		}
		parsedDiag := contents.GetInfo(warp.ZipPath)
		info := ParsedDiagInfo{ParsedDiag: parsedDiag} // wrapping to allow new methods to be added to the public  version

		warp.NewPrinter().PrintString(info.PublicReportInfo())

	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.PersistentFlags().BoolVarP(&wdc.Plain, "plain", "p", false, "Output the report in plain markdown")
	infoCmd.PersistentFlags().BoolVarP(&wdc.Offline, "offline", "o", false, "Force the use of the local YAML cache file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (info ParsedDiagInfo) PublicReportInfo() (string, error) {
	var markdown strings.Builder

	markdown.WriteString("## Warp Diag Information\n")

	markdown.WriteString(fmt.Sprintf("* Name: %s\n", info.DiagName))
	markdown.WriteString(fmt.Sprintf("* Platform: %s\n", info.PlatformType))
	switch info.PlatformType {
	case "mac":
		markdown.WriteString(fmt.Sprintf("* OS version: %s\n", info.PlatformDetails.OSversion))

	case "windows":

		markdown.WriteString(fmt.Sprintf("* OS version: %s\n", info.PlatformDetails.OSversion))
		markdown.WriteString(fmt.Sprintf("* OS Build: %s\n", info.PlatformDetails.OSbuild))

	case "linux":

		markdown.WriteString(fmt.Sprintf("* Linux Distro: %s\n", info.PlatformDetails.LinuxDistro))
		markdown.WriteString(fmt.Sprintf("* Linux Kernel: %s\n", info.PlatformDetails.LinuxKernel))
	}
	markdown.WriteString(fmt.Sprintf("* Installed version: %s\n", info.InstalledVersion))
	markdown.WriteString(fmt.Sprintf("* Account ID: %s\n", info.Account.AccountID))
	markdown.WriteString(fmt.Sprintf("* Team Name: %s\n", info.Account.Organization))

	markdown.WriteString(fmt.Sprintf("* Device ID: %s\n", info.Account.DeviceID))

	if wdc.Plain {
		return markdown.String(), nil
	}

	return glamour.Render(markdown.String(), "dark")
}
