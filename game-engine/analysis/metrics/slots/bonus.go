package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// BonusKind indicates the kind of bonus/free spins a game round has received.
type BonusKind uint8

const (
	NoFreeSpins BonusKind = iota

	BOTFreeSpins
	CCBSuperX
	CCBFreeSpins
	MGDWildRespin
	MGDBonus
	MGDSuperBonus
	LAMNorth
	LAMSouth
	LAMNorthBB
	LAMSouthBB
	OWLBonus1
	OWLBonus2
	OWLBonus3
	OWLBonus1BB
	OWLBonus2BB
	OWLBonus3BB
	FRMRefill
	FRMFreeSpins
	FRMFreeSpinsBB
	OFGFreeSpins
	OFGFreeSpinsBB
	FPRWildRespin
	FPRBonus
	FPRSuperBonus
	MGDBonus3
	MGDBonus4
	MGDSuperBonus3
	MGDSuperBonus4
	FPRBonus3
	FPRBonus4
	FPRSuperBonus3
	FPRSuperBonus4
	BTRFreeSpins
	BTRFreeSpins3
	BTRFreeSpins4
	BTRFreeSpins5
	BTRFreeSpins6
	BTRFreeSpins33
	OFGLevel1
	OFGLevel2
	OFGLevel3
	OFGLevel4
	OFGLevel1BB
	OFGLevel2BB
	OFGLevel3BB
	OFGLevel4BB

	// UnknownFreeSpins must always be last entry!
	UnknownFreeSpins
)

var names = []string{
	"none",
	"BOT free spins",
	"CCB SuperX spins",
	"CCB free spins",
	"MGD wild respin",
	"MGD bonus spins (all)",
	"MGD super bonus spins (all)",
	"LAM free spins North",
	"LAM free spins South",
	"LAM free spins North - bonus buy",
	"LAM free spins South - bonus buy",
	"OWL bonus spins (1)",
	"OWL bonus spins (2)",
	"OWL bonus spins (3)",
	"OWL bonus spins (1) - bonus bet",
	"OWL bonus spins (2) - bonus bet",
	"OWL bonus spins (3) - bonus bet",
	"FRM refill spins",
	"FRM free spins",
	"FRM free spins - bonus buy",
	"OFG free spins",
	"OFG free spins - bonus buy",
	"FPR wild respin",
	"FPR bonus spins (all)",
	"FPR super bonus spins (all)",
	"MGD bonus spins (3 scatters)",
	"MGD bonus spins (4 scatters)",
	"MGD super bonus spins (3 scatters)",
	"MGD super bonus spins (4 scatters)",
	"FPR bonus spins (3 scatters)",
	"FPR bonus spins (4 scatters)",
	"FPR super bonus spins (3 scatters)",
	"FPR super bonus spins (4 scatters)",
	"BTR free spins",
	"BTR free spins (3 scatters)",
	"BTR free spins (4 scatters)",
	"BTR free spins (5 scatters)",
	"BTR free spins (6 scatters)",
	"BTR free spins (3+3 scatters)",
	"OFG free spins (level 1)",
	"OFG free spins (level 2)",
	"OFG free spins (level 3)",
	"OFG free spins (level 4)",
	"OFG free spins (level 1) - bonus buy",
	"OFG free spins (level 2) - bonus buy",
	"OFG free spins (level 3) - bonus buy",
	"OFG free spins (level 4) - bonus buy",
	"Free spins",
}

// String implements the Stringer interface.
func (b BonusKind) String() string {
	if int(b) >= len(names) {
		return "??? free spins"
	}
	return names[b]
}

