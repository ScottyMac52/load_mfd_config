package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Stores the coordinates that are used to copy out portions of images for cropping
type Offsets struct {
	XOffsetStart  int `json:"xOffsetStart,omitempty"`
	XOffsetFinish int `json:"xOffsetFinish,omitempty"`
	YOffsetStart  int `json:"yOffsetStart,omitempty"`
	YOffsetFinish int `json:"yOffsetFinish,omitempty"`
}

// Stores base image properties
type ImageProperties struct {
	Center            bool    `json:"center,omitempty"`
	Opacity           float32 `json:"opacity,omitempty"`
	Enabled           bool    `json:"enabled,omitempty"`
	UseAsSwitch       bool    `json:"useAsSwitch,omitempty"`
	NeedsThrottleType bool    `json:"needsThrottleType,omitempty"`
	Image             *image.RGBA
}

// Stores a Configuration
type Configuration struct {
	Name           string `json:"name"`
	FileName       string `json:"fileName"`
	Module         *Module
	Parent         *Configuration
	Display        *Display
	Opacity        float32         `json:"opacity,omitempty"`
	Center         bool            `json:"center,omitempty"`
	Enabled        bool            `json:"enabled,omitempty"`
	Left           int             `json:"left,omitempty"`
	Top            int             `json:"top,omitempty"`
	Width          int             `json:"width,omitempty"`
	Height         int             `json:"height,omitempty"`
	XOffsetStart   int             `json:"xOffsetStart,omitempty"`
	XOffsetFinish  int             `json:"xOffsetFinish,omitempty"`
	YOffsetStart   int             `json:"yOffsetStart,omitempty"`
	YOffsetFinish  int             `json:"yOffsetFinish,omitempty"`
	Configurations []Configuration `json:"subConfigDef"`
}

// Stores a Module
type Module struct {
	Name           string          `json:"name"`
	Tag            string          `json:"tag"`
	DisplayName    string          `json:"displayName"`
	FileName       string          `json:"fileName"`
	Category       string          `json:"category"`
	Configurations []Configuration `json:"configurations"`
}

// Slice of Modules
type Modules []Module

// JSON structure of a module file
type JSONModuleData struct {
	Modules Modules `json:"modules"`
}

// Interface that establishes the contract for loading a JSON file and unmarshalling the data into objects or slices
type JSONModuleLoader interface {
	LoadJSONFile(filename string) ([]byte, error)
	UnmarshalData(data []byte) error
}

// Interface contract that defines how Configuration data is processed
type ConfigurationImpl interface {
	GetOffset() Offsets
	GetDisplayRef(displays Displays) (*Display, error)
	SetDefaults(display *Display)
	GetDimension() *Rectangle
	GetBaseDimension() *Rectangle
	SetFileName(module *Module)
	IsInside(config *Configuration) bool
	CenterIn(config *Configuration) Rectangle
}

// Determines if a Configuration is inside another Configuration
func (outer *Configuration) IsInside(inner *Configuration) bool {
	if inner.Width > outer.Width || inner.Height > outer.Height {
		return false
	} else {
		return true
	}
}

func (outer *Configuration) CenterIn(inner *Configuration) (*Rectangle, error) {
	// Check to see that the inner is smaller than the outer
	isInside := outer.IsInside(inner)
	if !isInside {
		return nil, errors.New("inner Configuration is bigger than the outer Configuration")
	}
	// Calculate the center position of the outer Configuration
	outerCenterX := outer.Left + (outer.Width / 2)
	outerCenterY := outer.Top + (outer.Height / 2)

	// Calculate the position to center the inner Configuration
	innerLeft := outerCenterX - (inner.Width / 2)
	innerTop := outerCenterY - (inner.Height / 2)

	// Create and return a Rectangle instance
	return &Rectangle{Left: innerLeft, Top: innerTop, Width: inner.Width, Height: inner.Height}, nil
}

func (config *Configuration) GetDisplayRef(displays Displays) (*Display, error) {
	if displays == nil {
		return nil, nil
	} else {
		matched := false
		for i := range displays {
			currentDisplay := displays[i]
			if strings.HasPrefix(config.Name, currentDisplay.Name) {
				matched = true
			}
			if matched {
				return &currentDisplay, nil
			}
		}
	}
	return nil, nil
}

// SetDefaults for a single Configuration
func (config *Configuration) SetDefaults(display *Display) {
	if display != nil {
		config.Opacity = display.Opacity
		config.Enabled = display.Enabled
		config.Center = display.Center
		config.Left = display.Left
		config.Top = display.Top
		config.Width = display.Width
		config.Height = display.Height
		config.XOffsetStart = display.XOffsetStart
		config.YOffsetStart = display.XOffsetFinish
		config.XOffsetFinish = display.YOffsetStart
		config.YOffsetFinish = display.YOffsetFinish
	} else {
		config.Opacity = 1.0
		config.Enabled = true
		config.Center = false
		config.Left = -1
		config.Top = -1
		config.Width = -1
		config.Height = -1
		config.XOffsetStart = -1
		config.YOffsetStart = -1
		config.XOffsetFinish = -1
		config.YOffsetFinish = -1
	}
}

