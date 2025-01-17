package magic

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// SpinPaylinesRange matches if the given spin has nr of paylines between min and max.
func SpinPaylinesRange(spin, min, max int) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		l := len(res[spin].Payouts)
		return l >= min && (max < min || l <= max)
	}
}

// AnyPaylinesRange matches any result with nr of paylines between min and max.
func AnyPaylinesRange(min, max int) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if l := len(res[ix].Payouts); l >= min && (max < min || l <= max) {
				return true
			}
		}
		return false
	}
}

// SpinPaylineCount matches if the given spin has the given payline with the given count.
func SpinPaylineCount(spin, payline, count int) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		payouts := res[spin].Payouts
		for ix := range payouts {
			if winline, ok := payouts[ix].(*comp.SpinPayout); ok && winline.PaylineID() == uint8(payline) && int(winline.Count()) == count {
				return true
			}
		}

		return false
	}
}

// AnyPaylineCount matches any result which has the given payline with the given count.
func AnyPaylineCount(payline, count int) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			payouts := res[ix].Payouts
			for iy := range payouts {
				if winline, ok := payouts[iy].(*comp.SpinPayout); ok && winline.PaylineID() == uint8(payline) && int(winline.Count()) == count {
					return true
				}
			}
		}
		return false
	}
}

// SpinPaylineSymbol matches if the given spin has the given payline with the given count of matching symbol.
func SpinPaylineSymbol(spin, payline int, symbol util.Index, count int) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		payouts := res[spin].Payouts
		for ix := range payouts {
			payout := payouts[ix]
			if winline, ok := payout.(*comp.SpinPayout); ok && winline.PaylineID() == uint8(payline) && symbolMatch(symbol, winline.Symbol()) && int(winline.Count()) == count {
				return true
			}
		}

		return false
	}
}

// AnyPaylineSymbol matches any result which has the given payline with the given count of matching symbol.
func AnyPaylineSymbol(payline int, symbol util.Index, count int) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			payouts := res[ix].Payouts
			for iy := range payouts {
				payout := payouts[iy]
				if winline, ok := payout.(*comp.SpinPayout); ok && winline.PaylineID() == uint8(payline) && symbolMatch(symbol, winline.Symbol()) && int(winline.Count()) == count {
					return true
				}
			}
		}
		return false
	}
}
