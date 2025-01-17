package slots

// SpinDataFilterer represents the function signature to filter actions based on spin data.
type SpinDataFilterer func(spin *Spin) bool

// OnFirstSpin filters based on the spin being the first spin.
func OnFirstSpin(spin *Spin) bool {
	return spin.kind == RegularSpin || spin.kind == FirstSpin
}

// OnSecondSpin filters based on the spin being the second spin.
func OnSecondSpin(spin *Spin) bool {
	return spin.kind == SecondSpin
}

// OnFreeSpin filters based on the spin being a free spin.
func OnFreeSpin(spin *Spin) bool {
	return spin.kind == FreeSpin || spin.kind == FirstFreeSpin || spin.kind == SecondFreeSpin
}

// OnRefillSpin filters based on the spin being a refill spin.
func OnRefillSpin(spin *Spin) bool {
	return spin.kind == RefillSpin
}

// OnSuperSpin filters based on the spin being a super spin.
func OnSuperSpin(spin *Spin) bool {
	return spin.kind == SuperSpin
}

// OnFreeSpins filters based on having free spins remaining after the current spin.
func OnFreeSpins(spin *Spin) bool {
	return spin.freeSpins > 0
}

// OnNoFreeSpins filters based on not having any free spins remaining after the current spin.
func OnNoFreeSpins(spin *Spin) bool {
	return spin.freeSpins == 0
}

// OnRemainingFreeSpins filters based on how many free spins are remaining after the current spin.
func OnRemainingFreeSpins(remain uint64) SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.freeSpins == remain
	}
}

// OnBetweenFreeSpins filters based on how many free spins are remaining after the current spin.
func OnBetweenFreeSpins(from, until uint64) SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.freeSpins >= from && spin.freeSpins <= until
	}
}

// OnStickies filters based on having sticky symbols in the grid.
func OnStickies(spin *Spin) bool {
	return spin.HasSticky()
}

// OnNoStickies filters based on not having sticky symbols in the grid.
func OnNoStickies(spin *Spin) bool {
	return !spin.HasSticky()
}

// OnSpinSequence filters based on the spin sequence number equal to the given value.
func OnSpinSequence(seq uint64) SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.spinSeq == seq
	}
}

// OnSpinSequenceBelow filters based on the spin sequence number below the given value.
func OnSpinSequenceBelow(seq uint64) SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.spinSeq < seq
	}
}

// OnSpinSequenceAbove filters based on the spin sequence number above the given value.
func OnSpinSequenceAbove(seq uint64) SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.spinSeq > seq
	}
}
