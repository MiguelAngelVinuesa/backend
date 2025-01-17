package data

import (
	"golang.org/x/text/language"
)

var (
	CurrencyFormats = map[language.Tag]string{
		language.Dutch:   "%s %.2f",
		language.Italian: "%.2f %s",
		language.English: "%s%.2f",
	}

	CurrencyAppend = map[language.Tag]bool{
		language.Italian: true,
	}
)
