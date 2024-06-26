package warp

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type ParsedDiag struct {
	DiagName         string
	InstalledVersion string
	PlatformType     string
	Settings         ParsedSettings
	Account          ParsedAccount
	Network          ParsedNetwork
	PlatformDetails  PlatformDetails
}

type PlatformDetails struct {
	OSversion   string
	OSbuild     string
	LinuxDistro string
	LinuxKernel string
}

type ParsedAccount struct {
	AccountType  string
	DeviceID     string
	PublicKey    string
	AccountID    string
	Organization string
}

type ParsedDaemonLog struct {
	DeviceProfile string
}

type ParsedNetwork struct {
	WarpNetIPv4 string
	WarpNetIPv6 string
}

type ParsedSettings struct {
	WarpConectionStatus   bool
	SplitTunnelMode       string
	SplitTunnelList       []string
	WarpMode              string
	FallbackDomains       []string
	AlwaysOn              bool
	SwitchLocked          bool
	WiFiDisabled          bool
	EthernetDisabled      bool
	ResolveVia            string
	OnboardingDialogShown bool
	TeamsAuth             bool
	AutoFallback          bool
	CaptivePortalTimeout  int
	AllowModeSwitch       bool
	AllowUpdates          bool
	AllowLeaveOrg         bool
	ProfileID             string
	WarpTunnelProtocol    string
}

// Default values for struct fields
const (
	DefaultStringValue = "unknown"
	DefaultBoolValue   = false
	DefaultIntValue    = 0
)

// Constructor for ParsedAccount
func NewParsedAccount() ParsedAccount {
	return ParsedAccount{
		AccountType:  DefaultStringValue,
		DeviceID:     DefaultStringValue,
		PublicKey:    DefaultStringValue,
		AccountID:    DefaultStringValue,
		Organization: DefaultStringValue,
	}
}

func NewPlatformDetails() PlatformDetails {
	return PlatformDetails{
		OSbuild:     DefaultStringValue,
		OSversion:   DefaultStringValue,
		LinuxDistro: DefaultStringValue,
		LinuxKernel: DefaultStringValue,
	}
}

// Constructor for ParsedNetwork
func NewParsedNetwork() ParsedNetwork {
	return ParsedNetwork{
		WarpNetIPv4: DefaultStringValue,
		WarpNetIPv6: DefaultStringValue,
	}
}

// Constructor for ParsedSettings
func NewParsedSettings() ParsedSettings {
	return ParsedSettings{
		WarpConectionStatus:   DefaultBoolValue,
		SplitTunnelMode:       DefaultStringValue,
		SplitTunnelList:       []string{},
		WarpMode:              DefaultStringValue,
		FallbackDomains:       []string{},
		AlwaysOn:              DefaultBoolValue,
		SwitchLocked:          DefaultBoolValue,
		WiFiDisabled:          DefaultBoolValue,
		EthernetDisabled:      DefaultBoolValue,
		ResolveVia:            DefaultStringValue,
		OnboardingDialogShown: DefaultBoolValue,
		TeamsAuth:             DefaultBoolValue,
		AutoFallback:          DefaultBoolValue,
		CaptivePortalTimeout:  DefaultIntValue,
		AllowModeSwitch:       DefaultBoolValue,
		AllowUpdates:          DefaultBoolValue,
		AllowLeaveOrg:         DefaultBoolValue,
		ProfileID:             DefaultStringValue,
		WarpTunnelProtocol:    DefaultStringValue,
	}
}

// Constructor for ParsedDiag
func NewParsedDiag() ParsedDiag {
	return ParsedDiag{
		DiagName:         DefaultStringValue,
		InstalledVersion: DefaultStringValue,
		PlatformType:     DefaultStringValue,
		Settings:         NewParsedSettings(),
		Account:          NewParsedAccount(),
		Network:          NewParsedNetwork(),
	}
}

type SysInfo struct {
	LongOsVersion  string `json:"long_os_version"`
	OsVersion      string `json:"os_version"`
	DistributionID string `json:"distribution_id"`
	KernelVersion  string `json:"kernel_version"`
}

