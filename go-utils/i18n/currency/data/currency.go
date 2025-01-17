package data

import (
	"strconv"
	"strings"
)

// CurrencyFromCode returns the currency details based on the given code.
// The code should be the currency code (case-insensitive), it's numeric id,
// or a country code (ISO639) for which a lookup results in a single currency for that country.
func CurrencyFromCode(code string) *Currency {
	if c := codes[strings.ToUpper(code)]; c != nil {
		return c
	}
	if i, err := strconv.Atoi(code); err == nil && i > 0 {
		if c := nums[i]; c != nil {
			return c
		}
	}

	list := iso3166[strings.ToUpper(code)]
	if len(list) == 1 {
		return list[0]
	}

	return nil
}

// CurrencyFromNum returns the currency details based on the given numeric id.
func CurrencyFromNum(num int) *Currency {
	return nums[num]
}

// CurrenciesForCountry returns the possible currencies for a country id (ISO639).
func CurrenciesForCountry(code string) []*Currency {
	return iso3166[strings.ToUpper(code)]
}

// Currency represents the details of a currency.
type Currency struct {
	num       int
	dec       int
	code      string
	name      string
	symbol    string
	countries []string
	iso3166   []string
}

// Code returns the code of the currency.
func (c *Currency) Code() string {
	return c.code
}

func (c *Currency) Num() int {
	return c.num
}

func (c *Currency) Dec() int {
	return c.dec
}

func (c *Currency) Name() string {
	return c.name
}

func (c *Currency) Symbol() string {
	if c.symbol != "" {
		return c.symbol
	}
	return c.code
}

func (c *Currency) Countries() []string {
	return c.countries
}

func (c *Currency) ISO3166() []string {
	return c.iso3166
}

var (
	codes   = make(map[string]*Currency, len(currencies))
	nums    = make(map[int]*Currency, len(currencies))
	iso3166 = make(map[string][]*Currency, len(currencies)*3/2)
)

func init() {
	for ix := range currencies {
		c := &currencies[ix]
		codes[c.code] = c
		if c.num > 0 {
			nums[c.num] = c
		}
		for iy := range c.iso3166 {
			code := c.iso3166[iy]
			if iso, ok := iso3166[code]; ok {
				iso3166[code] = append(iso, c)
			} else {
				iso3166[code] = []*Currency{c}
			}
		}
	}
}
