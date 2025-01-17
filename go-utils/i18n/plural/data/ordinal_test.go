package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/localization/data"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/plural"
)

func TestFindOrdinal(t *testing.T) {
	testCases := []struct {
		name    string
		locale  string
		value   string
		wantCat plural.Category
	}{
		{
			name:    "no locale",
			locale:  "",
			wantCat: plural.Other,
		},
		{
			name:    "short locale",
			locale:  "x",
			wantCat: plural.Other,
		},
		{
			name:    "bad locale",
			locale:  "qq",
			wantCat: plural.Other,
		},
		{
			name:    "NL 0",
			locale:  "nl",
			wantCat: plural.Other,
		},
		{
			name:    "NL 1",
			locale:  "NL",
			value:   "1",
			wantCat: plural.Other,
		},
		{
			name:    "NL 2",
			locale:  "NL_NL",
			value:   "2",
			wantCat: plural.Other,
		},
		{
			name:    "NL 1.00001",
			locale:  "NL",
			value:   "1.00001",
			wantCat: plural.Other,
		},
		{
			name:    "EN 0",
			locale:  "en_gb",
			wantCat: plural.Other,
		},
		{
			name:    "EN 1",
			locale:  "EN",
			value:   "1",
			wantCat: plural.One,
		},
		{
			name:    "EN 2",
			locale:  "EN_AU",
			value:   "2",
			wantCat: plural.Two,
		},
		{
			name:    "EN 3",
			locale:  "EN_us",
			value:   "3",
			wantCat: plural.Few,
		},
		{
			name:    "EN 1.00001",
			locale:  "en",
			value:   "1.00001",
			wantCat: plural.Other,
		},
		{
			name:    "IT 0",
			locale:  "it",
			wantCat: plural.Other,
		},
		{
			name:    "IT 1",
			locale:  "IT",
			value:   "1",
			wantCat: plural.Other,
		},
		{
			name:    "IT 8",
			locale:  "IT_it",
			value:   "8",
			wantCat: plural.Many,
		},
		{
			name:    "IT 11",
			locale:  "IT_it",
			value:   "11",
			wantCat: plural.Many,
		},
		{
			name:    "IT 80",
			locale:  "IT_it",
			value:   "80",
			wantCat: plural.Many,
		},
		{
			name:    "IT 800",
			locale:  "IT_it",
			value:   "800",
			wantCat: plural.Many,
		},
		{
			name:    "IT 1.00001",
			locale:  "IT",
			value:   "1.00001",
			wantCat: plural.Other,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			base, _ := data.Tags(tc.locale)
			rule := FindOrdinal(base, tc.value)
			require.NotNil(t, rule)
			assert.Equal(t, tc.wantCat, rule.Category)
		})
	}
}
