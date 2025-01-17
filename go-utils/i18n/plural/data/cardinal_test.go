package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/localization/data"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/plural"
)

func TestFindCardinal(t *testing.T) {
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
			wantCat: plural.One,
		},
		{
			name:    "NL 2",
			locale:  "nl_NL",
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
			locale:  "en_GB",
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
			wantCat: plural.Other,
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
			wantCat: plural.One,
		},
		{
			name:    "IT 2",
			locale:  "IT_it",
			value:   "2",
			wantCat: plural.Other,
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
			rule := FindCardinal(base, tc.value)
			require.NotNil(t, rule)
			assert.Equal(t, tc.wantCat, rule.Category)
		})
	}
}
