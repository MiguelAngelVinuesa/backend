package magic

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// SpinPayoutRange matches if the poyout for the given spin is between min and max.
func SpinPayoutRange(spin int, min, max float64) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		r := res[spin]
		return r.Total >= min && (max < min || r.Total <= max)
	}
}

// AnyPayoutRange matches if the poyout for any single spin is between min and max.
func AnyPayoutRange(min, max float64) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			r := res[ix]
			if r.Total >= min && (max < min || r.Total <= max) {
				return true
			}
		}
		return false
	}
}

// AllPayoutRange matches if the total poyout for the round is between min and max.
func AllPayoutRange(min, max float64) Matcher {
	return func(res rslt.Results) bool {
		var total float64
		for ix := range res {
			total += res[ix].Total
		}
		return total >= min && (max < min || total <= max)
	}
}

// MaxPayout matches if the last result reached the max payout.
func MaxPayout(max float64) Matcher {
	return func(res rslt.Results) bool {
		r := res[len(res)-1]
		return r.Total >= max
	}
}

// SpinScatterPayouts matches if the given spin scatter payouts count is between min and max.
func SpinScatterPayouts(spin int, min, max int) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		count := scatterPayouts(res[spin])
		return count >= min && (max < min || count <= max)
	}
}

// AnyScatterPayouts matches if any spin scatter payouts count is between min and max.
func AnyScatterPayouts(min, max int) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			count := scatterPayouts(res[ix])
			if count >= min && (max < min || count <= max) {
				return true
			}
		}
		return false
	}
}

func scatterPayouts(r *rslt.Result) int {
	var count int
	for iy := range r.Payouts {
		p := r.Payouts[iy]
		if p.Kind() == rslt.SlotScatters {
			count++
		}
	}
	return count
}

// SpinScatterPayout matches if the given spin contains a scatter payout between min and max.
func SpinScatterPayout(spin int, min, max float64) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}
		return scatterPayout(res[spin], min, max)
	}
}

// AnyScatterPayout matches if any spin contains a scatter payout between min and max.
func AnyScatterPayout(min, max float64) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if scatterPayout(res[ix], min, max) {
				return true
			}
		}
		return false
	}
}

func scatterPayout(r *rslt.Result, min, max float64) bool {
	for iy := range r.Payouts {
		p := r.Payouts[iy]
		if p.Kind() == rslt.SlotScatters {
			if p.Total() >= min && (max < min || r.Total <= max) {
				return true
			}
		}
	}
	return false
}

// SpinScatterPayoutSymbol matches if the given spin contains a scatter payout with a symbol count between min and max.
func SpinScatterPayoutSymbol(spin int, symbol util.Index, min, max int) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}
		return scatterPayoutSymbol(res[spin], symbol, uint8(min), uint8(max))
	}
}

// AnyScatterPayoutSymbol matches if any spin contains a scatter payout with a symbol count between min and max.
func AnyScatterPayoutSymbol(symbol util.Index, min, max int) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if scatterPayoutSymbol(res[ix], symbol, uint8(min), uint8(max)) {
				return true
			}
		}
		return false
	}
}

func scatterPayoutSymbol(r *rslt.Result, symbol util.Index, min, max uint8) bool {
	for iy := range r.Payouts {
		p := r.Payouts[iy]
		if p.Kind() == rslt.SlotScatters {
			if p2, ok := p.(*comp.SpinPayout); ok && symbolMatch(symbol, p2.Symbol()) {
				if p2.Count() >= min && (max < min || p2.Count() <= max) {
					return true
				}
			}
		}
	}
	return false
}

// SpinClusterPayouts matches if the given spin cluster payouts count is between min and max.
func SpinClusterPayouts(spin int, min, max int) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		count := clusterPayouts(res[spin])
		return count >= min && (max < min || count <= max)
	}
}

// AnyClusterPayouts matches if any spin cluster payouts count is between min and max.
func AnyClusterPayouts(min, max int) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			count := clusterPayouts(res[ix])
			if count >= min && (max < min || count <= max) {
				return true
			}
		}
		return false
	}
}

func clusterPayouts(r *rslt.Result) int {
	var count int
	for iy := range r.Payouts {
		p := r.Payouts[iy]
		if p.Kind() == rslt.SlotCluster {
			count++
		}
	}
	return count
}

// SpinClusterPayout matches if the given spin contains a cluster payout between min and max.
func SpinClusterPayout(spin int, min, max float64) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}
		return clusterPayout(res[spin], min, max)
	}
}

// AnyClusterPayout matches if any spin contains a cluster payout between min and max.
func AnyClusterPayout(min, max float64) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if clusterPayout(res[ix], min, max) {
				return true
			}
		}
		return false
	}
}

func clusterPayout(r *rslt.Result, min, max float64) bool {
	for iy := range r.Payouts {
		p := r.Payouts[iy]
		if p.Kind() == rslt.SlotCluster {
			if p.Total() >= min && (max < min || r.Total <= max) {
				return true
			}
		}
	}
	return false
}

// SpinClusterPayoutSymbol matches if the given spin contains a cluster payout with a symbol count between min and max.
func SpinClusterPayoutSymbol(spin int, symbol util.Index, min, max int) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}
		return clusterPayoutSymbol(res[spin], symbol, uint8(min), uint8(max))
	}
}

// AnyClusterPayoutSymbol matches if any spin contains a cluster payout with a symbol count between min and max.
func AnyClusterPayoutSymbol(symbol util.Index, min, max int) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if clusterPayoutSymbol(res[ix], symbol, uint8(min), uint8(max)) {
				return true
			}
		}
		return false
	}
}

func clusterPayoutSymbol(r *rslt.Result, symbol util.Index, min, max uint8) bool {
	for iy := range r.Payouts {
		p := r.Payouts[iy]
		if p.Kind() == rslt.SlotCluster {
			if p2, ok := p.(*comp.SpinPayout); ok && symbolMatch(symbol, p2.Symbol()) {
				if p2.Count() >= min && (max < min || p2.Count() <= max) {
					return true
				}
			}
		}
	}
	return false
}
