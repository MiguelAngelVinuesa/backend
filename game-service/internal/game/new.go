package game

import (
	"time"

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
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"

	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
)

func NewGame(game tg.GameNR, rtp int) *game.Regular {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.GeNewGame, started) }()

	switch game {
	case tg.BOTnr: // Book of Tomes
		return bot.NewLogged(rtp)
	case tg.CCBnr: // ChaCha Bomb
		return ccb.NewLogged(rtp)
	case tg.MGDnr: // Magic Devil
		return mgd.NewLogged(rtp)
	case tg.LAMnr: // La Modelo
		return lam.NewLogged(rtp)
	case tg.OWLnr: // Owl Kingdowm
		return owl.NewLogged(rtp)
	case tg.FRMnr: // Fruity Magic
		return frm.NewLogged(rtp)
	case tg.OFGnr: // 150 Ships
		return ofg.NewLogged(rtp)
	case tg.FPRnr: // Frost Princess
		return fpr.NewLogged(rtp)
	case tg.BTRnr: // Betic Riches
		return btr.NewLogged(rtp)
	case tg.BERnr: // Be Rich!
		return ber.NewLogged(rtp)
	case tg.CRWnr: // Cherry Reverse Win
		return crw.NewLogged(rtp)
	case tg.YYLnr: // Yin Yang Legacy
		return yyl.NewLogged(rtp)
	case tg.CASnr: // Casino Heist
		return cas.NewLogged(rtp)
	case tg.HOGnr: // Gunner's Hogs
	case tg.MOGnr: // Myth of Gunung
	case tg.BBSnr: // Barley-Broo's Saloon
	}
	return nil
}
