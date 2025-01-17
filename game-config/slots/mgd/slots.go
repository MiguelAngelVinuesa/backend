package mgd

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
)

// New return a Magic Devil slot-machine game with the given RTP if possible.
func New(rtp int) *game.Regular {
	switch rtp {
	case 92:
		return game.AcquireRegular(slots92params)
	case 94:
		return game.AcquireRegular(slots94params)
	case 96:
		return game.AcquireRegular(slots96params)
	}
	return nil
}

// NewLogged return a Magic Devil slot-machine game with the given RTP if possible and activates PRNG logging.
func NewLogged(rtp int) *game.Regular {
	switch rtp {
	case 92:
		s := slots92params
		s.PrngLog = true
		return game.AcquireRegular(s)
	case 94:
		s := slots94params
		s.PrngLog = true
		return game.AcquireRegular(s)
	case 96:
		s := slots96params
		s.PrngLog = true
		return game.AcquireRegular(s)
	}
	return nil
}

// NewWithRoundFlags returns an "Owl Kingdom" slot-machine game with the given RTP if possible and activates round flag logging.
func NewWithRoundFlags(rtp int) *game.Regular {
	switch rtp {
	case 92:
		s := slots92params
		s.ReturnFlags = true
		return game.AcquireRegular(s)
	case 94:
		s := slots94params
		s.ReturnFlags = true
		return game.AcquireRegular(s)
	case 96:
		s := slots96params
		s.ReturnFlags = true
		return game.AcquireRegular(s)
	}
	return nil
}

// AllSymbols returns the complete symbol set for Magic Devil.
func AllSymbols() *comp.SymbolSet {
	return symbols
}

// AllActions returns the complete list of actions for Magic Devil.
func AllActions(rtp int) comp.SpinActions {
	switch rtp {
	case 92:
		return actions92all
	case 94:
		return actions94all
	case 96:
		return actions96all
	}
	return nil
}

// Flags returns the flags for Owl Kingdom.
func Flags() comp.RoundFlags {
	return flags
}
