package warp_test

import (
	"testing"

	"github.com/peakefficiency/warp-diag-checker/warp"
	"github.com/peakefficiency/warp-diag-checker/wdc"
	"github.com/stretchr/testify/assert"
)

func TestMarkownCheckResult(t *testing.T) {
	t.Parallel()

	wdc.GetOrLoadConfig(wdc.WdcConfig)
	wdc.Plain = true
	result := warp.CheckResult{

		CheckName: "Warp Version Check",
		IssueType: "OUTDATED_VERSION",
		Evidence:  "Unable to check Linux version automatically, Please verify via package repo https://pkg.cloudflareclient.com/",
	}
	got, _ := result.MarkdownCheckResult()

	want := "## Warp Version Check\nIt appears that you are not running the latest version of the chosen release train.\nPlease attempt to replicate the error using the latest available version according to the details below.\n\n- Evidence: \n\n```\nUnable to check Linux version automatically, Please verify via package repo https://pkg.cloudflareclient.com/\n```\n\n"
	assert.Equal(t, want, got, "print check result error")
}
