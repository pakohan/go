package main

import "strings"

func removeUnderscores(s string) string {
	return strings.ReplaceAll(s, "_", "")
}

func plural(s string) string {
	if s == "" {
		return s
	}

	if s[len(s)-1] == 'y' && len(s) > 1 {
		switch s[len(s)-2] {
		case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		// do nothing, since 'y' only gets replaced by 'ie' for plurals in the English language,
		// if the preceding character is not a vowel
		default:
			s = s[:len(s)-1] + "ie"
		}
	}

	return s + "s"
}
