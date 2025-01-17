package data

import (
	"golang.org/x/text/language"
)

// FindCardinal finds the plural rule for the given locale and cardinal.
func FindCardinal(locale language.Tag, in string) *Rule {
	return Cardinals.Find(locale, in)
}

// FindOrdinal finds the plural rule for the given locale and ordinal.
func FindOrdinal(locale language.Tag, in string) *Rule {
	return Ordinals.Find(locale, in)
}

// FindRange finds the cardinal plural rule for the given locale and cardinal range.
func FindRange(locale language.Tag, in1, in2 string) *Rule {
	r1 := Cardinals.Find(locale, in1)
	if r1 == nil {
		return nil
	}

	r2 := Cardinals.Find(locale, in2)
	if r2 == nil {
		return nil
	}

	for ix := range Ranges {
		rule := Ranges[ix]
		if rule.Locale != locale || rule.Category != r1.Category || rule.Range != r2.Category {
			continue
		}
		for iy := range Cardinals {
			rule2 := Cardinals[iy]
			if rule2.Locale == locale && rule2.Category == rule.result {
				return rule2
			}
		}
		break
	}

	return nil
}
