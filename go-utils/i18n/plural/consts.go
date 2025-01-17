// see https://cldr.unicode.org/index/cldr-spec/plural-rules
// see http://unicode.org/reports/tr35/tr35-numbers.html#Language_Plural_Rules

package plural

// Category defines a plural rule category.
type Category uint8

// List of defined categories.
const (
	Zero Category = iota
	One
	Two
	Few
	Many
	Other
)

// String implements the Stringer interface.
func (c Category) String() string {
	switch c {
	case Zero:
		return "zero"
	case One:
		return "one"
	case Two:
		return "two"
	case Few:
		return "few"
	case Many:
		return "many"
	case Other:
		return "other"
	default:
		return "???"
	}
}
