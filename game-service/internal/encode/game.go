package encode

import (
	"strconv"

	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/clients/bo_backend"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/hashes"
)

func GameInfo(loc string, sess *tg.SessionKey, g *game.Regular, debugEnabled bool) *zjson.Encoder {
	prefs, casino, juris := bo_backend.GetGamePrefs(loc, sess)
	defer prefs.Release()

	s := g.Slots()

	info := gameInfoPool.Acquire().(*gameInfo)
	info.debugEnabled = debugEnabled
	info.doubleSpin = s.DoubleSpin()
	info.roundMultiplier = s.RoundMultiplier()
	info.progressMeter = s.ProgressMeter()
	info.playerChoice = s.PlayerChoice()
	info.bonusBuy = s.BonusBuy()
	info.symbolsState = s.SymbolsState()
	info.reels = s.ReelCount()
	info.rows = s.RowCount()
	info.mask = s.ReelMask()
	info.maxPayout = g.MaxPayout()
	info.targetRTP = s.RTP()
	info.semVer = consts.SemVerShort
	info.semVerFull = consts.SemVerFull
	info.buildDate = consts.SemDate
	info.bets = casino["bets"].([]int64)
	info.configHash = hashes.GameHashes[sess.GameID()+strconv.Itoa(int(info.targetRTP))]

	switch sess.GameNr() {
	case tg.OFGnr:
		info.progressMax = 13
	case tg.FRMnr:
		info.progressMax = 160
	default:
	}

	if paylines := s.Paylines(); paylines != nil {
		info.paylines = paylines.Paylines()
	}

	for len(info.mask) < info.reels {
		info.mask = append(info.mask, uint8(info.rows))
	}

	all := s.AltSymbols()
	if all == nil {
		all = s.Symbols()
	}

	maxID := all.GetMaxSymbolID()
	for ix := util.Index(1); ix <= maxID; ix++ {
		if symbol := all.GetSymbol(ix); symbol != nil {
			info.symbols = append(info.symbols, newSymbol(symbol))
		}
	}

	spin, _ := g.LastSpin().Data.(*comp.SpinResult)
	if spin != nil {
		initial := spin.Initial()
		for ix := range initial {
			info.initial.Append(uint16(initial[ix]))
		}
	}

	enc := zjson.AcquireEncoder(2048)
	enc.StartObject()

	encodePrefs(enc, juris)

	enc.StartObjectField("gameData")
	info.Encode(enc)
	enc.EndObject()

	var bet int64

	if len(prefs.GamePrefs()) > 0 {
		enc.StartObjectField("preferences")
		for k, v := range prefs.GamePrefs() {
			switch k {
			case consts.PrefMusic, consts.PrefEffects, consts.PrefVolume:
				i64, _ := strconv.ParseInt(v, 10, 64)
				enc.Int64Field(k, i64)
			case consts.PrefBet:
				bet, _ = strconv.ParseInt(v, 10, 64)
				enc.Int64Field(k, bet)
			default:
				enc.StringField(k, v)
			}
		}
		enc.EndObject()
	}

	if sess.GameNr() == tg.CCBnr && bet > 0 {
		if flags := prefs.GetStateCCB(bet); flags != nil {
			enc.ObjectField("state", flags)
		}
	}

	enc.BoolField("success", true)

	enc.EndObject()
	return enc
}

func encodePrefs(enc *zjson.Encoder, m map[string]any) {
	enc.StartObjectField("jurisdiction")

	jur, _ := m["code"].(string)
	if jur == "" {
		jur = "MGA"
	}

	wait, _ := m["spinWait"].(int64)
	if wait == 0 {
		wait = 1000
	}

	enc.StringFieldOpt("code", jur)
	enc.Int64Field("spinWait", wait)

	for k, v := range m {
		switch k {
		case "code", "spinWait": // already done
		default:
			switch t := v.(type) {
			case bool:
				enc.IntBoolFieldOpt(k, t)
			case string:
				enc.StringFieldOpt(k, t)
			case int:
				enc.IntFieldOpt(k, t)
			case int64:
				enc.Int64FieldOpt(k, t)
			case uint8:
				enc.Uint8FieldOpt(k, t)
			case uint16:
				enc.Uint16FieldOpt(k, t)
			case uint64:
				enc.Uint64FieldOpt(k, t)
			case float64:
				enc.FloatFieldOpt(k, t, 'g', -1)
			}
		}
	}

	enc.EndObject()
}

func newSymbol(s *comp.Symbol) *gameSymbol {
	s2 := gameSymbolPool.Acquire().(*gameSymbol)
	s2.id = int(s.ID())
	s2.name = s.Name()
	s2.resource = s.Resource()
	s2.multiplier = s.Multiplier()
	s2.kind = int(s.Kind())
	if len(s.Payouts()) > 0 {
		s2.payTable.Append(s.Payouts()...)
	} else if len(s.ScatterPayouts()) > 0 {
		s2.payTable.Append(s.ScatterPayouts()...)
	}
	return s2
}

