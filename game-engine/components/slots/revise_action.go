package slots

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// ReviseAction contains the details for actions that revise the grid.
// These actions exist to give more control in grid outcomes,
// and to manage the hit-rate and RTP in a more controlled manner.
type ReviseAction struct {
	SpinAction
	spinKinds []SpinKind // optional list of spin kinds when the action is allowed to trigger.

	// details for morphing one or more symbols into other symbols.
	morphSymbols bool          // true if symbol morphing is requested.
	morphReels   utils.UInt8s  // optional list of reels to consider (1-based).
	morphChances []float64     // chances of morph happening (max 2 decimals).
	morphFor     utils.Indexes // optional list of symbols to consider.

	// details for generating one or more special symbols on the grid.
	generateSymbol      bool                    // true if symbol generation is requested.
	symbolChances       []float64               // chances for 1,2,3,etc. symbols to be generated (max 2 decimals).
	generateReels       utils.UInt8s            // optional list of reels to consider (1-based).
	genAllowDupes       bool                    // indicates if duplicate symbols are allowed across reels.
	genAllowOld         bool                    // indicates if the symbol may already occur before the action is called.
	genMultiplier       utils.WeightedGenerator // optional weighting to generate multipliers.
	generateReelSymbols bool                    // true if reel weighted symbol generation is requested.
	symbolWeights       utils.WeightedGenerator // weighting for number of symbols to generate.
	reelWeights         utils.WeightedGenerator // weighting for reels to generate symbols on.
	generateSymbols     bool                    // true if multiple symbol generation is requested.
	genWeights          utils.WeightedGenerator // weighting for possible symbols to generate.

	// details for generating the bonus symbol on multiple reels of the grid.
	generateBonus bool      // true if bonus symbol generation is requested.
	bonusCount    uint8     // count of reels to contain the bonus symbol
	bonusChances  []float64 // chances (per symbolID) for the bonus symbol to be generated (max 2 decimals).

	// details for forcing a specific shape of symbols on the grid.
	generateShape bool                    // true if shape generation is requested.
	shapeChance   float64                 // chance of this action to be triggered (max 2 decimals).
	shapeCenters  GridOffsets             // center positions for the shape to choose from.
	shapeWeights  utils.WeightedGenerator // weights for the symbols to choose from.
	shapeGrid     GridOffsets             // the shape to generate.

	// details for all-symbol de-duplication.
	deduplicate bool                    // true if all symbol de-duplication is requested.
	dupWeights  utils.WeightedGenerator // weighting of the maximum duplicates.

	// details for single symbol de-duplication.
	dedupSymbol bool         // true if single symbol de-duplication is requested.
	dedupReels  utils.UInt8s // reels on which to perform de-duplication (if empty all reels will be checked; 1-based).

	// details for replacing one or more symbols from a reel based on the existence of another symbol.
	replaceSymbols bool          // true is symbols replacement is requested.
	detectReels    utils.UInt8s  // reels to detect the trigger symbol (1-based).
	replSymbols    utils.Indexes // symbols to be replaced.
	replChances    []float64     // chances for replacing the corresponding symbol (max 2 decimals).
	replReels      utils.UInt8s  // reels to replace the symbols from (1-based).

	// details for preventing payouts on the grid.
	preventPayouts   bool         // true if payline removal is requested.
	prevPayoutChance float64      // chance of this action to be triggered (max 2 decimals).
	prevPayoutDir    PayDirection // LTR, RTL or both.
	prevPayoutMech   int          // removal mechanism (1=dedupe 1st/2nd reels; 2=dedupe 2nd/3rd reel).
	prevPayoutWilds  bool         // always remove wild symbols.
}

// NewMorphSymbolAction instantiates a new symbol morphing action.
// Use the options to limit the reels and/or symbols being considered.
// Note that reels are 1-based (there is no reel 0)!
func NewMorphSymbolAction(reels utils.UInt8s, chances []float64, symbols ...utils.Index) *ReviseAction {
	a := newReviseAction()
	a.morphSymbols = true
	a.morphReels = reels
	a.morphChances = chances
	a.morphFor = symbols
	return a.finalizer()
}

// NewGenerateSymbolAction instantiates a new symbol generator action.
// The slice of chances indicates the maximum number of new symbols to be generated.
// E.g. if the slice contains 2 entries, zero, one or two symbols will be generated; The first entry gives the chance
// for one symbol, the second entry the chance for a second symbol *but* only if the first was already generated.
// Chances are a percentage, so maximum is 100% and minimum 0%; a maximum of 2 decimals is supported.
// The location in the grid where the new symbols may appear is chosen randomly.
// Note that reels are 1-based (there is no reel 0)!
func NewGenerateSymbolAction(symbol utils.Index, chances []float64, reels ...uint8) *ReviseAction {
	a := newReviseAction()
	a.generateSymbol = true
	a.symbol = symbol
	a.symbolChances = chances
	a.generateReels = reels
	return a.finalizer()
}