// Returns the coordinates that comprise a Configuration area
func (d *Configuration) GetDimension() *Rectangle {
	rect := Rectangle{
		Left:   d.Left,
		Top:    d.Top,
		Width:  d.Width,
		Height: d.Height,
	}
	return &rect
}

func (d *Configuration) GetBaseDimension() *Rectangle {
	rect := Rectangle{
		Left:   0,
		Top:    0,
		Width:  d.Width,
		Height: d.Height,
	}
	return &rect
}

func (config *Configuration) GetOffset() (Offsets, error) {
	return Offsets{XOffsetStart: config.XOffsetStart, XOffsetFinish: config.XOffsetFinish, YOffsetStart: config.YOffsetStart, YOffsetFinish: config.YOffsetFinish}, nil
}

// Reads all of the modules from the specified path and below
func readModuleFiles(startingPath string, displays *Displays) (Modules, error) {
	var modules Modules

	// Walk the directory tree starting from the specified path
	err := filepath.Walk(startingPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a JSON file
		if filepath.Ext(filePath) == ".json" {
			// Read the JSON file
			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			// Unmarshal the JSON data into a wrapper structure with the "Modules" array
			jsonData := JSONModuleData{}
			err = json.Unmarshal(data, &jsonData)
			if err != nil {
				return err
			}

			// Set the Category for each module and process configurations recursively
			for i := range jsonData.Modules {
				currentModule := &jsonData.Modules[i]
				// Calculate the relative Category based on the starting path
				dir, _ := path.Split(strings.ReplaceAll(filePath, "\\", "/"))
				relativePath, err := filepath.Rel(getBaseDirectory(), dir)
				if err != nil {
					return err
				}
				currentModule.Category = relativePath
				processConfigurationsRecursively(currentModule, nil, currentModule.Configurations, displays)
				//unmarshalConfigurations(currentModule.Configurations, displays)
			}

			// Append the modules from the wrapper to the main modules slice
			modules = append(modules, jsonData.Modules...)
			return nil
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return modules, nil
}

func (currentConfig *Configuration) SetFileName(module *Module) error {
	if len(currentConfig.FileName) > 0 {
		isInCorrectPath, err := path.Match(configurationInstance.FilePath, currentConfig.FileName)
		if err != nil {
			return err
		}
		if !isInCorrectPath {
			tempPath := path.Join(configurationInstance.FilePath, currentConfig.FileName)
			currentConfig.FileName = strings.ReplaceAll(os.ExpandEnv(tempPath), "/", "\\")
		}
	} else {
		if module != nil && len(module.FileName) > 0 {
			isInCorrectPath, err := path.Match(configurationInstance.FilePath, module.FileName)
			if err != nil {
				return err
			}
			if !isInCorrectPath {
				tempPath := path.Join(configurationInstance.FilePath, module.FileName)
				currentConfig.FileName = strings.ReplaceAll(os.ExpandEnv(tempPath), "/", "\\")
			}
		}
	}
	return nil
}

// Recursively process configurations
func processConfigurationsRecursively(module *Module, parent *Configuration, configs []Configuration, displays *Displays) error {
	for i := range configs {
		currentConfig := &configs[i]
		if module != nil {
			currentConfig.Module = module
		}
		if parent != nil {
			currentConfig.Parent = parent
		}
		logger.Log(fmt.Sprintf("Getting Display for %s\n", currentConfig.Name))
		displayRef, _ := currentConfig.GetDisplayRef(*displays)
		currentConfig.SetDefaults(displayRef)
		currentConfig.Display = displayRef
		currentConfig.SetFileName(module)
		err := unmarshalConfigurations(currentConfig.Configurations, displays)
		if err != nil {
			return err
		}
		processConfigurationsRecursively(nil, currentConfig, currentConfig.Configurations, displays)
	}
	return nil
}

// Unmarshal configurations
func unmarshalConfigurations(configs []Configuration, displays *Displays) error {
	for i := range configs {
		currentConfig := &configs[i]
		configData, err := json.Marshal(currentConfig)
		if err != nil {
			return err
		}
		displayRef, err := currentConfig.GetDisplayRef(*displays)
		if err != nil {
			return err
		}

		currentConfig.SetDefaults(displayRef)
		err = json.Unmarshal(configData, &currentConfig)
		if err != nil {
			return err
		}
		currentConfig.Display = displayRef
		currentConfig.SetFileName(nil)
		unmarshalConfigurations(currentConfig.Configurations, displays)
	}
	return nil
}