func (zipContent FileContentMap) GetInfo(zipPath string) (info ParsedDiag) {
	info = NewParsedDiag() // Initialize with default values

	info.DiagName = filepath.Base(zipPath)

	if content, ok := zipContent["platform.txt"]; ok {
		info.PlatformType = strings.ToLower(string(content.Data))
		if strings.Contains(info.PlatformType, "mac") {
			info.PlatformType = "mac" // this corrects instances where the platform type is macos causing issues on switch statements
		}
	}

	if content, ok := zipContent["warp-account.txt"]; ok {
		accountLines := strings.Split(string(content.Data), "\n")

		for _, line := range accountLines {

			if strings.Contains(line, "Account type:") {
				parts := strings.SplitN(line, ":", 3)
				if len(parts) >= 2 {
					info.Account.AccountType = strings.TrimSpace(parts[1])
				}
				continue

			}
			if strings.Contains(line, "Device ID:") {
				parts := strings.SplitN(line, ":", 3)
				if len(parts) >= 2 {
					info.Account.DeviceID = strings.TrimSpace(parts[1])
				}
				continue
			}
			if strings.Contains(line, "Public key:") {
				parts := strings.SplitN(line, ":", 3)
				if len(parts) >= 2 {
					info.Account.PublicKey = strings.TrimSpace(parts[1])
				}
				continue
			}
			if strings.Contains(line, "Account ID:") {
				parts := strings.SplitN(line, ":", 3)
				if len(parts) >= 2 {
					info.Account.AccountID = strings.TrimSpace(parts[1])

				}
				continue
			}
			if strings.Contains(line, "Organization:") {
				parts := strings.SplitN(line, ":", 3)
				if len(parts) >= 2 {
					info.Account.Organization = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	if content, ok := zipContent["warp-settings.txt"]; ok {

		settingsLines := strings.Split(string(content.Data), "\n")

		var splitTunnelStart, fallbackDomainsStart, postFallbackSettings int

		for i, line := range settingsLines {
			if strings.Contains(line, "Exclude mode") || strings.Contains(line, "Include mode") {
				splitTunnelStart = i
				info.Settings.SplitTunnelMode = line

			}
			if strings.Contains(line, "Fallback domains") {
				fallbackDomainsStart = i
			}

			if !strings.HasPrefix(line, "  ") {
				postFallbackSettings = i
			}
			// if statements above determine the sections of the settings file.
			// below actually sets the values.

			if strings.Contains(line, "Always On:") {
				if strings.Contains(line, "true") {
					info.Settings.AlwaysOn = true
					continue
				}
				info.Settings.AlwaysOn = false
				continue
			}
			if strings.Contains(line, "Switch Locked:") {
				if strings.Contains(line, "true") {
					info.Settings.SwitchLocked = true
					continue
				}
				info.Settings.SwitchLocked = false
				continue
			}
			if strings.Contains(line, "Mode:") {
				parts := strings.SplitN(line, ":", 3)
				if len(parts) >= 2 {
					info.Settings.WarpMode = strings.TrimSpace(parts[1])
				}
				continue
			}

			if strings.Contains(line, "Profile ID:") {
				parts := strings.SplitN(line, ":", 3)
				if len(parts) >= 2 {
					info.Settings.ProfileID = strings.TrimSpace(parts[1])
				}
				continue

			}
			if strings.Contains(line, "WARP tunnel protocol:") {
				parts := strings.SplitN(line, ":", 3)
				if len(parts) >= 2 {
					info.Settings.WarpTunnelProtocol = strings.TrimSpace(parts[1])
				}
				continue

			}
			if strings.Contains(line, "Disabled for Wifi:") {
				if strings.Contains(line, "true") {
					info.Settings.WiFiDisabled = true
					continue
				}
				info.Settings.WiFiDisabled = false
				continue
			}
			if strings.Contains(line, "Disabled for Ethernet:") {
				if strings.Contains(line, "true") {
					info.Settings.EthernetDisabled = true
					continue
				}
				info.Settings.EthernetDisabled = false
				continue
			}

			if strings.Contains(line, "Resolve via:") {
				info.Settings.ResolveVia = line
				continue
			}

			if strings.Contains(line, "Onboarding:") {
				if strings.Contains(line, "true") {
					info.Settings.OnboardingDialogShown = true
					continue
				}
				info.Settings.OnboardingDialogShown = false
				continue
			}
			if strings.Contains(line, "Daemon Teams Auth:") {
				if strings.Contains(line, "true") {
					info.Settings.TeamsAuth = true
					continue
				}
				info.Settings.TeamsAuth = false
				continue
			}
			if strings.Contains(line, "Disable Auto Fallback:") {
				if strings.Contains(line, "true") {
					info.Settings.AutoFallback = true
					continue
				}
				info.Settings.AutoFallback = false
				continue
			}

			if strings.Contains(line, "Allow Mode Switch:") {
				if strings.Contains(line, "true") {
					info.Settings.AllowModeSwitch = true
					continue
				}
				info.Settings.AllowModeSwitch = false
				continue
			}
			if strings.Contains(line, "Allow Updates:") {
				if strings.Contains(line, "true") {
					info.Settings.AllowUpdates = true
					continue
				}
				info.Settings.AllowUpdates = false
				continue

			}
			if strings.Contains(line, "Allowed to Leave Org:") {
				if strings.Contains(line, "true") {
					info.Settings.AllowLeaveOrg = true
					continue
				}
				info.Settings.AllowLeaveOrg = false
				continue
			}

		}

		// Check if the indices are within bounds before slicing
		if splitTunnelStart+1 < len(settingsLines) && splitTunnelStart+1 < fallbackDomainsStart {
			for _, line := range settingsLines[splitTunnelStart+1 : fallbackDomainsStart] {
				if strings.HasPrefix(line, "  ") {
					splitTunnelEntry := strings.TrimSpace(line)
					info.Settings.SplitTunnelList = append(info.Settings.SplitTunnelList, splitTunnelEntry)
				}
			}
		}

		if fallbackDomainsStart+1 < len(settingsLines) && fallbackDomainsStart+1 < postFallbackSettings {
			for _, line := range settingsLines[fallbackDomainsStart+1 : postFallbackSettings] {
				if strings.HasPrefix(line, "  ") {
					fallbackEntry := strings.TrimSpace(line)
					info.Settings.FallbackDomains = append(info.Settings.FallbackDomains, fallbackEntry)
				}
			}
		}
	}

	for _, line := range info.Settings.SplitTunnelList {

		cidr := strings.Split(line, " ")[0] // Only use the first part of the split line as the CIDR ignores comments
		Cidrs = append(Cidrs, cidr)

	}

	if content, ok := zipContent["version.txt"]; ok {
		versionContent := strings.Split(string(content.Data), "\n")
		for _, line := range versionContent {
			if strings.HasPrefix(line, "Version:") {
				// Extract the version number after the "Version:" prefix
				versionParts := strings.Split(line, " ")
				if len(versionParts) >= 2 {
					info.InstalledVersion = strings.TrimSpace(versionParts[1])
					break // Exit the loop after finding the version
				}
			} else if strings.Contains(line, ".") && !strings.Contains(line, "Commit:") {
				// This handles the case where the version line doesn't start with "Version:"
				// and is not a commit line, which is typical for Windows and Mac.
				cleanedVersion := strings.TrimSpace(line)
				if cleanedVersion != "" {
					info.InstalledVersion = cleanedVersion
					break // Exit the loop after finding the version
				}
			}
		}
	}

	if content, ok := zipContent["warp-network.txt"]; ok {

		var warpNetworkData map[string]interface{}
		err := json.Unmarshal(content.Data, &warpNetworkData)

		if err != nil {
			fmt.Println(errors.New("failed to parse warp-network.txt"))
		}

		if warpNetworkData["v4_iface"] == nil {
			info.Network.WarpNetIPv4 = ""
		} else {
			info.Network.WarpNetIPv4 = warpNetworkData["v4_iface"].(map[string]interface{})["addr"].(string)
		}
		if warpNetworkData["v6_iface"] == nil {
			info.Network.WarpNetIPv6 = ""
		} else {
			info.Network.WarpNetIPv6 = warpNetworkData["v6_iface"].(map[string]interface{})["addr"].(string)
		}

	}

	if content, ok := zipContent["sysinfo.json"]; ok {

		var SysInfo SysInfo
		switch info.PlatformType {
		case "mac":

			if err := json.Unmarshal(content.Data, &SysInfo); err == nil {
				info.PlatformDetails.OSversion = SysInfo.LongOsVersion
				info.PlatformDetails.OSbuild = SysInfo.OsVersion
			}

		case "windows":

			if err := json.Unmarshal(content.Data, &SysInfo); err == nil {
				info.PlatformDetails.OSversion = SysInfo.LongOsVersion
				info.PlatformDetails.OSbuild = SysInfo.OsVersion
			}

		case "linux":

			if err := json.Unmarshal(content.Data, &SysInfo); err == nil {
				info.PlatformDetails.LinuxDistro = SysInfo.DistributionID
				info.PlatformDetails.LinuxKernel = SysInfo.KernelVersion
			}

		}
	}
	return info
}
