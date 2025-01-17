package slots

import (
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/ber"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/bot"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/btr"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/cas"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/ccb"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/crw"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/fpr"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/frm"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/lam"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/mgd"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/ofg"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/owl"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/yyl"
)

func NewGame(game tg.GameNR, rtp int, rngBuf, flags bool) *game.Regular {
	switch game {
	case tg.BOTnr:
		return NewGameBOT(rtp, rngBuf, flags)
	case tg.CCBnr:
		return NewGameCCB(rtp, rngBuf, flags)
	case tg.MGDnr:
		return NewGameMGD(rtp, rngBuf, flags)
	case tg.LAMnr:
		return NewGameLAM(rtp, rngBuf, flags)
	case tg.OWLnr:
		return NewGameOWL(rtp, rngBuf, flags)
	case tg.FRMnr:
		return NewGameFRM(rtp, rngBuf, flags)
	case tg.OFGnr:
		return NewGameOFG(rtp, rngBuf, flags)
	case tg.FPRnr:
		return NewGameFPR(rtp, rngBuf, flags)
	case tg.BTRnr:
		return NewGameBTR(rtp, rngBuf, flags)
	case tg.BERnr:
		return NewGameBER(rtp, rngBuf, flags)
	case tg.CRWnr:
		return NewGameCRW(rtp, rngBuf, flags)
	case tg.YYLnr:
		return NewGameYYL(rtp, rngBuf, flags)
	case tg.CASnr:
		return NewGameCAS(rtp, rngBuf, flags)
	default:
		return nil
	}
}

func NewGameBOT(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return bot.NewLogged(rtp)
	}
	return bot.New(rtp)
}

func NewGameCCB(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return ccb.NewLogged(rtp)
	}
	return ccb.New(rtp)
}

func NewGameMGD(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return mgd.NewLogged(rtp)
	}
	if flags {
		return mgd.NewWithRoundFlags(rtp)
	}
	return mgd.New(rtp)
}

func NewGameLAM(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return lam.NewLogged(rtp)
	}
	return lam.New(rtp)
}

func NewGameOWL(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return owl.NewLogged(rtp)
	}
	return owl.NewWithRoundFlags(rtp)
}

func NewGameFRM(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return frm.NewLogged(rtp)
	}
	return frm.New(rtp)
}

func NewGameOFG(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return ofg.NewLogged(rtp)
	}
	if flags {
		return ofg.NewWithRoundFlags(rtp)
	}
	return ofg.New(rtp)
}

func NewGameFPR(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return fpr.NewLogged(rtp)
	}
	if flags {
		return fpr.NewWithRoundFlags(rtp)
	}
	return fpr.New(rtp)
}

func NewGameBTR(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return btr.NewLogged(rtp)
	}
	if flags {
		return btr.NewWithRoundFlags(rtp)
	}
	return btr.New(rtp)
}

func NewGameBER(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return ber.NewLogged(rtp)
	}
	if flags {
		return ber.NewWithRoundFlags(rtp)
	}
	return ber.New(rtp)
}

func NewGameCRW(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return crw.NewLogged(rtp)
	}
	return crw.New(rtp)
}

func NewGameYYL(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return yyl.NewLogged(rtp)
	}
	return yyl.New(rtp)
}

func NewGameCAS(rtp int, rngBuf, flags bool) *game.Regular {
	if rngBuf {
		return cas.NewLogged(rtp)
	}
	return cas.New(rtp)
}