// DetermineBonusKind determines the kind of bonus (if any) from the game nr and the set of results.
func DetermineBonusKind(gameNR tg.GameNR, res results.Results) BonusKind {
	l := len(res)
	if l <= 1 {
		return NoFreeSpins
	}

	s, ok := res[0].Data.(*slots.SpinResult)

	switch gameNR {
	case tg.BOTnr:
		return BOTFreeSpins

	case tg.CCBnr:
		switch {
		case ok && s.Kind() == slots.SuperSpin:
			return CCBSuperX
		default:
			return CCBFreeSpins
		}

	case tg.MGDnr:
		wr := IsMgdWildRespin(s.Initial())
		sc := SymbolCount(s.Initial(), 9)
		switch {
		case wr && sc < 3:
			return MGDWildRespin

		case wr:
			if sc == 4 {
				return MGDSuperBonus4
			}
			return MGDSuperBonus3

		default:
			if sc == 4 {
				return MGDBonus4
			}
			return MGDBonus3
		}

	case tg.LAMnr:
		switch LamWing(res) {
		case "north":
			if IsBonusBuy(res) {
				return LAMNorthBB
			}
			return LAMNorth

		case "south":
			if IsBonusBuy(res) {
				return LAMSouthBB
			}
			return LAMSouth
		}

	case tg.FPRnr:
		wr := IsMgdWildRespin(s.Initial())
		sc := SymbolCount(s.Initial(), 9)
		switch {
		case wr && sc < 3:
			return FPRWildRespin

		case wr:
			if sc == 4 {
				return FPRSuperBonus4
			}
			return FPRSuperBonus3

		default:
			if sc == 4 {
				return FPRBonus4
			}
			return FPRBonus3
		}

	case tg.OFGnr:
		return ofgBonus(res)

	case tg.OWLnr:
		// TODO: ...

	case tg.FRMnr:
		// TODO: ...

	case tg.BTRnr, tg.BERnr:
		return btrBonus(s, res)

	case tg.BBSnr:
		// TODO: ...

	case tg.HOGnr:
		// TODO: ...

	case tg.MOGnr:
		// TODO: ...
	}

	return UnknownFreeSpins
}

// IsMgdWildRespin returns true if the spin grid matches the MGD wild respin bonus.
func IsMgdWildRespin(initial utils.Indexes) bool {
	if len(initial) < 22 {
		return false
	}
	return initial[0] == 10 && initial[1] == 10 && initial[20] == 10 && initial[21] == 10
}

// LamWing returns the player choice of wing for a bonus game.
func LamWing(res results.Results) string {
	if len(res) < 2 {
		return ""
	}
	if s, ok := res[1].Data.(*slots.SpinResult); ok {
		return s.Choices()["wing"]
	}
	return ""
}

// IsBonusBuy returns true if the game round was started with the bonus buy/bonus bet feature.
func IsBonusBuy(res results.Results) bool {
	for ix := 0; ix < 5 && ix < len(res); ix++ {
		if s, ok := res[ix].Data.(*slots.SpinResult); ok {
			if bb, _ := s.BonusBuy(); bb != 0 {
				return true
			}
		}
	}
	return false
}

func btrBonus(s *slots.SpinResult, res results.Results) BonusKind {
	sc := SymbolCount(s.Initial(), 11)
	switch sc {
	case 3:
		if btrBonus33(res) {
			return BTRFreeSpins33
		}
		return BTRFreeSpins3
	case 4:
		return BTRFreeSpins4
	case 5:
		return BTRFreeSpins5
	case 6:
		return BTRFreeSpins6
	}

	return UnknownFreeSpins
}

func ofgBonus(res results.Results) BonusKind {
	last := len(res) - 1
	if s, ok := res[last].Data.(*slots.SpinResult); ok {
		if IsBonusBuy(res) {
			switch s.ProgressLevel() {
			case 1, 2, 3:
				return OFGLevel1BB
			case 4, 5, 6, 7:
				return OFGLevel2BB
			case 8, 9, 10, 11, 12:
				return OFGLevel3BB
			default:
				return OFGLevel4BB
			}
		} else {
			switch s.ProgressLevel() {
			case 1, 2, 3:
				return OFGLevel1
			case 4, 5, 6, 7:
				return OFGLevel2
			case 8, 9, 10, 11, 12:
				return OFGLevel3
			default:
				return OFGLevel4
			}
		}
	}

	if IsBonusBuy(res) {
		return OFGFreeSpinsBB
	}
	return OFGFreeSpins
}

func btrBonus33(res results.Results) bool {
	for ix := range res {
		if ix > 0 {
			r := res[ix].Data.(*slots.SpinResult)
			if SymbolCount(r.Initial(), 11) == 3 {
				return true
			}
		}
	}
	return false
}

func SymbolCount(initial utils.Indexes, symbol utils.Index) int {
	var count int
	for _, id := range initial {
		if id == symbol {
			count++
		}
	}
	return count
}
