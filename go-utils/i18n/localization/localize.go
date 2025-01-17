package localization

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	curData "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/currency/data"
	locData "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/localization/data"
	pluData "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/plural/data"
)

// Localize localizes the given string using the given locale and parameters.
func Localize(locale, in string, params map[string]string) string {
	base, loc := locData.Tags(locale)
	p := &processor{base: base, locale: loc, params: params}
	return rxMF.ReplaceAllStringFunc(in, p.processRule)
}

// Localizer is the interface to localize multiple strings using the given parameters.
type Localizer interface {
	Localize(in string, params map[string]string) string
}

// NewLocalizer returns a Localizer interface for localizing strings using the given locale.
// The returned interface is not safe for use across multiple go-routines.
func NewLocalizer(locale string) Localizer {
	base, loc := locData.Tags(locale)
	return &processor{base: base, locale: loc}
}

type processor struct {
	base   language.Tag
	locale language.Tag
	params map[string]string
}

// Localize implements the Localizer interface.
func (p *processor) Localize(in string, params map[string]string) string {
	p.params = params
	return rxMF.ReplaceAllStringFunc(in, p.processRule)
}

func (p *processor) processRule(rule string) string {
	short := rule[1 : len(rule)-1]
	switch {
	case rxParam.MatchString(short):
		if v, ok := p.params[short]; ok {
			return v
		}

	case rxRule.MatchString(short):
		parts := rxRule.FindStringSubmatch(short)
		v1, v2, ok := p.testValues(parts[1])
		if !ok {
			return rule
		}

		switch parts[2] {
		case "select":
			return p.selector(newDefinitions(parts[3]), nil, v1, v1, v2)
		case "plural":
			return p.selector(newDefinitions(parts[3]), pluData.FindCardinal(p.base, v1), "other", v1, v2)
		case "selectordinal":
			return p.selector(newDefinitions(parts[3]), pluData.FindOrdinal(p.base, v1), "other", v1, v2)
		case "pluralrange":
			return p.selector(newDefinitions(parts[3]), pluData.FindRange(p.base, v1, v2), "other", v1, v2)
		case "number":
			return p.number(parts[3], v1)
		default:
			return rule
		}
	}

	return rule
}

func (p *processor) testValues(in string) (string, string, bool) {
	parts := strings.Split(in, "|")

	var v1, v2 string
	var ok bool

	switch len(parts) {
	case 1:
		v1, ok = p.params[parts[0]]
	case 2:
		if v1, ok = p.params[parts[0]]; ok {
			v2, ok = p.params[parts[1]]
		}
	}

	return v1, v2, ok
}

func (p *processor) selector(def definitions, rule *pluData.Rule, cat, v1, v2 string) string {
	if rule != nil {
		cat = rule.Category.String()
	}

	for ix := range def {
		d := def[ix]
		if d.selector == cat || (strings.HasPrefix(d.selector, "=") && d.selector[1:] == v1) {
			return p.processText(d.result, v1, v2)
		}
	}
	return p.processText(def[len(def)-1].result, v1, v2)
}

func (p *processor) processText(in, v1, v2 string) string {
	out := strings.Replace(in, "#", v1, 1)
	return strings.Replace(out, "#", v2, 1)
}

func (p *processor) number(spec, value string) string {
	forceInt := spec == "integer"
	forcePct := spec == "percent"
	cur := "EUR"

	if c, ok := p.params["currency"]; ok && len(c) == 3 {
		cur = c
	}

	var forceCur bool
	if rxCurrency.MatchString(spec) {
		forceCur = true
		if e := rxCurrency.FindStringSubmatch(spec); len(e) > 1 && len(e[1]) == 3 {
			cur = e[1]
		}
	}

	f, _ := strconv.ParseFloat(value, 64)
	i := int64(math.Round(f))

	prt := message.NewPrinter(p.locale)
	switch {
	case forceInt:
		return prt.Sprintf("%d", i)
	case forcePct:
		return prt.Sprintf("%.12g%%", f)
	case forceCur:
		return p.currency(prt, f, cur)
	default:
		return prt.Sprintf("%.12g", f)
	}
}

func (p *processor) currency(prt *message.Printer, amount float64, cur string) string {
	fmt := curData.CurrencyFormats[p.base]
	if fmt == "" {
		fmt = curData.CurrencyFormats[language.English]
	}

	symbol := strings.ToUpper(cur)
	if c := curData.CurrencyFromCode(symbol); c != nil {
		symbol = c.Symbol()
	}

	if curData.CurrencyAppend[p.base] {
		return prt.Sprintf(fmt, amount, symbol)
	} else {
		return prt.Sprintf(fmt, symbol, amount)
	}
}

type definition struct {
	selector string
	result   string
}

type definitions []*definition

func newDefinitions(in string) definitions {
	out := make(definitions, 0, 4)
	parts := rxSelect.FindAllStringSubmatch(in, -1)
	for ix := range parts {
		if e := parts[ix]; len(e) == 3 {
			out = append(out, &definition{selector: e[1], result: e[2]})
		}
	}
	return out
}

var (
	rxMF       = regexp.MustCompile("\\{(\\s*\\{[^{}]+}\\s*|\\s*[^{}]+\\s*)+}")
	rxParam    = regexp.MustCompile("^[^{}\\s]+$")
	rxRule     = regexp.MustCompile("\\s*([^,]+)\\s*,\\s*([^,]+)\\s*(?:,\\s*([^,]+)\\s*)?")
	rxSelect   = regexp.MustCompile("(\\S+)\\s*\\{([^}]*)}\\s*")
	rxCurrency = regexp.MustCompile("^currency(?::([A-Z]{3}))?")
)
