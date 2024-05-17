package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/GlebYaltchik/sc-keybind-extract/internal/kbd"
	"github.com/GlebYaltchik/sc-keybind-extract/internal/l10n"

	"github.com/spf13/pflag"
)

func main() {
	var (
		profileFileName string
		l10nFileName    string
	)

	pflag.StringVarP(&profileFileName, "profile", "p", "", "profile file name (usually Data/Libs/Config/defaultProfile.xml)")
	pflag.StringVarP(&l10nFileName, "localization", "l", "", "localization file name (usually Data/Localization/english/global.ini )")

	pflag.Parse()

	if profileFileName == "" || l10nFileName == "" {
		pflag.Usage()
		os.Exit(1)
	}

	pData, err := os.ReadFile(profileFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read profile: %v\n", err)
		os.Exit(1)
	}

	var v Profile

	if err := xml.Unmarshal(pData, &v); err != nil {
		fmt.Fprintf(os.Stderr, "can't unmarshal profile data: %v\n", err)
		os.Exit(1)
	}

	tr, err := l10n.NewFromFile(l10nFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	out := csv.NewWriter(os.Stdout)

	header := []string{
		"Group ID",
		"Category",
		"Label",
		"Ver",
		"Action ID",
		"Action",
		"Hotkey",
		"Mode",
		"Description",
	}

	out.Write(header)

	for _, group := range v.ActionMaps {
		for _, action := range group.Actions {
			line := append(
				[]string(nil),
				group.Name,
				tr.Translate(group.Category),
				tr.Translate(group.Label),
				group.Version,
				action.Name,
				tr.Translate(action.Label),
				kbd.Normalize(action.Keyboard),
				action.ActivationMode,
				tr.Translate(action.Description),
			)

			_ = out.Write(line)
		}
	}

	out.Flush()
}

type Profile struct {
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
