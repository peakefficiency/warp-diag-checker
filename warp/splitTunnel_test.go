package warp_test

import (
	"reflect"
	"testing"

	"github.com/peakefficiency/warp-diag-checker/warp"
)

func TestSplitTunnelCIDRdetection(t *testing.T) {
	t.Parallel()
	info := warp.ParsedDiag{}

	info.Network.WarpNetIPv4 = "192.168.10.1"

}

func TestVerifyDefaultExcludedCIDRs(t *testing.T) {

	cidrsWithDefaults := append([]string{}, warp.DefaultExcludedCIDRs...)
	cidrsWithoutDefaults := []string{
		"192.0.2.0/24",
		"198.51.100.0/24",
		"203.0.113.0/24",
	}

	missing, ok := warp.VerifyDefaultExcludedCIDRs(cidrsWithDefaults)
	if !ok {
		t.Errorf("Expected true, got false. Missing CIDRs: %v", missing)
	}

	expectedMissing := []string{
		"100.64.0.0/10",
		"169.254.0.0/16",
		"192.0.0.0/24",
	}
	cidrsWithSomeDefaults := append(cidrsWithoutDefaults, warp.DefaultExcludedCIDRs[0], warp.DefaultExcludedCIDRs[3])
	cidrsWithSomeDefaults = append(cidrsWithSomeDefaults, warp.DefaultExcludedCIDRs[5:]...)
	missing, ok = warp.VerifyDefaultExcludedCIDRs(cidrsWithSomeDefaults)
	if ok {
		t.Errorf("Expected false, got true. Missing CIDRs: %v", missing)
	}
	if !reflect.DeepEqual(missing, expectedMissing) {
		t.Errorf("Expected missing CIDRs: %v, got: %v", expectedMissing, missing)
	}
}
