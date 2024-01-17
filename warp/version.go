package warp

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/go-version"
)

const (
	WindowsReleaseURL      = "https://warp-diag-checker.pages.dev/api/warp-version/windows-release"
	WindowsBetaURL         = "https://warp-diag-checker.pages.dev/api/warp-version/windows-beta"
	MacReleaseURL          = "https://warp-diag-checker.pages.dev/api/warp-version/mac-release"
	MacBetaURL             = "https://warp-diag-checker.pages.dev/api/warp-version/mac-beta"
	LinuxVersionURL        = "https://warp-diag-checker.pages.dev/api/warp-version/linux"
	WindowsDownloadURL     = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-windows-1/distribution_groups/release"
	WindowsBetaDownloadURL = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-windows/distribution_groups/beta"
	MacDownloadURL         = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-macos-1/distribution_groups/release"
	MacBetaDownloadURL     = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-macos/distribution_groups/beta"
	LinuxPKGurl            = "https://pkg.cloudflareclient.com/"
)

type Releases struct {
	ID              int       `json:"id"`
	ShortVersion    string    `json:"short_version"`
	Version         string    `json:"version"`
	UploadedAt      time.Time `json:"uploaded_at"`
	MandatoryUpdate bool      `json:"mandatory_update"`
	Enabled         bool      `json:"enabled"`
}

type LatestVersions struct {
	Release string
	Beta    string
}

func FetchVersionFrom(url string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10, // Set a timeout for the request
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version: %v", err)
	}
	defer resp.Body.Close()

	// Check if the status code is OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Convert the body to a string assuming the body is just a version string
	version := string(bodyBytes)

	return version, nil
}

func LatestWinVersions() (WinVersions LatestVersions, err error) {
	WinRelease, err := FetchVersionFrom(WindowsReleaseURL)
	if err != nil {
		return LatestVersions{}, err
	}
	WinBeta, err := FetchVersionFrom(WindowsBetaURL)
	if err != nil {
		return LatestVersions{}, err
	}

	WinVersions.Release = WinRelease
	WinVersions.Beta = WinBeta

	return WinVersions, nil
}

func LatestMacVersions() (MacVersions LatestVersions, err error) {
	MacRelease, err := FetchVersionFrom(MacReleaseURL)
	if err != nil {
		return LatestVersions{}, err
	}
	MacBeta, err := FetchVersionFrom(MacBetaURL)
	if err != nil {
		return LatestVersions{}, err
	}

	MacVersions.Release = MacRelease
	MacVersions.Beta = MacBeta

	return MacVersions, nil
}
func LatestLinuxVersion() (string, error) {
	LinuxVersion, err := FetchVersionFrom(LinuxVersionURL)
	if err != nil {
		return "", err
	}
	return LinuxVersion, nil
}

func ParseLinuxVersion(content string) (string, error) {
	// Define a regular expression pattern to match the version number
	versionRegex := regexp.MustCompile(`Version:\s*(\d+\.\d+\.\d+)`)

	// Find the first match for the version pattern
	matches := versionRegex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return "", fmt.Errorf("version string not found")
	}

	// The first match is the entire line, the second match is the captured version number
	version := matches[1]
	return version, nil
}

func (info ParsedDiag) VersionCheck() (VersionCheckResult CheckResult, err error) {

	VersionCheckResult = CheckResult{
		CheckName: "Warp Version Check",
		IssueType: "OUTDATED_VERSION",
		CheckPass: true,
	}

	if Debug {
		fmt.Printf("installed version %s", info.InstalledVersion)
	}

	switch info.PlatformType {
	case "linux":
		{
			// Fetch the latest Linux version from the API
			LinuxVersion, err := LatestLinuxVersion()
			if err != nil {
				return CheckResult{}, err
			}

			// Parse the installed version using the custom parsing logic
			parsedInstalledVersion, err := ParseLinuxVersion(info.InstalledVersion)
			if err != nil {
				return CheckResult{}, err
			}

			// Create version.Version objects for comparison
			LinuxInstalled, err := version.NewVersion(parsedInstalledVersion)
			if err != nil {
				return CheckResult{}, err
			}
			LinuxLatest, err := version.NewVersion(LinuxVersion)
			if err != nil {
				return CheckResult{}, err
			}

			// Compare the installed version with the latest version
			if LinuxInstalled.LessThan(LinuxLatest) {
				VersionCheckResult.CheckPass = false
				VersionCheckResult.Evidence = fmt.Sprintf("Installed version: %s, latest version: %s. Please update at %s", LinuxInstalled, LinuxLatest, LinuxPKGurl)
			}
		}

	case "windows":
		{
			WinVersions, err := LatestWinVersions()
			if err != nil {
				return CheckResult{}, err
			}
			WinBeta, err := version.NewVersion(WinVersions.Beta)
			if err != nil {
				return CheckResult{}, err
			}
			WinRelease, err := version.NewVersion(WinVersions.Release)
			if err != nil {
				return CheckResult{}, err
			}
			WinInstalled, err := version.NewVersion(info.InstalledVersion)

			if err != nil {
				return CheckResult{}, err
			}

			if WinInstalled.LessThan(WinRelease) {
				VersionCheckResult.CheckPass = false
				VersionCheckResult.Evidence = fmt.Sprintf("installed version: %s, Latest Release version: %s Please update at %s", WinInstalled, WinRelease, WindowsDownloadURL)
			}

			if WinInstalled.GreaterThan(WinRelease) && WinInstalled.LessThan(WinBeta) {
				VersionCheckResult.CheckPass = false
				VersionCheckResult.Evidence = fmt.Sprintf("installed version: %s, Which appears to be a beta as it is newer than the latest release: %s,  but not the latest beta which is: %s Please update at %s", WinInstalled, WinRelease, WinBeta, WindowsBetaDownloadURL)

			}

		}
	case "mac":
		{
			MacVersions, err := LatestMacVersions()
			if err != nil {
				return CheckResult{}, err
			}
			MacBeta, err := version.NewVersion(MacVersions.Beta)
			if err != nil {
				return CheckResult{}, err
			}
			MacRelease, err := version.NewVersion(MacVersions.Release)
			if err != nil {
				return CheckResult{}, err
			}
			MacInstalled, err := version.NewVersion(info.InstalledVersion)
			if err != nil {
				return CheckResult{}, err
			}

			if MacInstalled.LessThan(MacRelease) {
				VersionCheckResult.CheckPass = false
				VersionCheckResult.Evidence = fmt.Sprintf("installed version: %s, Latest Release version: %s Please update at %s", MacInstalled, MacRelease, MacReleaseURL)
			}

			if MacInstalled.GreaterThan(MacRelease) && MacInstalled.LessThan(MacBeta) {
				VersionCheckResult.CheckPass = false
				VersionCheckResult.Evidence = fmt.Sprintf("installed version: %s, Which appears to be a beta as it is newer than the latest release: %s,  but not the latest beta which is: %s Please update at %s", MacInstalled, MacRelease, MacBeta, MacBetaDownloadURL)

			}

		}
	}

	return VersionCheckResult, nil
}
