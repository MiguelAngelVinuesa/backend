package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrencyFromCode(t *testing.T) {
	testCases := []struct {
		name string
		code string
		fail bool
	}{
		{name: "XYZ", code: "xyz", fail: true},
		{name: "QQQ", code: "qqq", fail: true},
		{name: "EUR", code: "eur"},
		{name: "USD", code: "usd"},
		{name: "GBP", code: "GbP"},
		{name: "JPY", code: "jpY"},
		{name: "044", code: "044"},
		{name: "44", code: "44"},
		{name: "NL", code: "nl"},
		{name: "GB", code: "gB"},
		{name: "US", code: "Us"},
		{name: "JP", code: "JP"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := CurrencyFromCode(tc.code)
			if tc.fail {
				require.Nil(t, c)
			} else {
				require.NotNil(t, c)
				assert.NotEmpty(t, c.Code())
				assert.NotZero(t, c.Num())
				assert.GreaterOrEqual(t, c.Dec(), 0)
				assert.LessOrEqual(t, c.Dec(), 4)
				assert.NotEmpty(t, c.Name())
				assert.NotEmpty(t, c.Symbol())
				assert.NotEmpty(t, c.Countries())
				assert.NotEmpty(t, c.ISO3166())
			}
		})
	}
}

func TestCurrencyFromNum(t *testing.T) {
	testCases := []struct {
		name string
		num  int
		fail bool
	}{
		{name: "0", num: 0, fail: true},
		{name: "44", num: 44},
		{name: "978", num: 978},
		{name: "826", num: 826},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := CurrencyFromNum(tc.num)
			if tc.fail {
				require.Nil(t, c)
			} else {
				require.NotNil(t, c)
				assert.NotEmpty(t, c.Code())
				assert.NotZero(t, c.Num())
				assert.GreaterOrEqual(t, c.Dec(), 0)
				assert.LessOrEqual(t, c.Dec(), 4)
				assert.NotEmpty(t, c.Name())
			}
		})
	}
}

func TestCurrenciesForCountry(t *testing.T) {
	testCases := []struct {
		name  string
		code  string
		fail  bool
		count int
	}{
		{name: "YZ", code: "yz", fail: true},
		{name: "QQ", code: "qq", fail: true},
		{name: "44", code: "44", fail: true},
		{name: "NL", code: "nl", count: 1},
		{name: "GB", code: "gB", count: 1},
		{name: "US", code: "Us", count: 1},
		{name: "JP", code: "JP", count: 1},
		{name: "VE", code: "VE", count: 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := CurrenciesForCountry(tc.code)
			if tc.fail {
				require.Nil(t, c)
			} else {
				require.NotNil(t, c)
				assert.Equal(t, tc.count, len(c))
			}
		})
	}
}
