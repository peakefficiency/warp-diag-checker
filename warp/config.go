package warp

import (
	_ "embed"
	"errors"
	"fmt"
	"io"

	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ID   string
	URL  string
	Path string
	File string
}

type WDCYaml struct {
	AppReleaseVersion  string   `yaml:"wdc_latest_version"`
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

//go:embed wdc-config.yaml
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

	yamlFile, err = os.ReadFile(c.Path)
	if err != nil {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Failed to get current user:", err)
			return
		}
		configPath := filepath.Join(usr.HomeDir, c.File)
		yamlFile, err = os.ReadFile(configPath)
		if err != nil {
			fmt.Println("Failed to read local YAML file:", err)
			yamlFile = embeddedConfig // use the embedded config as fallback
		}
	}
}

func GetOrLoadConfig(c Config) {

	if Offline {

		LocalConfig(c)
		LoadConfig(c)
		return
	}

	RemoteConfig(c)
	LoadConfig(c)

}

func RemoteConfig(c Config) {
	resp, err := http.Get(c.URL)
	if err != nil {
		fmt.Println(errors.New("unable to get remote config"))
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to download YAML file: HTTP %d\n", resp.StatusCode)

		LocalConfig(c)
	}
	yamlFile, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		LocalConfig(c)
		return
	}

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
