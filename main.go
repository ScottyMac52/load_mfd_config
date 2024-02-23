package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

//var displays Displays
//var modules Modules

func loadApplicationConfiguration() {
	currentUser, err := user.Current()
	if err != nil {
		logger.Log(fmt.Sprintf("Error: %s", err))
		return
	}
	configFilePath := filepath.Join(currentUser.HomeDir, "\\Saved Games\\MFDMF\\appsettings.json")
	LoadConfiguration(configFilePath)
}

func loadDisplayDefinitions() (Displays, error) {

	displayJsonPath := configurationInstance.DisplayConfigurationFile
	displays := Displays{}
	// Load JSON data.
	data, err := displays.LoadJSONFile(displayJsonPath)
	if err != nil {
		logger.Log(fmt.Sprintf("Error loading JSON file: %v\n", err))
		os.Exit(1)
	}

	// Unmarshal data into displays.
	if err := displays.UnmarshalData(data); err != nil {
		logger.Log(fmt.Sprintf("Error unmarshalling JSON data: %v\n", err))
		os.Exit(2)
	}
	return displays, nil
}

func loadModuleDefinitions(displays Displays) (Modules, error) {
	loadPath := configurationInstance.Modules
	modules, err := readModuleFiles(loadPath, &displays)
	if err != nil {
		return nil, err
	}
	return modules, nil
}

func processArguments() {
	flag.Parse()
	if verbose {
		fmt.Println("Verbose mode is enabled.")
	}
	if clearCache {
		statusMessage := clearCacheFolder()
		instance.Log(statusMessage)
		fmt.Println(statusMessage)
		return
	}
}

var logger = GetLogger()

var (
	module     string
	subModule  string
	verbose    bool
	clearCache bool
)

func init() {
	flag.StringVar(&module, "mod", "", "Module to select")
	flag.StringVar(&subModule, "sub", "", "Sub-Module to select")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose mode")
	flag.BoolVar(&clearCache, "clear", false, "Clears the cache")
}

func main() {
	logger.Log("Starting GOMFD!")

	// Parse and process the command line arguments
	processArguments()

	// load the configuration
	loadApplicationConfiguration()

	// load the display configurations
	displays, err := loadDisplayDefinitions()
	if err != nil {
		logger.Log("Unable to load display configuration")
	} else {
		displayCount := len(displays)
		logger.Log(fmt.Sprintf("Loaded %d display configurations", displayCount))
	}

	mods, err := loadModuleDefinitions(displays)
	if err != nil {
		logger.Log("Unable to load modules")
	} else {
		moduleCount := len(mods)
		logger.Log(fmt.Sprintf("Loaded %d modules", moduleCount))
	}

	// Display loaded data.
	//	fmt.Printf("Display data: %+v\n", displays)
	fmt.Printf("Modules data: %+v\n", mods)
}
