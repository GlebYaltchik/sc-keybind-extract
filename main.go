package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/GlebYaltchik/sc-keybind-extract/internal/kbd"
	"github.com/GlebYaltchik/sc-keybind-extract/internal/l10n"
	"github.com/GlebYaltchik/sc-keybind-extract/internal/profile"

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

	v, err := profile.DecodeFile(profileFileName)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	tr, err := l10n.NewFromFile(l10nFileName)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
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

	_ = out.Write(header)

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
