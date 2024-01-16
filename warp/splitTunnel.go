package warp

import (
	"fmt"
	"net"
	"strings"
)

var DefaultExcludedCIDRs = []string{
	"10.0.0.0/8",
	"100.64.0.0/10",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.0.0.0/24",
	"192.168.0.0/16",
	"224.0.0.0/24",
	"240.0.0.0/4",
	"255.255.255.255/32",
	"fe80::/10",
	"fd00::/8",
	"ff01::/16",
	"ff02::/16",
	"ff03::/16",
	"ff04::/16",
	"ff05::/16",
}
var Cidrs []string

func (info ParsedDiag) SplitTunnelCheck() (CheckResult, error) {
	SplitTunnelResult := CheckResult{
		CheckName: "IP Address Split Tunnel Check",
		IssueType: "SPLITTUNNEL",
	}

	// Check if the SplitTunnelMode or WarpNetIPv4 have default values
	if info.Settings.SplitTunnelMode == DefaultStringValue || info.Network.WarpNetIPv4 == "" {
		SplitTunnelResult.CheckPass = false
		SplitTunnelResult.IssueType = "default values found"
		SplitTunnelResult.Evidence = "Could not set SplitTunnelMode or WarpNetIPv4 to valid values while parsing \n manually validate warp-settings.txt warp-network.txt files."
		return SplitTunnelResult, nil
	}

	ip := net.ParseIP(info.Network.WarpNetIPv4)
	isInCIDR := false
	var matchedCIDR string
	for _, cidr := range Cidrs { //this is extracted as part of GetInfo to ensure order of split tunnel related checks is irrelevant
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(ip) {
			isInCIDR = true
			matchedCIDR = cidr
			break
		}
	}

	mode := info.Settings.SplitTunnelMode
	if (strings.Contains(mode, "Exclude mode") && isInCIDR) || (strings.Contains(mode, "Include mode") && !isInCIDR) {
		SplitTunnelResult.CheckPass = true
	} else {
		SplitTunnelResult.CheckPass = false
	}

	if !SplitTunnelResult.CheckPass {
		SplitTunnelResult.Evidence = fmt.Sprintf("Mode: %s\nAssigned IP: %s, Not Matched in Split tunnel CIDRS ", mode, info.Network.WarpNetIPv4)
	} else {
		SplitTunnelResult.Evidence = fmt.Sprintf("Mode: %s\nAssigned IP: %s, Matched in Split tunnel CIDRS: %s", mode, info.Network.WarpNetIPv4, matchedCIDR)
	}

	return SplitTunnelResult, nil
}

func (info ParsedDiag) DefaultExcludeCheck() (CheckResult, error) {

	DefaultExcludeResult := CheckResult{

		CheckName: "Default Exclude Check",
		IssueType: "EXCLUDE_EDITED",
	}

	if strings.Contains(info.Settings.SplitTunnelMode, "Exclude mode") {
		missingCIDRs, allDefaultCIDRsPresent := VerifyDefaultExcludedCIDRs(Cidrs)
		if !allDefaultCIDRsPresent {

			DefaultExcludeResult.IssueType = "EXCLUDE_EDITED"

			DefaultExcludeResult.CheckPass = false
			missingCIDRStr := strings.Join(missingCIDRs, ", ")
			DefaultExcludeResult.Evidence += fmt.Sprintf("Missing default excluded CIDRs: %s", missingCIDRStr)
		} else {
			DefaultExcludeResult.CheckPass = true
			DefaultExcludeResult.Evidence += "All default excluded CIDRs are present"
		}
	}

	return DefaultExcludeResult, nil
}

func VerifyDefaultExcludedCIDRs(cidrs []string) ([]string, bool) {
	missingCIDRs := make([]string, 0)

	for _, defaultCIDR := range DefaultExcludedCIDRs {
		found := false
		for _, cidr := range cidrs {
			if cidr == defaultCIDR {
				found = true
				break
			}
		}
		if !found {
			missingCIDRs = append(missingCIDRs, defaultCIDR)
		}
	}

	return missingCIDRs, len(missingCIDRs) == 0
}
