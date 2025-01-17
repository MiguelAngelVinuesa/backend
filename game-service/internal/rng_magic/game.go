package rng_magic

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/ber"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/bot"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/btr"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/ccb"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/crw"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/fpr"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/frm"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/lam"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/mgd"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/ofg"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/owl"
	magic "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/simulation/slots"
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
)

func GameData(gameNR tg.GameNR) (map[string]*magic.Condition, *comp.SymbolSet) {
	switch gameNR {
	case tg.BOTnr: // Book of Tomes
		return bot.Conditions(), bot.AllSymbols()
	case tg.CCBnr: // ChaCha Bomb
		return ccb.Conditions(), ccb.AllSymbols()
	case tg.MGDnr: // Magic Devil
		return mgd.Conditions(), mgd.AllSymbols()
	case tg.LAMnr: // La Modelo
		return lam.Conditions(), lam.AllSymbols()
	case tg.OWLnr: // Owl Kingdom
		return owl.Conditions(), owl.AllSymbols()
	case tg.FRMnr: // Fruity Magic
		return frm.Conditions(), frm.AllSymbols()
	case tg.OFGnr: // 150 Ships
		return ofg.Conditions(), ofg.AllSymbols()
	case tg.FPRnr: // Frosty Princess
		return fpr.Conditions(), fpr.AllSymbols()
	case tg.BTRnr: // Betic Riches
		return btr.Conditions(), btr.AllSymbols()
	case tg.HOGnr: // Gunner's Hogs
	case tg.MOGnr: // Myth of Gunung
	case tg.BBSnr: // Barley-Broo's Saloon
	case tg.BERnr: // Be Rich!
		return ber.Conditions(), ber.AllSymbols()
	case tg.CRWnr: // Cherry Reverse Win
		return crw.Conditions(), crw.AllSymbols()
	}

	return nil, nil
}
