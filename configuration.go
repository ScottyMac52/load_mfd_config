package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type MfdConfig struct {
	DisplayConfigurationFile string `json:"displayConfigurationFile"`
	DefaultConfiguration     string `json:"defaultConfiguration"`
	DcsSavedGamesPath        string `json:"dcsSavedGamesPath"`
	SaveCroppedImages        bool   `json:"saveCroppedImages"`
	Modules                  string `json:"modules"`
	FilePath                 string `json:"filePath"`
	UseCougar                bool   `json:"useCougar"`
	ShowRulers               bool   `json:"showRulers"`
	RulerSize                int    `json:"rulerSize"`
}

// LoadConfig loads the configuration from a JSON file.
func LoadConfiguration(filename string) *MfdConfig {
	configOnce.Do(func() {
		// Read the JSON file
		data, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		// Unmarshal JSON into the configuration struct
		var config MfdConfig
		if err := json.Unmarshal(data, &config); err != nil {
			panic(err)
		}
		fixupConfigurationPaths(&config)
		configurationInstance = &config
	})
	return configurationInstance
}

func fixupConfigurationPaths(config *MfdConfig) {
	config.FilePath = strings.ReplaceAll(os.ExpandEnv(config.FilePath), "/", "\\")
	config.DcsSavedGamesPath = strings.ReplaceAll(os.ExpandEnv(config.FilePath), "/", "\\")
	config.DisplayConfigurationFile = strings.ReplaceAll(os.ExpandEnv(config.DisplayConfigurationFile), "/", "\\")
	config.Modules = strings.ReplaceAll(os.ExpandEnv(config.Modules), "/", "\\")
}

func getCacheBaseDirectroy() string {
	return filepath.Join(getSavedGamesFolder(), "MFDMF", "Cache")
}

func clearCacheFolder() string {
	cacheFolder := getCacheBaseDirectroy()
	removeContents(cacheFolder)
	return fmt.Sprintf("The cache has been cleared at %s", cacheFolder)
}

func removeContents(path string) error {
	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	// Get a list of all files and subdirectories in the directory
	entries, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	// Loop through each entry and remove it, handling subdirectories recursively
	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			// If it's a subdirectory, call removeContents recursively
			err := removeContents(entryPath)
			if err != nil {
				return err
			}
			// Remove the empty directory
			err = os.Remove(entryPath)
			if err != nil {
				return err
			}
		} else {
			// If it's a file, remove the file
			err := os.Remove(entryPath)
			if err != nil {
				return err
			}
		}
	}

	os.Remove(path)
	return nil
}

var configurationInstance *MfdConfig
var configOnce sync.Once
