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
		profileFileName     string
		prevProfileFileName string
		l10nFileName        string
	)

	pflag.StringVarP(&profileFileName, "profile", "p", "", "profile file name (usually Data/Libs/Config/defaultProfile.xml)")
	pflag.StringVar(&prevProfileFileName, "prev-profile", "", "previous profile file name")
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

	var prev profile.Actions

	if prevProfileFileName != "" {
		prev, err = profile.DecodeFile(prevProfileFileName)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
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
		"Status",
		"Old Definition",
	}

	_ = out.Write(header)

	for _, item := range v {
		line := append(
			[]string(nil),
			item.Group.Name,
			tr.Translate(item.Group.Category),
			tr.Translate(item.Group.Label),
			item.Group.Version,
			item.Action.Name,
			tr.Translate(item.Action.Label),
			kbd.Normalize(item.Keyboard),
			item.ActivationMode,
			tr.Translate(item.Description),
		)

		line = append(line, compare(item, prev)...)

		_ = out.Write(line)
	}

	out.Flush()
}

func compare(curr profile.ActionInfo, prevData profile.Actions) []string {
	if prevData == nil {
		return nil
	}

	old, ok := prevData.Lookup(curr)
	if !ok {
		return []string{"NEW"}
	}

	curr.Action.Label = ""
	curr.Action.Description = ""
	old.Action.Label = ""
	old.Action.Description = ""

	if curr.Action == old.Action {
		return nil
	}

	if curr.Keyboard == "" && old.Keyboard != "" {
		return []string{
			"UNASSIGNED",
			fmt.Sprintf("key: %s, mode: %s", kbd.Normalize(old.Keyboard), old.ActivationMode),
		}
	}

	if curr.Keyboard != "" && old.Keyboard == "" {
		return []string{"ASSIGNED"}
	}

	return []string{
		"CHANGED",
		fmt.Sprintf("key: %s, mode: %s", kbd.Normalize(old.Keyboard), old.ActivationMode),
	}
}
