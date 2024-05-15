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

	v.Translate(tr)

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
				group.Category,
				group.Label,
				group.Version,
				action.Name,
				action.Label,
				action.Keyboard,
				action.ActivationMode,
				action.Description,
			)

			_ = out.Write(line)
		}
	}

	out.Flush()
}

type Profile struct {
	ActionMaps []ActionMap `xml:"actionmap"`
}

func (p *Profile) Translate(tr l10n.L10N) {
	for i := 0; i < len(p.ActionMaps); i++ {
		p.ActionMaps[i].Translate(tr)
	}
}

type ActionMap struct {
	Name     string   `xml:"name,attr"`
	Version  string   `xml:"version,attr"`
	Label    string   `xml:"UILabel,attr"`
	Category string   `xml:"UICategory,attr"`
	Actions  []Action `xml:"action"`
}

func (m *ActionMap) Translate(tr l10n.L10N) {
	m.Label = tr.Translate(m.Label)
	m.Category = tr.Translate(m.Category)

	for i := 0; i < len(m.Actions); i++ {
		m.Actions[i].Translate(tr)
	}
}

type Action struct {
	Name           string `xml:"name,attr"`
	Keyboard       string `xml:"keyboard,attr"`
	Label          string `xml:"UILabel,attr"`
	Description    string `xml:"UIDescription,attr"`
	ActivationMode string `xml:"activationMode,attr"`
}

func (a *Action) Translate(tr l10n.L10N) {
	a.Label = tr.Translate(a.Label)
	a.Description = tr.Translate(a.Description)
	a.Keyboard = kbd.Normalize(a.Keyboard)
}
