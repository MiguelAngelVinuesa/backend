package slots

func (a *ReviseAction) doRemovePayouts(spin *Spin) bool {
	if !spin.TestChance2(a.ModifyChance(a.prevPayoutChance, spin)) {
		return false
	}

	switch a.prevPayoutMech {
	case 1: // remove paylines by de-duplicate symbols on first adjoining reels.
		return spin.PreventPaylines(a.prevPayoutDir, a.prevPayoutWilds, a.genAllowDupes)
	case 2: // remove paylines by de-duplicate symbols on 2nd & 3rd reels.
		return spin.PreventPaylines2(a.prevPayoutDir, a.prevPayoutWilds, a.genAllowDupes)
	default:
		return false
	}
}
