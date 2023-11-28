package warp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-version"
)

const (
	MacReleaseURL          = "https://install.appcenter.ms/api/v0.1/apps/cloudflare/1.1.1.1-macos-1/distribution_groups/release/public_releases?scope=tester"
	MacBetaURL             = "https://install.appcenter.ms/api/v0.1/apps/cloudflare/1.1.1.1-macos/distribution_groups/beta/public_releases?scope=tester"
	WindowsReleaseURL      = "https://install.appcenter.ms/api/v0.1/apps/cloudflare/1.1.1.1-windows-1/distribution_groups/release/public_releases?scope=tester"
	WindowsBetaURL         = "https://install.appcenter.ms/api/v0.1/apps/cloudflare/1.1.1.1-windows/distribution_groups/beta/public_releases?scope=tester"
	LinuxPKGurl            = "https://pkg.cloudflareclient.com/"
	WindowsDownloadURL     = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-windows-1/distribution_groups/release"
	WindowsBetaDownloadURL = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-windows/distribution_groups/beta"
	MacDownloadURL         = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-macos-1/distribution_groups/release"
	MacBetaDownloadURL     = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-macos/distribution_groups/beta"
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

func FetchReleasesFrom(url string) (ReleaseDetails []Releases, err error) {

	client := &http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Releases{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1 Safari/605.1.15")

	resp, err := client.Do(req)
	if err != nil {
		return []Releases{}, fmt.Errorf("failed to fetch latest version: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Releases{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Releases{}, fmt.Errorf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(bodyBytes, &ReleaseDetails)
	if err != nil {
		return []Releases{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	if len(ReleaseDetails) == 0 {
		return []Releases{}, fmt.Errorf("no releases found")
	}

	return ReleaseDetails, nil

}

func LatestWinVersions() (WinVersions LatestVersions, err error) {

	WinBetaReleases, err := FetchReleasesFrom(WindowsBetaURL)
	if err != nil {
		return LatestVersions{}, err
	}
	WinReleases, err := FetchReleasesFrom(WindowsReleaseURL)
	if err != nil {
		return LatestVersions{}, err
	}

	WinVersions.Release = WinReleases[0].Version
	WinVersions.Beta = WinBetaReleases[0].Version

	return WinVersions, nil

}

func LatestMacVersions() (MacVersions LatestVersions, err error) {

	MacBetaReleases, err := FetchReleasesFrom(MacBetaURL)
	if err != nil {
		return LatestVersions{}, err
	}
	MacReleases, err := FetchReleasesFrom(MacReleaseURL)
	if err != nil {
		return LatestVersions{}, err
	}

	MacVersions.Release = MacReleases[0].ShortVersion
	MacVersions.Beta = MacBetaReleases[0].ShortVersion

	return MacVersions, nil

}

func (info ParsedDiag) VersionCheck() (VersionCheckResult CheckResult, err error) {

	VersionCheckResult = CheckResult{
		CheckName: "Warp Version Check",
		IssueType: "OUTDATED_VERSION",
		CheckPass: true,
	}

	if Debug {
		fmt.Println(info.InstalledVersion)
		fmt.Println("debug")
	}

	switch info.PlatformType {
	case "linux":
		{
			VersionCheckResult.Evidence = fmt.Sprintf("Unable to check Linux version automatically, Please verify via package repo %s", LinuxPKGurl)
			VersionCheckResult.CheckPass = false
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
				VersionCheckResult.Evidence = fmt.Sprintf("installed version: %s, Latest Release version: %s", WinInstalled, WinRelease)
			}

			if WinInstalled.GreaterThan(WinRelease) && WinInstalled.LessThan(WinBeta) {
				VersionCheckResult.CheckPass = false
				VersionCheckResult.Evidence = fmt.Sprintf("installed version: %s, Which appears to be a beta as it is newer than the latest release: %s,  but not the latest beta which is: %s", WinInstalled, WinRelease, WinBeta)

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
				VersionCheckResult.Evidence = fmt.Sprintf("installed version: %s, Latest Release version: %s", MacInstalled, MacRelease)
			}

			if MacInstalled.GreaterThan(MacRelease) && MacInstalled.LessThan(MacBeta) {
				VersionCheckResult.CheckPass = false
				VersionCheckResult.Evidence = fmt.Sprintf("installed version: %s, Which appears to be a beta as it is newer than the latest release: %s,  but not the latest beta which is: %s", MacInstalled, MacRelease, MacBeta)

			}

		}
	}

	return VersionCheckResult, nil
}
