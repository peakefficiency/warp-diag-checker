package wdc

import (
	_ "embed"
	"errors"
	"fmt"
	"io"

	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/hashicorp/go-version"
	"gopkg.in/yaml.v2"
)

const AppVersion = "0.7.1"

type Config struct {
	ID   string
	URL  string
	Path string
	File string
}

type WDCYaml struct {
	WDCReleaseVersion  string   `yaml:"wdc_latest_version"`
	WDDReleaseVersion  string   `yaml:"wdd_latest_version"`
	ConfigVersion      string   `yaml:"config_version"`
	BadVersions        []string `yaml:"bad_versions"`
	LogPatternsByIssue []struct {
		SearchFile string `yaml:"search_file"`
		Issue      map[string]struct {
			SearchTerms []string `yaml:"search_term"`
		} `yaml:"issue_type"`
		ReplyType map[string]struct {
			Message string `yaml:"message"`
		} `yaml:"reply_type"`
	} `yaml:"log_patterns_by_issue"`
	ReplyByIssueType map[string]struct {
		Message string `yaml:"message"`
	} `yaml:"reply_by_issue_type"`
}

//go:embed website/public/wdc-config.yaml
var embeddedConfig []byte

var yamlFile []byte
var err error
var WdcConf WDCYaml

var WdcConfig = Config{
	ID:   "WDC",
	URL:  "https://warp-diag-checker.pages.dev/wdc-config.yaml",
	Path: "./wdc-config.yaml",
	File: "wdc-config.yaml",
}

var SaveReport, Verbose, Debug, Offline, Plain bool

func LocalConfig(c Config) {
	// First, try to read the YAML file from the specified path
	localYamlFile, err := os.ReadFile(c.Path)
	if err != nil {
		// If reading from c.Path fails, try to read from the user's home directory
		usr, err := user.Current()
		if err != nil {
			if Debug {
				fmt.Printf("Debug: Failed to get current user: %s\n", err)
			}
			fmt.Println("Using embedded configuration as fallback")
			yamlFile = embeddedConfig
			return
		}
		configPath := filepath.Join(usr.HomeDir, c.File)
		localYamlFile, err = os.ReadFile(configPath)
		if err != nil {
			if Debug {
				fmt.Printf("Debug: Failed to read local YAML file from home directory: %s\n", err)
			}
			fmt.Println("Using embedded configuration as fallback")
			yamlFile = embeddedConfig
			return
		}
	}

	// If the local file is read successfully, unmarshal it to check its version
	var localConfig WDCYaml
	if err := yaml.Unmarshal(localYamlFile, &localConfig); err != nil {
		if Debug {
			fmt.Printf("Debug: Error unmarshalling local configuration: %s\n", err)
		}
		fmt.Println("Using embedded configuration as fallback")
		yamlFile = embeddedConfig
		return
	}

	// Unmarshal the embedded configuration to check its version
	var embeddedConfigStruct WDCYaml
	if err := yaml.Unmarshal(embeddedConfig, &embeddedConfigStruct); err != nil {
		if Debug {
			fmt.Printf("Debug: Error unmarshalling embedded configuration: %s\n", err)
		}
		fmt.Println("Using embedded configuration as fallback")
		yamlFile = embeddedConfig
		return
	}

	// Compare the versions of the local and embedded configurations
	localVersion, errLocalVersion := version.NewVersion(localConfig.ConfigVersion)
	embeddedVersion, errEmbeddedVersion := version.NewVersion(embeddedConfigStruct.ConfigVersion)
	if errLocalVersion == nil && errEmbeddedVersion == nil {
		if localVersion.GreaterThanOrEqual(embeddedVersion) {
			// If the local version is newer or the same, use the local configuration
			yamlFile = localYamlFile
			return
		}
	} else {
		if Debug {
			if errLocalVersion != nil {
				fmt.Printf("Debug: Error parsing local config version: %s\n", errLocalVersion)
			}
			if errEmbeddedVersion != nil {
				fmt.Printf("Debug: Error parsing embedded config version: %s\n", errEmbeddedVersion)
			}
		}
	}

	// If there was an error reading the local file or unmarshalling the versions, use the embedded configuration
	fmt.Println("Using embedded configuration as fallback")
	yamlFile = embeddedConfig
}

