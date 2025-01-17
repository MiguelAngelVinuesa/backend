package slots

// OnZeroPayouts returns a filter to test if there are no payouts.
func OnZeroPayouts() SpinDataFilterer {
	return func(spin *Spin) bool {
		for ix := range spin.payouts {
			if spin.payouts[ix] != 0 {
				return false
			}
		}
		return true
	}
}
