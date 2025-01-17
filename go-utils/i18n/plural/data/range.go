package data

import (
	"golang.org/x/text/language"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/plural"
)

// Ranges contains the supported locales with their cardinal range rules.
var Ranges = Rules{
	{
		Locale:   language.Dutch,
		Category: plural.One,
		Range:    plural.Other,
		Rule:     "one + other -> other",
		result:   plural.Other,
	},
	{
		Locale:   language.Dutch,
		Category: plural.Other,
		Range:    plural.One,
		Rule:     "other + one -> one",
		result:   plural.One,
	},
	{
		Locale:   language.Dutch,
		Category: plural.Other,
		Range:    plural.Other,
		Rule:     "other + other -> other",
		result:   plural.Other,
	},
	{
		Locale:   language.English,
		Category: plural.One,
		Range:    plural.Other,
		Rule:     "one + other -> other",
		result:   plural.Other,
	},
	{
		Locale:   language.English,
		Category: plural.Other,
		Range:    plural.One,
		Rule:     "other + one -> other",
		result:   plural.Other,
	},
	{
		Locale:   language.English,
		Category: plural.Other,
		Range:    plural.Other,
		Rule:     "other + other -> other",
		result:   plural.Other,
	},
	{
		Locale:   language.Italian,
		Category: plural.One,
		Range:    plural.Other,
		Rule:     "one + other -> other",
		result:   plural.Other,
	},
	{
		Locale:   language.Italian,
		Category: plural.Other,
		Range:    plural.One,
		Rule:     "other + one -> one",
		result:   plural.One,
	},
	{
		Locale:   language.Italian,
		Category: plural.Other,
		Range:    plural.Other,
		Rule:     "other + other -> other",
		result:   plural.Other,
	},
}
