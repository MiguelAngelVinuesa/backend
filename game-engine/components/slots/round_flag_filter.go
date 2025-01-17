package slots

// OnRoundFlagValue returns a filter to test the spin round flag contains the indicated value.
func OnRoundFlagValue(flag int, value int) SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.roundFlags[flag] == value
	}
}

// OnRoundFlagValues returns a filter to test the spin round flag contains one of the indicated values.
func OnRoundFlagValues(flag int, values ...int) SpinDataFilterer {
	return func(spin *Spin) bool {
		for _, value := range values {
			if spin.roundFlags[flag] == value {
				return true
			}
		}
		return false
	}
}

// OnNotRoundFlagValue returns a filter to test the spin round flag does not contain the indicated value.
func OnNotRoundFlagValue(flag int, value int) SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.roundFlags[flag] != value
	}
}

// OnRoundFlagBelow returns a filter to test the spin round flag contains less than the indicated value.
func OnRoundFlagBelow(flag int, value int) SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.roundFlags[flag] < value
	}
}

// OnRoundFlagAbove returns a filter to test the spin round flag contains more than the indicated value.
func OnRoundFlagAbove(flag int, value int) SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.roundFlags[flag] > value
	}
}