// NewGenerateReelSymbolAction instantiates a new symbol generator action.
// Note that reels are 1-based (there is no reel 0)!
func NewGenerateReelSymbolAction(symbol utils.Index, weights, reels utils.WeightedGenerator) *ReviseAction {
	a := newReviseAction()
	a.generateReelSymbols = true
	a.symbol = symbol
	a.symbolWeights = weights
	a.reelWeights = reels
	return a.finalizer()
}

// NewGenerateSymbolsAction instantiates a new symbol generator action.
// The slice of chances indicates the maximum number of new symbols to be generated.
// E.g. if the slice contains 2 entries, zero, one or two symbols will be generated; The first entry gives the chance
// for one symbol, the second entry the chance for a second symbol *but* only if the first was already generated.
// Chances are a percentage, so maximum is 100% and minimum 0%; a maximum of 2 decimals is supported.
// The weights determine which of the symbols have a higher chance of geting selected.
// The location in the grid where the new symbols may appear is chosen randomly.
// Note that reels are 1-based (there is no reel 0)!
func NewGenerateSymbolsAction(weights utils.WeightedGenerator, chances []float64, reels ...uint8) *ReviseAction {
	a := newReviseAction()
	a.generateSymbols = true
	a.symbolChances = chances
	a.generateReels = reels
	a.genWeights = weights
	return a.finalizer()
}

// NewGenOrMorphSymbolAction instantiates a new symbol generator action.
// It either morphs from an existing symbol, or generates a new symbol.
// The slice of chances indicates the maximum number of new symbols to be generated or morphed.
// E.g. if the slice contains 2 entries, zero, one or two symbols will be generated; The first entry gives the chance
// for one symbol, the second entry the chance for a second symbol *but* only if the first was already generated.
// Chances are a percentage, so maximum is 100% and minimum 0%; a maximum of 2 decimals is supported.
// The location in the grid where the new symbols may appear is chosen randomly.
// Note that reels are 1-based (there is no reel 0)!
func NewGenOrMorphSymbolAction(symbol, from utils.Index, chances []float64, reels ...uint8) *ReviseAction {
	a := newReviseAction()
	a.morphSymbols = true // we try to morph first
	a.morphReels = reels
	a.morphChances = chances
	a.morphFor = utils.Indexes{from}
	a.generateSymbol = true // if the morph fails, we generate :)
	a.symbol = symbol
	a.symbolChances = chances
	a.generateReels = reels
	return a.finalizer()
}

// NewGenerateBonusAction instantiates a new bonus symbol generator action.
// The chance indicates the percent of times the generation will take place (max 2 decimals).
// Up to count bonus symbols will be generated, depending on how many are already present and the number of hot reels.
// The reel in the grid where the new symbols are generated is chosen randomly.
func NewGenerateBonusAction(count uint8, chances []float64) *ReviseAction {
	a := newReviseAction()
	a.generateBonus = true
	a.bonusCount = count
	a.bonusChances = chances
	return a.finalizer()
}

// NewGenerateShapeAction instantiates a new shape generator action.
// The chance is a percentage, so maximum is 100% and minimum 0%; a maximum of 4 decimals is supported.
// The chance can not be met exactly. Read below for the reason.
// The location in the grid where the shape may appear is chosen randomly from the centers slice.
// The symbol is chosen using the given weighting table.
// Note that sometimes the chosen symbol needs to be changed into the sticky symbol when the shape touches a sticky tile.
// Also note that sometimes the shape can not be generated because unchanged tiles within the shape already contain the same symbol.
// This latter situation actually decreases the chance of generating the shape!
func NewGenerateShapeAction(chance float64, shape, centers GridOffsets, weights utils.WeightedGenerator) *ReviseAction {
	a := newReviseAction()
	a.generateShape = true
	a.shapeChance = chance
	a.shapeGrid = shape
	a.shapeCenters = centers
	a.shapeWeights = weights
	return a.finalizer()
}

// NewDeduplicationAction instantiates a grid de-duplication action.
// It uses the weights to determine how many duplicates are allowed, and then tests the grid.
// If a symbol has too many duplicates they will randomly be transformed into another randomly chosen symbol
// that has less than the maximum allowed duplicates.
// The deduplication will panic if there are more grid positions than symbols, and the weighting produces a 1.
func NewDeduplicationAction(weights utils.WeightedGenerator) *ReviseAction {
	a := newReviseAction()
	a.deduplicate = true
	a.dupWeights = weights
	return a.finalizer()
}

