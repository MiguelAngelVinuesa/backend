package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestCurrFormats(t *testing.T) {
	testCases := []struct {
		name   string
		locale language.Tag
		fmt    string
		app    bool
	}{
		{name: "undefined", locale: language.Und},
		{name: "Dutch", locale: language.Dutch, fmt: "%s %.2f"},
		{name: "English", locale: language.English, fmt: "%s%.2f"},
		{name: "British", locale: language.BritishEnglish},
		{name: "American", locale: language.AmericanEnglish},
		{name: "Italian", locale: language.Italian, fmt: "%.2f %s", app: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt := CurrencyFormats[tc.locale]
			app := CurrencyAppend[tc.locale]

			assert.Equal(t, tc.fmt, fmt)
			assert.Equal(t, tc.app, app)
		})
	}
}
