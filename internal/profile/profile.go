package profile

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Info struct {
	ActionMaps []ActionMap `xml:"actionmap"`
}

type ActionMap struct {
	Name     string   `xml:"name,attr"`
	Version  string   `xml:"version,attr"`
	Label    string   `xml:"UILabel,attr"`
	Category string   `xml:"UICategory,attr"`
	Actions  []Action `xml:"action"`
}

type Action struct {
	Name           string `xml:"name,attr"`
	Keyboard       string `xml:"keyboard,attr"`
	Label          string `xml:"UILabel,attr"`
	Description    string `xml:"UIDescription,attr"`
	ActivationMode string `xml:"activationMode,attr"`
}

func Decode(data []byte) (Info, error) {
	var info Info

	if err := xml.Unmarshal(data, &info); err != nil {
		return Info{}, fmt.Errorf("can't unmarshal profile data: %w", err)
	}

	return info, nil
}

func DecodeFile(name string) (Info, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return Info{}, fmt.Errorf("can't read profile: %w", err)
	}

	return Decode(data)
}
