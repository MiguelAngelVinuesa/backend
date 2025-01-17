// see https://cldr.unicode.org/index/cldr-spec/plural-rules
// see http://unicode.org/reports/tr35/tr35-numbers.html#Language_Plural_Rules

package data

import (
	"golang.org/x/text/language"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/plural"
)

// Cardinals contains the supported locales with their cardinal rules.
var Cardinals = Rules{
	{
		Locale:   language.Dutch,
		Category: plural.One,
		Rule:     "i = 1 and v = 0",
		rule:     func(in string) bool { return plural.StringI(in) == 1 && plural.StringV(in) == 0 },
	},
	{
		Locale:   language.Dutch,
		Category: plural.Other,
	},
	{
		Locale:   language.English,
		Category: plural.One,
		Rule:     "i = 1 and v = 0",
		rule:     func(in string) bool { return plural.StringI(in) == 1 && plural.StringV(in) == 0 },
	},
	{
		Locale:   language.English,
		Category: plural.Other,
	},
	{
		Locale:   language.Italian,
		Category: plural.One,
		Rule:     "i = 1 and v = 0",
		rule:     func(in string) bool { return plural.StringI(in) == 1 && plural.StringV(in) == 0 },
	},
	{
		Locale:   language.Italian,
		Category: plural.Other,
	},
}
