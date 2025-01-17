package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestFixLocale(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty", want: "en-GB"},
		{name: "short", input: "x", want: "en-GB"},
		{name: "bad", input: "xyz", want: "en-GB"},
		{name: "no sep", input: "ENGB", want: "en-GB"},
		{name: "en", input: "en", want: "en"},
		{name: "EN", input: "EN", want: "en"},
		{name: "it", input: "it", want: "it"},
		{name: "IT", input: "IT", want: "it"},
		{name: "nl", input: "nl", want: "nl"},
		{name: "NL", input: "NL", want: "nl"},
		{name: "en-gb", input: "en-gb", want: "en-GB"},
		{name: "EN_gb", input: "EN_gb", want: "en-GB"},
		{name: "en-us", input: "en-us", want: "en-US"},
		{name: "EN_us", input: "EN_us", want: "en-US"},
		{name: "it-it", input: "it-it", want: "it-IT"},
		{name: "IT_it", input: "IT_it", want: "it-IT"},
		{name: "nl-nl", input: "nl-nl", want: "nl-NL"},
		{name: "NL_nl", input: "NL_nl", want: "nl-NL"},
		{name: "nl-be", input: "nl-be", want: "nl-BE"},
		{name: "NL_be", input: "NL_be", want: "nl-BE"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := FixLocale(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTags(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		base  language.Tag
		loc   language.Tag
	}{
		{name: "empty", base: language.English, loc: language.BritishEnglish},
		{name: "short", input: "x", base: language.English, loc: language.BritishEnglish},
		{name: "bad", input: "xyz", base: language.English, loc: language.BritishEnglish},
		{name: "no sep", input: "ENGB", base: language.English, loc: language.BritishEnglish},
		{name: "en", input: "en", base: language.English, loc: language.English},
		{name: "EN", input: "EN", base: language.English, loc: language.English},
		{name: "it", input: "it", base: language.Italian, loc: language.Italian},
		{name: "IT", input: "IT", base: language.Italian, loc: language.Italian},
		{name: "nl", input: "nl", base: language.Dutch, loc: language.Dutch},
		{name: "NL", input: "NL", base: language.Dutch, loc: language.Dutch},
		{name: "en-gb", input: "en-gb", base: language.English, loc: language.BritishEnglish},
		{name: "EN_gb", input: "EN_gb", base: language.English, loc: language.BritishEnglish},
		{name: "en-us", input: "en-us", base: language.English, loc: language.AmericanEnglish},
		{name: "EN_us", input: "EN_us", base: language.English, loc: language.AmericanEnglish},
		{name: "it-it", input: "it-it", base: language.Italian, loc: language.MustParse("it-IT")},
		{name: "IT_it", input: "IT_it", base: language.Italian, loc: language.MustParse("it-IT")},
		{name: "nl-nl", input: "nl-nl", base: language.Dutch, loc: language.MustParse("nl-NL")},
		{name: "NL_nl", input: "NL_nl", base: language.Dutch, loc: language.MustParse("nl-NL")},
		{name: "nl-be", input: "nl-be", base: language.Dutch, loc: language.MustParse("nl-BE")},
		{name: "NL_be", input: "NL_be", base: language.Dutch, loc: language.MustParse("nl-BE")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			base, loc := Tags(tc.input)
			assert.Equal(t, tc.base, base)
			assert.Equal(t, tc.loc, loc)
		})
	}
}
