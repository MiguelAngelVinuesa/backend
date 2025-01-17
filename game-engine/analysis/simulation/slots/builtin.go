package magic

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"

	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
)

var (
	Builtins = make(Conditions, 0, 64)
)

func MakeMatcher(key string, params map[string]any, symbols *comp.SymbolSet, game *game.Regular) Matcher {
	switch key {
	case KeyResultCountRange:
		return ResultCountRange(
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
		)

	case KeyPlayerChoice:
		return PlayerChoice(
			conv.StringFromAny(params[fieldChoice]),
			conv.StringFromAny(params[fieldValue]),
		)

	case KeySpinSymbolCount:
		return SpinSymbolCount(
			conv.IntFromAny(params[fieldSeq]),
			conv.IntFromAny(params[fieldCount]),
			SymbolIndex(symbols, params[fieldSymbol]),
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeyAnySymbolCount:
		return AnySymbolCount(
			conv.IntFromAny(params[fieldCount]),
			SymbolIndex(symbols, params[fieldSymbol]),
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeyTotalSymbolCount:
		return AllSymbolCount(
			conv.IntFromAny(params[fieldCount]),
			SymbolIndex(symbols, params[fieldSymbol]),
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeySpinSymbolRange:
		return SpinSymbolRange(
			conv.IntFromAny(params[fieldSeq]),
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
			SymbolIndex(symbols, params[fieldSymbol]),
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeyAnySymbolRange:
		return AnySymbolRange(
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
			SymbolIndex(symbols, params[fieldSymbol]),
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeyTotalSymbolRange:
		return AllSymbolRange(
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
			SymbolIndex(symbols, params[fieldSymbol]),
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeySpinSymbolTypeCount:
		return SpinHasSymbolTypeCount(
			conv.IntFromAny(params[fieldSeq]),
			conv.IntFromAny(params[fieldCount]),
			conv.StringFromAny(params[fieldSymbolType]),
			symbols,
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeyAnySymbolTypeCount:
		return AnyHasSymbolTypeCount(
			conv.IntFromAny(params[fieldCount]),
			conv.StringFromAny(params[fieldSymbolType]),
			symbols,
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeyTotalSymbolTypeCount:
		return AllHasSymbolTypeCount(
			conv.IntFromAny(params[fieldCount]),
			conv.StringFromAny(params[fieldSymbolType]),
			symbols,
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeySpinSymbolTypeRange:
		return SpinHasSymbolTypeRange(
			conv.IntFromAny(params[fieldSeq]),
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
			conv.StringFromAny(params[fieldSymbolType]),
			symbols,
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeyAnySymbolTypeRange:
		return AnyHasSymbolTypeRange(
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
			conv.StringFromAny(params[fieldSymbolType]),
			symbols,
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeyTotalSymbolTypeRange:
		return AllHasSymbolTypeRange(
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
			conv.StringFromAny(params[fieldSymbolType]),
			symbols,
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeySpinPaylinesRange:
		return SpinPaylinesRange(
			conv.IntFromAny(params[fieldSeq]),
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
		)

	case KeyAnyPaylinesRange:
		return AnyPaylinesRange(
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
		)

	case KeySpinPayoutRange:
		return SpinPayoutRange(
			conv.IntFromAny(params[fieldSeq]),
			conv.FloatFromAny(params[fieldMin]),
			conv.FloatFromAny(params[fieldMax]),
		)

	case KeyAnyPayoutRange:
		return AnyPayoutRange(
			conv.FloatFromAny(params[fieldMin]),
			conv.FloatFromAny(params[fieldMax]),
		)

	case KeyTotalPayoutRange:
		return AllPayoutRange(
			conv.FloatFromAny(params[fieldMin]),
			conv.FloatFromAny(params[fieldMax]),
		)

	case KeyMaxPayout:
		return MaxPayout(
			game.MaxPayout(),
		)

	case KeySpinPaylineCount:
		return SpinPaylineCount(
			conv.IntFromAny(params[fieldSeq]),
			conv.IntFromAny(params[fieldPayline]),
			conv.IntFromAny(params[fieldCount]),
		)

	case KeyAnyPaylineCount:
		return AnyPaylineCount(
			conv.IntFromAny(params[fieldPayline]),
			conv.IntFromAny(params[fieldCount]),
		)

	case KeySpinPaylineSymbol:
		return SpinPaylineSymbol(
			conv.IntFromAny(params[fieldSeq]),
			conv.IntFromAny(params[fieldPayline]),
			SymbolIndex(symbols, params[fieldSymbol]),
			conv.IntFromAny(params[fieldCount]),
		)

	case KeyAnyPaylineSymbol:
		return AnyPaylineSymbol(
			conv.IntFromAny(params[fieldPayline]),
			SymbolIndex(symbols, params[fieldSymbol]),
			conv.IntFromAny(params[fieldCount]),
		)

	case KeySpinStickyCount:
		return SpinStickyCount(
			conv.IntFromAny(params[fieldSeq]),
			conv.IntFromAny(params[fieldCount]),
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeyAnyStickyCount:
		return AnyStickyCount(
			conv.IntFromAny(params[fieldSeq]),
			game,
			conv.IntsFromAny(params[fieldReels]),
		)

	case KeySpinStickySymbol:
		return SpinStickySymbol(
			conv.IntFromAny(params[fieldSeq]),
			SymbolIndex(symbols, params[fieldSymbol]),
		)

	case KeyAnyStickySymbol:
		return AnyStickySymbol(
			SymbolIndex(symbols, params[fieldSymbol]),
		)

	case KeySpinStickyWild:
		return SpinStickyWild(
			conv.IntFromAny(params[fieldSeq]),
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
			symbols,
		)

	case KeyAnyStickyWild:
		return AnyStickyWild(
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
			symbols,
		)

	case KeySpinMultiplierCount:
		return SpinMultiplierCount(
			conv.IntFromAny(params[fieldSeq]),
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
			symbols,
		)

	case KeyAnyMultiplierCount:
		return AnyMultiplierCount(
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
			symbols,
		)

	case KeySpinMultiplierRange:
		return SpinMultiplierRange(
			conv.IntFromAny(params[fieldSeq]),
			uint16(conv.IntFromAny(params[fieldMin])),
			uint16(conv.IntFromAny(params[fieldMax])),
		)

	case KeyAnyMultiplierRange:
		return AnyMultiplierRange(
			uint16(conv.IntFromAny(params[fieldMin])),
			uint16(conv.IntFromAny(params[fieldMax])),
		)

	case KeyLastMultiplierRange:
		return LastMultiplierRange(
			uint16(conv.IntFromAny(params[fieldMin])),
			uint16(conv.IntFromAny(params[fieldMax])),
		)

	case KeySpinMultiplierSymbol:
		return SpinMultiplierSymbol(
			conv.IntFromAny(params[fieldSeq]),
			SymbolIndex(symbols, params[fieldSymbol]),
			float64(conv.IntFromAny(params[fieldMin])),
			float64(conv.IntFromAny(params[fieldMax])),
			symbols,
		)

	case KeyAnyMultiplierSymbol:
		return AnyMultiplierSymbol(
			SymbolIndex(symbols, params[fieldSymbol]),
			float64(conv.IntFromAny(params[fieldMin])),
			float64(conv.IntFromAny(params[fieldMax])),
			symbols,
		)

	case KeySpinMultiplierWild:
		return SpinMultiplierWild(
			conv.IntFromAny(params[fieldSeq]),
			float64(conv.IntFromAny(params[fieldMin])),
			float64(conv.IntFromAny(params[fieldMax])),
			symbols,
		)

	case KeyAnyMultiplierWild:
		return AnyMultiplierWild(
			float64(conv.IntFromAny(params[fieldMin])),
			float64(conv.IntFromAny(params[fieldMax])),
			symbols,
		)

	case KeyRoundMultiplier:
		return RoundMultiplier(
			conv.FloatFromAny(params[fieldMin]),
			conv.FloatFromAny(params[fieldMax]),
		)

	case KeyProgressMeter:
		return ProgressMeter(
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
		)

	case KeySpinScatterPayouts:
		return SpinScatterPayouts(
			conv.IntFromAny(params[fieldSeq]),
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
		)

	case KeyAnyScatterPayouts:
		return AnyScatterPayouts(
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
		)

	case KeySpinScatterPayout:
		return SpinScatterPayout(
			conv.IntFromAny(params[fieldSeq]),
			conv.FloatFromAny(params[fieldMin]),
			conv.FloatFromAny(params[fieldMax]),
		)

	case KeyAnyScatterPayout:
		return AnyScatterPayout(
			conv.FloatFromAny(params[fieldMin]),
			conv.FloatFromAny(params[fieldMax]),
		)

	case KeySpinScatterPayoutSymbol:
		return SpinScatterPayoutSymbol(
			conv.IntFromAny(params[fieldSeq]),
			SymbolIndex(symbols, params[fieldSymbol]),
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
		)

	case KeyAnyScatterPayoutSymbol:
		return AnyScatterPayoutSymbol(
			SymbolIndex(symbols, params[fieldSymbol]),
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
		)

	case KeySpinClusterPayouts:
		return SpinClusterPayouts(
			conv.IntFromAny(params[fieldSeq]),
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
		)

	case KeyAnyClusterPayouts:
		return AnyClusterPayouts(
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
		)

	case KeySpinClusterPayout:
		return SpinClusterPayout(
			conv.IntFromAny(params[fieldSeq]),
			conv.FloatFromAny(params[fieldMin]),
			conv.FloatFromAny(params[fieldMax]),
		)

	case KeyAnyClusterPayout:
		return AnyClusterPayout(
			conv.FloatFromAny(params[fieldMin]),
			conv.FloatFromAny(params[fieldMax]),
		)

	case KeySpinClusterPayoutSymbol:
		return SpinClusterPayoutSymbol(
			conv.IntFromAny(params[fieldSeq]),
			SymbolIndex(symbols, params[fieldSymbol]),
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
		)

	case KeyAnyClusterPayoutSymbol:
		return AnyClusterPayoutSymbol(
			SymbolIndex(symbols, params[fieldSymbol]),
			conv.IntFromAny(params[fieldMin]),
			conv.IntFromAny(params[fieldMax]),
		)

	case KeySpinGrid:
		return SpinGrid(
			conv.IntFromAny(params[fieldSeq]),
			indexesFromAny(params[fieldGrid]),
		)

	case KeyAnyGrid:
		return AnyGrid(
			indexesFromAny(params[fieldGrid]),
		)
	}

	return nil
}

func init() {
	Builtins = append(Builtins,
		// symbols
		NewCondition(KeySpinSymbolCount, NewSequenceParam(), NewSymbolParam(), NewIntParam(fieldCount, 0, 100), NewReelsParam()),
		NewCondition(KeyAnySymbolCount, NewSymbolParam(), NewIntParam(fieldCount, 0, 100), NewReelsParam()),
		NewCondition(KeyTotalSymbolCount, NewSymbolParam(), NewIntParam(fieldCount, 0, 1000), NewReelsParam()),
		NewCondition(KeySpinSymbolRange, NewSequenceParam(), NewSymbolParam(), NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100), NewReelsParam()),
		NewCondition(KeyAnySymbolRange, NewSymbolParam(), NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100), NewReelsParam()),
		NewCondition(KeyTotalSymbolRange, NewSymbolParam(), NewIntParam(fieldMin, 0, 1000), NewIntParam(fieldMax, 0, 1000), NewReelsParam()),

		// symbol types
		NewCondition(KeySpinSymbolTypeCount, NewSequenceParam(), NewSymbolTypeParam(), NewIntParam(fieldCount, 0, 100), NewReelsParam()),
		NewCondition(KeyAnySymbolTypeCount, NewSymbolTypeParam(), NewIntParam(fieldCount, 0, 100), NewReelsParam()),
		NewCondition(KeyTotalSymbolTypeCount, NewSymbolTypeParam(), NewIntParam(fieldCount, 0, 1000), NewReelsParam()),
		NewCondition(KeySpinSymbolTypeRange, NewSequenceParam(), NewSymbolTypeParam(), NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100), NewReelsParam()),
		NewCondition(KeyAnySymbolTypeRange, NewSymbolTypeParam(), NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100), NewReelsParam()),
		NewCondition(KeyTotalSymbolTypeRange, NewSymbolTypeParam(), NewIntParam(fieldMin, 0, 1000), NewIntParam(fieldMax, 0, 1000), NewReelsParam()),

		// paylines
		NewCondition(KeySpinPaylinesRange, NewSequenceParam(), NewIntParam(fieldMin, 0, 99999), NewIntParam(fieldMax, 0, 99999)),
		NewCondition(KeyAnyPaylinesRange, NewIntParam(fieldMin, 0, 99999), NewIntParam(fieldMax, 0, 99999)),
		NewCondition(KeySpinPaylineCount, NewSequenceParam(), NewPaylineParam(), NewIntParam(fieldCount, 2, 10)),
		NewCondition(KeyAnyPaylineCount, NewPaylineParam(), NewIntParam(fieldCount, 2, 10)),
		NewCondition(KeySpinPaylineSymbol, NewSequenceParam(), NewSymbolParam(), NewPaylineParam(), NewIntParam(fieldCount, 2, 10)),
		NewCondition(KeyAnyPaylineSymbol, NewPaylineParam(), NewSymbolParam(), NewIntParam(fieldCount, 2, 10)),

		// payouts
		NewCondition(KeySpinPayoutRange, NewSequenceParam(), NewFloatParam(fieldMin, 0, 99999), NewFloatParam(fieldMax, 0, 99999)),
		NewCondition(KeyAnyPayoutRange, NewFloatParam(fieldMin, 0, 99999), NewFloatParam(fieldMax, 0, 99999)),
		NewCondition(KeyTotalPayoutRange, NewFloatParam(fieldMin, 0, 99999), NewFloatParam(fieldMax, 0, 99999)),
		NewCondition(KeySpinScatterPayouts, NewSequenceParam(), NewIntParam(fieldMin, 1, 10), NewIntParam(fieldMax, 0, 10)),
		NewCondition(KeyAnyScatterPayouts, NewIntParam(fieldMin, 1, 10), NewIntParam(fieldMax, 0, 10)),
		NewCondition(KeySpinScatterPayout, NewSequenceParam(), NewFloatParam(fieldMin, 0.1, 999), NewFloatParam(fieldMax, 0, 999)),
		NewCondition(KeyAnyScatterPayout, NewFloatParam(fieldMin, 0.1, 999), NewFloatParam(fieldMax, 0, 999)),
		NewCondition(KeySpinScatterPayoutSymbol, NewSequenceParam(), NewSymbolParam(), NewIntParam(fieldMin, 1, 99), NewIntParam(fieldMax, 0, 99)),
		NewCondition(KeyAnyScatterPayoutSymbol, NewSymbolParam(), NewIntParam(fieldMin, 1, 99), NewIntParam(fieldMax, 0, 99)),
		NewCondition(KeySpinClusterPayouts, NewSequenceParam(), NewIntParam(fieldMin, 1, 10), NewIntParam(fieldMax, 0, 10)),
		NewCondition(KeyAnyClusterPayouts, NewIntParam(fieldMin, 1, 10), NewIntParam(fieldMax, 0, 10)),
		NewCondition(KeySpinClusterPayout, NewSequenceParam(), NewFloatParam(fieldMin, 0.1, 999), NewFloatParam(fieldMax, 0, 9999)),
		NewCondition(KeyAnyClusterPayout, NewFloatParam(fieldMin, 0.1, 999), NewFloatParam(fieldMax, 0, 9999)),
		NewCondition(KeySpinClusterPayoutSymbol, NewSequenceParam(), NewSymbolParam(), NewIntParam(fieldMin, 1, 99), NewIntParam(fieldMax, 0, 99)),
		NewCondition(KeyAnyClusterPayoutSymbol, NewSymbolParam(), NewIntParam(fieldMin, 1, 99), NewIntParam(fieldMax, 0, 99)),
		NewCondition(KeyMaxPayout),

		// stickiness
		NewCondition(KeySpinStickyCount, NewSequenceParam(), NewIntParam(fieldCount, 2, 10), NewReelsParam()),
		NewCondition(KeyAnyStickyCount, NewIntParam(fieldCount, 2, 10), NewReelsParam()),
		NewCondition(KeySpinStickySymbol, NewSequenceParam(), NewSymbolParam()),
		NewCondition(KeyAnyStickySymbol, NewSymbolParam()),
		NewCondition(KeySpinStickyWild, NewSequenceParam(), NewIntParam(fieldMin, 0, 99), NewIntParam(fieldMax, 0, 99)),
		NewCondition(KeyAnyStickyWild, NewIntParam(fieldMin, 0, 99), NewIntParam(fieldMax, 0, 99)),

		// multipliers
		NewCondition(KeySpinMultiplierCount, NewSequenceParam(), NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100)),
		NewCondition(KeyAnyMultiplierCount, NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100)),
		NewCondition(KeySpinMultiplierRange, NewSequenceParam(), NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100)),
		NewCondition(KeyAnyMultiplierRange, NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100)),
		NewCondition(KeyLastMultiplierRange, NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100)),
		NewCondition(KeySpinMultiplierSymbol, NewSequenceParam(), NewSymbolParam(), NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100)),
		NewCondition(KeyAnyMultiplierSymbol, NewSymbolParam(), NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100)),
		NewCondition(KeySpinMultiplierWild, NewSequenceParam(), NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100)),
		NewCondition(KeyAnyMultiplierWild, NewIntParam(fieldMin, 0, 100), NewIntParam(fieldMax, 0, 100)),
		NewCondition(KeyRoundMultiplier, NewFloatParam(fieldMin, 0, 9999), NewFloatParam(fieldMax, 0, 9999)),

		// misc
		NewCondition(KeyResultCountRange, NewIntParam(fieldMin, 0, 999), NewIntParam(fieldMax, 0, 999)),
		NewCondition(KeyPlayerChoice, NewStringParam(fieldChoice), NewStringParam(fieldValue)),
		NewCondition(KeyProgressMeter, NewIntParam(fieldMin, 1, 999), NewIntParam(fieldMax, 1, 999)),
		NewCondition(KeySpinGrid, NewSequenceParam(), NewGridParam()),
		NewCondition(KeyAnyGrid, NewGridParam()),
	)
}