// NewDedupSymbolAction instantiates a single symbol de-duplication action.
// The given reels will be checked for duplicates. If the slice is empty all reels will be checked.
// Note that reels are 1-based (there is no reel 0)!
func NewDedupSymbolAction(symbol utils.Index, reels utils.UInt8s) *ReviseAction {
	a := newReviseAction()
	a.symbol = symbol
	a.dedupSymbol = true
	a.dedupReels = reels
	return a.finalizer()
}

// NewReplaceSymbolsAction instantiates a symbols remove action.
// Removal takes place if the trigger symbol is found in the given detection reels.
// For each symbol to replace there must be a corresponding chance, expressed as % (with up to 2 decimals).
// If the chance is awarded, the corresponding symbol will be replaced with another symbol using the PRNG.
// Note that reels are 1-based (there is no reel 0)!
func NewReplaceSymbolsAction(symbol utils.Index, reels utils.UInt8s, replaceSymbols utils.Indexes,
	replaceChances []float64, replaceReels utils.UInt8s) *ReviseAction {
	a := newReviseAction()
	a.symbol = symbol
	a.replaceSymbols = true
	a.detectReels = reels
	a.replSymbols = replaceSymbols
	a.replChances = replaceChances
	a.replReels = replaceReels
	return a.finalizer()
}

// NewPreventPayoutsAction instantiates a payout prevention action.
func NewPreventPayoutsAction(chance float64, direction PayDirection, mechanism int, removeWilds bool) *ReviseAction {
	a := newReviseAction()
	a.preventPayouts = true
	a.prevPayoutChance = chance
	a.prevPayoutDir = direction
	a.prevPayoutMech = mechanism
	a.prevPayoutWilds = removeWilds
	return a.finalizer()
}

// WithSpinKinds limits the action to the given spin kinds.
func (a *ReviseAction) WithSpinKinds(kinds []SpinKind) *ReviseAction {
	a.spinKinds = kinds
	return a.finalizer()
}

// WithAlternate adds an alternative morphing action when this action doesn't trigger.
func (a *ReviseAction) WithAlternate(alt *ReviseAction) *ReviseAction {
	a.alternate = alt
	return a.finalizer()
}

// WithMultipliers adds multiplier weightings used when generating symbols.
func (a *ReviseAction) WithMultipliers(w utils.WeightedGenerator) *ReviseAction {
	a.genMultiplier = w
	return a.finalizer()
}

// GenerateNoDupes disallows to generate duplicate symbols across reels.
func (a *ReviseAction) GenerateNoDupes() *ReviseAction {
	a.genAllowDupes = false
	return a.finalizer()
}

// AllowPrevious allows to generate symbols when the symbol already appears in teh grid.
func (a *ReviseAction) AllowPrevious() *ReviseAction {
	a.genAllowOld = true
	return a.finalizer()
}

// Triggered implements the SpinAction.Triggered interface.
func (a *ReviseAction) Triggered(spin *Spin) SpinActioner {
	if len(a.spinKinds) > 0 && !a.kindAllowed(spin) {
		return nil
	}

	switch {
	case a.generateSymbol:
		if a.doGenerateSymbol(spin) {
			return a
		}
	case a.generateReelSymbols:
		if a.doGenerateReelSymbols(spin) {
			return a
		}
	case a.generateBonus:
		if a.doGenerateBonus(spin) {
			return a
		}
	case a.generateShape:
		if a.doGenerateShape(spin) {
			return a
		}
	case a.deduplicate:
		if a.doDeduplication(spin) {
			return a
		}
	case a.dedupSymbol:
		if a.doDedupSymbol(spin) {
			return a
		}
	case a.replaceSymbols:
		if a.doReplaceSymbols(spin) {
			return a
		}
	case a.preventPayouts:
		if a.doRemovePayouts(spin) {
			return a
		}
	case a.morphSymbols:
		if a.doMorphSymbols(spin) {
			return a
		}
	case a.generateSymbols:
		if a.doGenerateSymbols(spin) {
			return a
		}
	}

	if a.alternate != nil {
		return a.alternate.Triggered(spin)
	}
	return nil
}

func (a *ReviseAction) kindAllowed(spin *Spin) bool {
	for ix := range a.spinKinds {
		if spin.kind == a.spinKinds[ix] {
			return true
		}
	}
	return false
}

// newReviseAction instantiates a new basic morphing action.
func newReviseAction() *ReviseAction {
	a := &ReviseAction{}
	a.init(ReviseGrid, Processed, reflect.TypeOf(a).String())
	a.genAllowDupes = true
	a.prevPayoutDir = PayLTR
	return a
}

