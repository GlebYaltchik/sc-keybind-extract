package profile

import (
	"encoding/xml"
	"fmt"
	"os"
)

type info struct {
	ActionMaps []actionMap `xml:"actionmap"`
}

type actionMap struct {
	Group
	Actions []Action `xml:"action"`
}

type Group struct {
	Name     string `xml:"name,attr"`
	Version  string `xml:"version,attr"`
	Label    string `xml:"UILabel,attr"`
	Category string `xml:"UICategory,attr"`
}

type Action struct {
	Name           string `xml:"name,attr"`
	Keyboard       string `xml:"keyboard,attr"`
	Label          string `xml:"UILabel,attr"`
	Description    string `xml:"UIDescription,attr"`
	ActivationMode string `xml:"activationMode,attr"`
}

type Actions []ActionInfo

type ActionInfo struct {
	Group
	Action
}

func Decode(data []byte) (Actions, error) {
	var info info

	if err := xml.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("can't unmarshal profile data: %w", err)
	}

	var res Actions

	for _, group := range info.ActionMaps {
		for _, action := range group.Actions {
			res = append(res, ActionInfo{
				Group:  group.Group,
				Action: action,
			})
		}
	}

	return res, nil
}

func DecodeFile(name string) (Actions, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("can't read profile: %w", err)
	}

	return Decode(data)
}
