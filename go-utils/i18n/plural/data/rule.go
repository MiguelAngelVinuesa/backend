package data

import (
	"golang.org/x/text/language"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/plural"
)

// Rule contains the details of a plural rule.
type Rule struct {
	Locale   language.Tag    `json:"locale"`
	Category plural.Category `json:"category"`
	Range    plural.Category `json:"range,omitempty"`
	Rule     string          `json:"rule,omitempty"`
	// hidden
	rule   func(in string) bool
	result plural.Category
}

// Rules is a convenience type for a slice of plural rules.
type Rules []*Rule

func (r Rules) Find(locale language.Tag, in string) *Rule {
	for ix := range r {
		rule := r[ix]
		if rule.Locale == locale && (rule.rule == nil || rule.rule(in)) {
			return rule
		}
	}
	return nil
}
