package utils

// ValidMultiplier determines if the given multiplier has a valid value.
// A multiplier can have at most 2 decimals and cannot be zero.
// The default multiplier is 1.0
func ValidMultiplier(multiplier float64) bool {
	return (multiplier > 0.005 && multiplier < 0.995) || multiplier > 1.005
}

// NewMultiplier determines a valid multiplier from the inputs
func NewMultiplier(in ...float64) float64 {
	out := 1.0
	for ix := range in {
		if m := in[ix]; ValidMultiplier(m) {
			out *= m
		}
	}
	return out
}

// FixArraySize calculates the best size for a new array.
func FixArraySize(in, min int) int {
	if in <= min {
		return min
	}

	t := in / min
	return (t + 1) * min
}