type gameInfo struct {
	debugEnabled    bool
	doubleSpin      bool
	roundMultiplier bool
	progressMeter   bool
	playerChoice    bool
	bonusBuy        bool
	symbolsState    bool
	reels           int
	rows            int
	progressMax     int
	maxPayout       float64
	targetRTP       float64
	semVer          string
	semVerFull      string
	buildDate       string
	configHash      string
	mask            []uint8
	bets            []int64
	symbols         []*gameSymbol
	paylines        comp.Paylines
	initial         *object.Uint16sManager
	pool.Object
}

var gameInfoPool = pool.NewProducer(func() (pool.Objecter, func()) {
	g := &gameInfo{
		initial: indexesPool.Acquire().(*object.Uint16sManager),
	}
	return g, g.reset
})

func (i *gameInfo) reset() {
	for ix := range i.symbols {
		i.symbols[ix].Release()
		i.symbols[ix] = nil
	}

	i.debugEnabled = false
	i.doubleSpin = false
	i.roundMultiplier = false
	i.progressMeter = false
	i.playerChoice = false
	i.bonusBuy = false
	i.symbolsState = false
	i.reels = 0
	i.rows = 0
	i.progressMax = 0
	i.maxPayout = 0
	i.targetRTP = 0
	i.semVer = ""
	i.semVerFull = ""
	i.buildDate = ""
	i.configHash = ""
	i.mask = nil
	i.bets = nil
	i.symbols = i.symbols[:0]
	i.paylines = nil
	i.initial.Items = object.ResetUint16s(i.initial.Items, 16, 0, true)
}

func (i *gameInfo) Encode(enc *zjson.Encoder) {
	enc.IntField("reels", i.reels)
	enc.IntField("rows", i.rows)

	if len(i.mask) > 0 {
		enc.StartArrayField("mask")
		for _, m := range i.mask {
			enc.Uint64(uint64(m))
		}
		enc.EndArray()
	}

	enc.Key("initial")
	i.initial.Encode(enc)

	enc.StartArrayField("symbols")
	for ix := range i.symbols {
		enc.StartObject()
		i.symbols[ix].Encode(enc)
		enc.EndObject()
	}
	enc.EndArray()

	if i.paylines != nil {
		enc.StartArrayField("paylines")
		for ix := range i.paylines {
			p := i.paylines[ix]
			enc.StartObject()
			enc.Uint8Field("id", p.ID())
			enc.StartArrayField("rowMap")
			for _, r := range p.RowMap() {
				enc.Uint64(uint64(r))
			}
			enc.EndArray()
			enc.EndObject()
		}
		enc.EndArray()
	}

	enc.StringField("version", i.semVer)
	enc.StringField("versionFull", i.semVerFull)
	enc.StringFieldOpt("buildDate", i.buildDate)
	enc.IntBoolFieldOpt("debugEnabled", i.debugEnabled)
	enc.StringField("configHash", i.configHash)
	enc.IntBoolFieldOpt("doubleSpin", i.doubleSpin)
	enc.IntBoolFieldOpt("roundMultiplier", i.roundMultiplier)
	enc.IntBoolFieldOpt("progressMeter", i.progressMeter)
	enc.IntBoolFieldOpt("playerChoice", i.playerChoice)
	enc.IntBoolFieldOpt("bonusBuy", i.bonusBuy)
	enc.IntBoolFieldOpt("symbolsState", i.symbolsState)
	enc.FloatFieldOpt("maxPayout", i.maxPayout, 'f', 2)
	enc.FloatFieldOpt("targetRTP", i.targetRTP, 'f', 2)
	enc.IntFieldOpt("progressMax", i.progressMax)

	enc.StartArrayField("bets")
	for _, b := range i.bets {
		enc.Int64(b)
	}
	enc.EndArray()
}

type gameSymbol struct {
	id         int
	kind       int
	multiplier float64
	name       string
	resource   string
	payTable   *object.FloatsManager
	pool.Object
}

var gameSymbolPool = pool.NewProducer(func() (pool.Objecter, func()) {
	g := &gameSymbol{
		payTable: floatsPool.Acquire().(*object.FloatsManager),
	}
	return g, g.reset
})

func (s *gameSymbol) reset() {
	s.id = 0
	s.kind = 0
	s.multiplier = 0.0
	s.name = ""
	s.resource = ""
	s.payTable.Items = object.ResetFloats(s.payTable.Items, 16, 0, true)
}

func (s *gameSymbol) Encode(enc *zjson.Encoder) {
	enc.IntFieldOpt("id", s.id)
	enc.IntFieldOpt("kind", s.kind)
	enc.StringFieldOpt("name", s.name)
	enc.StringFieldOpt("resource", s.resource)

	if s.multiplier != 0.0 && s.multiplier != 1.0 {
		enc.FloatField("multiplier", s.multiplier, 'f', 2)
	}

	if !s.payTable.IsEmpty() {
		enc.Key("payTable")
		s.payTable.Encode(enc)
	}
}

var (
	indexesPool = object.NewUint16sProducer(16, 64, true)
	floatsPool  = object.NewFloatsProducer(16, 64, true, 'f', 2)
)
