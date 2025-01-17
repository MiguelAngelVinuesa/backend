package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/localization/data"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/plural"
)

func TestFindRange(t *testing.T) {
	testCases := []struct {
		name    string
		locale  string
		value1  string
		value2  string
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
			name:    "NL 0 1",
			locale:  "nl",
			value2:  "1",
			wantCat: plural.One,
		},
		{
			name:    "NL 0 2",
			locale:  "nl_be",
			value2:  "2",
			wantCat: plural.Other,
		},
		{
			name:    "NL 1 2",
			locale:  "NL_BE",
			value1:  "1",
			value2:  "2",
			wantCat: plural.Other,
		},
		{
			name:    "EN 0 1",
			locale:  "en",
			value2:  "1",
			wantCat: plural.Other,
		},
		{
			name:    "EN 0 2",
			locale:  "en_ca",
			value2:  "2",
			wantCat: plural.Other,
		},
		{
			name:    "EN 1 2",
			locale:  "EN_CA",
			value1:  "1",
			value2:  "2",
			wantCat: plural.Other,
		},
		{
			name:    "IT 0 1",
			locale:  "it",
			value2:  "1",
			wantCat: plural.One,
		},
		{
			name:    "IT 0 2",
			locale:  "IT",
			value2:  "2",
			wantCat: plural.Other,
		},
		{
			name:    "IT 1 2",
			locale:  "IT_it",
			value1:  "1",
			value2:  "2",
			wantCat: plural.Other,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			base, _ := data.Tags(tc.locale)
			rule := FindRange(base, tc.value1, tc.value2)
			require.NotNil(t, rule)
			assert.Equal(t, tc.wantCat, rule.Category)
		})
	}
}
