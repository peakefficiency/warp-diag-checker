package warp_test

import (
	"testing"

	"github.com/peakefficiency/warp-diag-checker/warp"
	"github.com/stretchr/testify/assert"
)

func TestZipToInfo(t *testing.T) {
	t.Parallel()

	realZipPath := "testdata/warp-debugging-info-20230831-185328.zip"
	content, err := warp.ExtractToMemory(realZipPath)
	if err != nil {
		t.Error("Some error extracting zip", err)
	}

	info := content.GetInfo(realZipPath)

	if info.DiagName != "warp-debugging-info-20230831-185328.zip" {
		t.Errorf("Expected DiagName to be %s, got %s", "warp-debugging-info-20230831-185328.zip", info.DiagName)
	}
	if info.PlatformType != "mac" {
		t.Errorf("Expected PlatformType to be %s, got %s", "mac", info.PlatformType)
	}
	assert.Containsf(t, info.Settings.SplitTunnelMode, "Exclude", "expected Split Tunne mode to be Exclude got %s", info.Settings.SplitTunnelMode)

	assert.Equal(t, true, info.Settings.AlwaysOn, "always on not detected correctly")

	expectedSplitTunnelIPs := []string{
		"10.0.0.0/8",
		"100.64.0.0/10",
		"169.254.0.0/16 (DHCP Unspecified)",
		"172.16.0.0/12",
		"192.0.0.0/24",
		"192.168.0.0/16",
		"224.0.0.0/24",
		"240.0.0.0/4",
		"255.255.255.255/32 (DHCP Broadcast)",
		"fe80::/10 (IPv6 Link Local)",
		"fd00::/8",
		"ff01::/16",
		"ff02::/16",
		"ff03::/16",
		"ff04::/16",
		"ff05::/16",
		"*.wikipedia.org",
		"*.en.wikipedia.org",
	}

	assert.Equal(t, expectedSplitTunnelIPs, info.Settings.SplitTunnelList, "Split tunnel list doesnt match")
	expectedFallbackDomains := []string{
		"home.arpa",
		"wikipedia.org	-> [8.8.8.8]",
		"intranet",
		"internal",
		"private",
		"localdomain",
		"domain",
		"lan",
		"home",
		"host",
		"corp",
		"local",
		"localhost",
		"invalid",
		"test",
	}
	assert.Equal(t, expectedFallbackDomains, info.Settings.FallbackDomains, "Fallback domains dont match")

	assert.Equal(t, "2023.7.159.0", info.InstalledVersion, "installed version not detected correctly")

	expectedIP4 := "192.168.20.116"

	assert.Equal(t, expectedIP4, info.Network.WarpNetIPv4, "Warp IPv4 address not detected correctly")
}
