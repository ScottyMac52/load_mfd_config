package main

import (
	"encoding/json"
	"os"
)

type Display struct {
	Name          string  `json:"name"`
	Center        bool    `json:"center,omitempty"`
	Left          int     `json:"left,omitempty"`
	Top           int     `json:"top,omitempty"`
	Width         int     `json:"width,omitempty"`
	Height        int     `json:"height,omitempty"`
	XOffsetStart  int     `json:"xOffsetStart,omitempty"`
	XOffsetFinish int     `json:"xOffsetFinish,omitempty"`
	YOffsetStart  int     `json:"yOffsetStart,omitempty"`
	YOffsetFinish int     `json:"yOffsetFinish,omitempty"`
	Opacity       float32 `json:"opacity,omitempty"`
	Enabled       bool    `json:"enabled,omitempty"`
}

// Returns the coordinates that comprise a display area
func (d *Display) GetDimension() *Rectangle {
	rect := Rectangle{
		Left:   d.Left,
		Top:    d.Top,
		Width:  d.Width,
		Height: d.Height,
	}
	return &rect
}

func (d *Display) GetBaseDimension() *Rectangle {
	rect := Rectangle{
		Left:   0,
		Top:    0,
		Width:  d.Width,
		Height: d.Height,
	}
	return &rect
}

// SetDefaults for a single Display
func (d *Display) SetDefaults() {
	d.Opacity = 1.0
	d.Center = false
	d.Enabled = true
	d.Left = -1
	d.Top = -1
	d.Width = -1
	d.Height = -1
	d.XOffsetStart = -1
	d.YOffsetStart = -1
	d.XOffsetFinish = -1
	d.YOffsetFinish = -1
}

// Displays represents a slice of Display
type Displays []Display

type JSONFileLoader interface {
	LoadJSONFile(filename string) ([]byte, error)
	UnmarshalData(data []byte) error
}

func (ms *Displays) LoadJSONFile(filename string) ([]byte, error) {
	// Read JSON file.
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ms *Displays) UnmarshalData(data []byte) error {
	// Temporary slice to unmarshal JSON data.
	var temp []json.RawMessage

	// Unmarshal JSON into the temporary slice of raw JSON messages.
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Initialize Displays slice with the correct size.
	*ms = make(Displays, len(temp))

	// Unmarshal each item into a Display and set defaults.
	for i, itemData := range temp {
		var display Display
		display.SetDefaults() // Set defaults before unmarshalling.

		// Unmarshal the item data into the display object.
		if err := json.Unmarshal(itemData, &display); err != nil {
			return err
		}

		// Assign the display to the Displays slice.
		(*ms)[i] = display
	}

	return nil
}
