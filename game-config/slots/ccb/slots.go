package ccb

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
)

// New return a ChaCha Bomb slot-machine game with the given RTP if possible.
func New(rtp int) *game.Regular {
	switch rtp {
	case 92:
		return game.AcquireRegular(slots92param)
	case 94:
		return game.AcquireRegular(slots94param)
	case 96:
		return game.AcquireRegular(slots96param)
	}
	return nil
}

// NewLogged return a ChaCha Bomb slot-machine game with the given RTP if possible and activates PRNG logging.
func NewLogged(rtp int) *game.Regular {
	switch rtp {
	case 92:
		s := slots92param
		s.PrngLog = true
		return game.AcquireRegular(s)
	case 94:
		s := slots94param
		s.PrngLog = true
		return game.AcquireRegular(s)
	case 96:
		s := slots96param
		s.PrngLog = true
		return game.AcquireRegular(s)
	}
	return nil
}

// AllSymbols returns the complete symbol set for ChaCha Bomb.
func AllSymbols() *comp.SymbolSet {
	return symbols
}

// AllActions returns the complete symbol set for ChaCha Bomb.
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