func (a *ReviseAction) finalizer() *ReviseAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	switch {
	case a.morphSymbols || a.generateSymbol: // joined here as they may both be true.
		if a.morphSymbols {
			b.WriteString(",morphSymbols=true")
			if len(a.morphReels) > 0 {
				b.WriteString(",reels=")
				j, _ := json.Marshal(a.morphSymbols)
				b.Write(j)
			}
			b.WriteString(",chances=")
			j, _ := json.Marshal(a.morphChances)
			b.Write(j)
			b.WriteString(",symbols=")
			j, _ = json.Marshal(a.morphSymbols)
			b.Write(j)
		}

		if a.generateSymbol {
			b.WriteString(",generateSymbol=true")
			b.WriteString(",symbol=")
			b.WriteString(strconv.Itoa(int(a.symbol)))
			b.WriteString(",chances=")
			j, _ := json.Marshal(a.symbolChances)
			b.Write(j)
			if len(a.generateReels) > 0 {
				b.WriteString(",reels=")
				j, _ := json.Marshal(a.generateReels)
				b.Write(j)
			}
			b.WriteString(",dupes=")
			if a.genAllowDupes {
				b.WriteString("true")
			} else {
				b.WriteString("false")
			}
			if a.genMultiplier != nil {
				b.WriteString(",multiplier=")
				b.WriteString(a.genMultiplier.String())
			}
		}

	case a.generateReelSymbols:
		b.WriteString(",generateReelSymbols=true")
		b.WriteString(",symbol=")
		b.WriteString(strconv.Itoa(int(a.symbol)))
		b.WriteString(",chances=")
		j, _ := json.Marshal(a.symbolChances)
		b.Write(j)
		b.WriteString(",reels=")
		b.WriteString(a.reelWeights.String())
		b.WriteString(",dupes=")
		if a.genAllowDupes {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
		if a.genMultiplier != nil {
			b.WriteString(",multiplier=")
			b.WriteString(a.genMultiplier.String())
		}

	case a.generateSymbols:
		b.WriteString(",generateSymbols=true")
		b.WriteString(",weights=")
		b.WriteString(a.genWeights.String())
		b.WriteString(",chances=")
		j, _ := json.Marshal(a.symbolChances)
		b.Write(j)
		if len(a.generateReels) > 0 {
			b.WriteString(",reels=")
			j, _ := json.Marshal(a.generateReels)
			b.Write(j)
		}
		b.WriteString(",dupes=")
		if a.genAllowDupes {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}

	case a.generateBonus:
		b.WriteString(",generateBonus=true")
		b.WriteString(",count=")
		b.WriteString(strconv.Itoa(int(a.bonusCount)))
		b.WriteString(",chances=")
		b.WriteString(fmt.Sprintf("%v", a.bonusChances))

	case a.generateShape:
		b.WriteString(",generateShape=true")
		b.WriteString(",chance=")
		b.WriteString(strconv.FormatFloat(a.shapeChance, 'f', 2, 64))
		b.WriteString(",centers=")
		j, _ := json.Marshal(a.shapeCenters)
		b.Write(j)
		b.WriteString(",shape=")
		j, _ = json.Marshal(a.shapeGrid)
		b.Write(j)

	case a.deduplicate:
		b.WriteString(",deduplicate=true")

	case a.dedupSymbol:
		b.WriteString(",dedupSymbol=true")
		if len(a.dedupReels) > 0 {
			b.WriteString(",reels=")
			j, _ := json.Marshal(a.dedupReels)
			b.Write(j)
		}

	case a.replaceSymbols:
		b.WriteString(",replaceSymbols=true")
		if len(a.detectReels) > 0 {
			b.WriteString(",detect=")
			j, _ := json.Marshal(a.detectReels)
			b.Write(j)
		}
		b.WriteString(",symbols=")
		j, _ := json.Marshal(a.replSymbols)
		b.Write(j)
		b.WriteString(",chances=")
		j, _ = json.Marshal(a.replChances)
		b.Write(j)
		if len(a.replReels) > 0 {
			b.WriteString(",reels=")
			j, _ = json.Marshal(a.replReels)
			b.Write(j)
		}

	case a.preventPayouts:
		b.WriteString(",preventPayouts=true")
		b.WriteString(",chance=")
		b.WriteString(strconv.FormatFloat(a.prevPayoutChance, 'f', 2, 64))
		b.WriteString(",dir=")
		b.WriteString(strconv.Itoa(int(a.prevPayoutDir)))
		b.WriteString(",mech=")
		b.WriteString(strconv.Itoa(a.prevPayoutMech))
		b.WriteString(",wilds=")
		if a.prevPayoutWilds {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	}

	a.config = b.String()
	return a
}
