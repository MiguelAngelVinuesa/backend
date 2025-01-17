package slots

// OnProgressLevel returns a filter to test the progress meter contains the indicated value.
func OnProgressLevel(value int) SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.progressLevel == value
	}
}

// OnProgressLevelAbove returns a filter to test the progress meter is above the indicated value.
func OnProgressLevelAbove(value int) SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.progressLevel > value
	}
}

// OnProgressLevelBelow returns a filter to test the progress meter is below the indicated value.
func OnProgressLevelBelow(value int) SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.progressLevel < value
	}
}
