package l10n

import (
	"bytes"
	"fmt"
	"os"
)

type L10N struct {
	dict map[string]string
}

func New(data []byte) (L10N, error) {
	const strangeSuffix = ",P"

	dict := map[string]string{}

	for n, line := range bytes.Split(data, []byte{'\n'}) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		k, v, ok := bytes.Cut(line, []byte{'='})
		if !ok {
			return L10N{}, fmt.Errorf("l10n: invalid localization string on line %d: %s", n+1, line[:min(50, len(line))])
		}

		k = bytes.TrimSuffix(k, []byte(strangeSuffix))

		dict[string(k)] = string(v)
	}

	return L10N{
		dict: dict,
	}, nil
}

func NewFromFile(fileName string) (L10N, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return L10N{}, fmt.Errorf("l10n: %w", err)
	}

	return New(data)
}

func (l L10N) Translate(v string) string {
	if len(v) == 0 || v[0] != '@' {
		return v
	}

	vv, ok := l.dict[v[1:]]
	if !ok {
		return v
	}

	return vv
}
