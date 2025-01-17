// see https://cldr.unicode.org/index/cldr-spec/plural-rules
// see http://unicode.org/reports/tr35/tr35-numbers.html#Language_Plural_Rules

package plural

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Int64N returns the plural operand N for an int64.
func Int64N(in int64) float64 {
	if in < 0 {
		return float64(-in)
	}
	return float64(in)
}

// FloatN returns the plural operand N for a float.
func FloatN(in float64) float64 {
	return math.Abs(in)
}

// StringN returns the plural operand N for a string.
func StringN(in string) float64 {
	f, _ := strconv.ParseFloat(in, 64)
	return math.Abs(f)
}

// Int64I returns the plural operand I for an int64.
func Int64I(in int64) int64 {
	return in
}

// FloatI returns the plural operand I for a float.
func FloatI(in float64) int64 {
	f := math.Abs(in)
	return int64(math.Floor(f))
}

// StringI returns the plural operand I for a string.
func StringI(in string) int64 {
	f, _ := strconv.ParseFloat(in, 64)
	return FloatI(f)
}

// Int64V returns the plural operand V for an int64.
func Int64V(in int64) int64 {
	return 0
}

// FloatV returns the plural operand V for a float.
func FloatV(in float64) int64 {
	return StringV(floatToString(in))
}

// StringV returns the plural operand V for a string.
func StringV(in string) int64 {
	return int64(len(fraction(in)))
}

// Int64W returns the plural operand W for an int64.
func Int64W(in int64) int64 {
	return 0
}

// FloatW returns the plural operand W for a float.
func FloatW(in float64) int64 {
	return StringW(floatToString(in))
}

// StringW returns the plural operand W for a string.
func StringW(in string) int64 {
	return int64(len(fractionNoTrailing(in)))
}

// Int64F returns the plural operand F for an int64.
func Int64F(in int64) int64 {
	return 0
}

// FloatF returns the plural operand F for a float.
func FloatF(in float64) int64 {
	return StringF(floatToString(in))
}

// StringF returns the plural operand F for a string.
func StringF(in string) int64 {
	i, _ := strconv.ParseInt(fraction(in), 10, 64)
	return i
}

// Int64T returns the plural operand T for an int64.
func Int64T(in int64) int64 {
	return 0
}

// FloatT returns the plural operand T for a float.
func FloatT(in float64) int64 {
	return StringT(fmt.Sprintf("%g", in))
}

// StringT returns the plural operand T for a string.
func StringT(in string) int64 {
	i, _ := strconv.ParseInt(fractionNoTrailing(in), 10, 64)
	return i
}

func fraction(in string) string {
	ix := strings.Index(in, ".")
	if ix >= 0 {
		return in[ix+1:]
	}
	return ""
}

func fractionNoTrailing(in string) string {
	return removeTrailingZero(fraction(in))
}

func floatToString(in float64) string {
	out := removeTrailingZero(fmt.Sprintf("%.12f", in))
	if out == "" {
		return "0"
	}
	return out
}

func removeTrailingZero(in string) string {
	l := len(in)
	for l > 0 && in[l-1] == '0' {
		l--
	}
	if l > 0 && in[l-1] == '.' {
		l--
	}
	return in[:l]
}
