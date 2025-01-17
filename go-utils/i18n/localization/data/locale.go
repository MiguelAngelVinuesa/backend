package data

import (
	"strings"

	"golang.org/x/text/language"
)

// DefaultLocale represents the default fallback locale.
const DefaultLocale = "en-GB"

// Tags returns the language tags for the given locale string.
func Tags(locale string) (language.Tag, language.Tag) {
	loc, err := language.Parse(FixLocale(locale))
	if err != nil || loc.IsRoot() {
		loc, err = language.Parse(DefaultLocale)
		if err != nil || loc.IsRoot() {
			loc = language.English
		}
	}

	base, err2 := language.Parse(loc.String()[:2])
	if err2 != nil || base.IsRoot() {
		base = loc
	}

	return base, loc
}

// FixLocale returns the fixed locale code from the input parameter.
// We only support 2-char ISO639-1, or 2-char ISO639-1 + 2-char ISO3166-1 with "-" or "_" as separator.
// Input is case-insensitive and converted to proper case on output.
// For invalid input, the function returns the default locale.
func FixLocale(locale string) string {
	switch len(locale) {
	case 2:
		return strings.ToLower(locale[:2])
	case 5:
		return strings.ToLower(locale[:2]) + "-" + strings.ToUpper(locale[3:])
	default:
		return DefaultLocale
	}
}
