package kbd

import (
	"regexp"
	"strings"
)

var kbRx = regexp.MustCompile(`^([^+_\n]+)([+_]?)(?:([^+_\n]+)([+_]?))?(?:([^+_\n]+)([+_]?))?$`)

func Normalize(k string) string {
	parts := kbRx.FindStringSubmatch(k)
	if len(parts) < 1 {
		return k
	}

	for i, v := range parts[1:] {
		parts[i+1] = normalize(v)
	}

	return strings.Join(parts[1:], "")
}

func normalize(v string) string {
	if len(v) <= 1 {
		return strings.ToUpper(v)
	}

	switch v := strings.ToLower(v); v {
	case
		"lshift", "rshift",
		"lalt", "ralt",
		"lctrl", "rctrl",
		"lbracket", "rbracket",
		"mwheel":
		return strings.ToUpper(v[:2]) + v[2:]
	}

	if len(v) >= 4 && strings.ToLower(v[:3]) == "np_" {
		return strings.ToUpper(v[:4]) + v[4:]
	}

	return strings.ToUpper(v[:1]) + v[1:]
}