func GetOrLoadConfig(c Config) {
	if Offline {
		LocalConfig(c) // Use local configuration when offline
		LoadConfig(c)  // Load the configuration into the application
		return
	}

	// Attempt to fetch and use the remote configuration
	if !RemoteConfig(c) {
		// If fetching remote configuration fails, check local configuration
		LocalConfig(c)
	}
	LoadConfig(c) // Load the configuration into the application
}

func RemoteConfig(c Config) bool {
	resp, err := http.Get(c.URL)
	if err != nil {
		fmt.Println("Unable to get remote config:", err)
		return false // Indicate failure to fetch remote config
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to download YAML file: HTTP %d\n", resp.StatusCode)
		return false // Indicate failure to fetch remote config
	}

	remoteYamlFile, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return false // Indicate failure to fetch remote config
	}

	// Unmarshal the remote configuration to check its version
	var remoteConfig WDCYaml
	if err := yaml.Unmarshal(remoteYamlFile, &remoteConfig); err != nil {
		fmt.Println("Failed to parse remote YAML file:", err)
		return false // Indicate failure to fetch remote config
	}

	// Unmarshal the embedded configuration to check its version
	var embeddedConfigStruct WDCYaml
	if err := yaml.Unmarshal(embeddedConfig, &embeddedConfigStruct); err != nil {
		fmt.Println("Failed to parse embedded YAML file:", err)
		return false // Indicate failure to fetch remote config
	}

	// Compare the versions of the remote and embedded configurations
	remoteVersion, errRemoteVersion := version.NewVersion(remoteConfig.ConfigVersion)
	embeddedVersion, errEmbeddedVersion := version.NewVersion(embeddedConfigStruct.ConfigVersion)
	if errRemoteVersion == nil && errEmbeddedVersion == nil && remoteVersion.LessThan(embeddedVersion) {
		// If the remote version is older, use the embedded configuration
		yamlFile = embeddedConfig
		return true // Indicate success, but with embedded config due to version check
	}

	// If the remote version is newer or the same, use the remote configuration
	yamlFile = remoteYamlFile
	return true // Indicate success
}

func LoadConfig(c Config) {
	var config WDCYaml
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println("Failed to parse YAML file:", err)

	}
	WdcConf = config

}

func SaveConfig(c Config) error {

	var yamlFile []byte

	resp, err := http.Get(c.URL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message := fmt.Sprintf("Failed to download YAML file: HTTP %d\n", resp.StatusCode)

		return errors.New(message)
	}
	yamlFile, err = io.ReadAll(resp.Body)
	if err != nil {
		return err

	}
	usr, err := user.Current()
	if err != nil {
		return err

	}
	configPath := filepath.Join(usr.HomeDir, c.File)
	err = os.WriteFile(configPath, yamlFile, 0600)
	if err != nil {

		return err
	}
	fmt.Println("Configuration saved to:", configPath)

	return nil
}

func CheckForAppUpdate() {
	currentVersion, err := version.NewVersion(AppVersion)
	if err != nil {
		fmt.Printf("Error parsing current version: %s\n", err)
		return
	}

	remoteVersion, err := version.NewVersion(WdcConf.WDCReleaseVersion)
	if err != nil {
		fmt.Printf("Error parsing remote version: %s\n", err)
		return
	}

	if remoteVersion.GreaterThan(currentVersion) {
		fmt.Printf("A newer version of the application is available: %s. Please update to the latest version.\n", WdcConf.WDCReleaseVersion)
		fmt.Printf("If you are not able to update at the current time. Please use the -o (--offline) flag to use the local configuration.")
		fmt.Println("You can use the saveconfig command to save the remote configuration to the local device")

		os.Exit(1)
	}
}
