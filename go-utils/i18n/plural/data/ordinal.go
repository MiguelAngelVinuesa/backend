package data

import (
	"math"

	"golang.org/x/text/language"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/plural"
)

// Ordinals contains the supported locales with their ordinal rules.
var Ordinals = Rules{
	{
		Locale:   language.Dutch,
		Category: plural.Other,
	},
	{
		Locale:   language.English,
		Category: plural.One,
		Rule:     "n % 10 = 1 and n % 100 != 11",
		rule: func(in string) bool {
			return math.Remainder(plural.StringN(in), 10) == 1 && math.Remainder(plural.StringN(in), 100) != 11
		},
	},
	{
		Locale:   language.English,
		Category: plural.Two,
		Rule:     "n % 10 = 2 and n % 100 != 12",
		rule: func(in string) bool {
			return math.Remainder(plural.StringN(in), 10) == 2 && math.Remainder(plural.StringN(in), 100) != 12
		},
	},
	{
		Locale:   language.English,
		Category: plural.Few,
		Rule:     "n % 10 = 3 and n % 100 != 13",
		rule: func(in string) bool {
			return math.Remainder(plural.StringN(in), 10) == 3 && math.Remainder(plural.StringN(in), 100) != 13
		},
	},
	{
		Locale:   language.English,
		Category: plural.Other,
	},
	{
		Locale:   language.Italian,
		Category: plural.Many,
		Rule:     "n = 11,8,80,800",
		rule:     func(in string) bool { n := plural.StringN(in); return n == 11 || n == 8 || n == 80 || n == 800 },
	},
	{
		Locale:   language.Italian,
		Category: plural.Other,
	},
}
