package slots

import (
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// MaxSymbolID is the highest supported symbol id.
const MaxSymbolID = 99

// SymbolKind represents a kind of symbol.
type SymbolKind uint8

// list of pre-defined symbol kinds.
const (
	// Standard is a normal symbol.
	Standard SymbolKind = iota + 1
	// Split is a symbol representing multiple symbols (partial wild).
	Split
	// Wild is a wild symbol.
	Wild
	// Hero is a hero symbol.
	Hero
	// Scatter is a scatter symbol.
	Scatter
	// Bomb is a bomb symbol.
	Bomb
	// Shooter is a shooter symbol.
	Shooter
	// Prize is a prize symbol.
	Prize
	// WildScatter is a wild and scatter symbol.
	WildScatter
	// HeroScatter is a hero and scatter symbol.
	HeroScatter
	// WildBomb is a wild and bomb symbol.
	WildBomb
	// WildShooter is a wild and shooter symbol.
	WildShooter
	// ScatterShooter is a scatter and shooter symbol.
	ScatterShooter
)

// String implements the Stringer interface.
func (k SymbolKind) String() string {
	switch k {
	case Standard:
		return "Standard"
	case Split:
		return "Split"
	case Wild:
		return "Wild"
	case Hero:
		return "Hero"
	case Scatter:
		return "Scatter"
	case Bomb:
		return "Bomb"
	case Shooter:
		return "Shooter"
	case Prize:
		return "Prize"
	case WildScatter:
		return "Wild-Scatter"
	case HeroScatter:
		return "Hero-Scatter"
	case WildBomb:
		return "Wild-Bomb"
	case WildShooter:
		return "Wild-Shooter"
	case ScatterShooter:
		return "Scatter-Shooter"
	default:
		return "[unknown]"
	}
}

func (k SymbolKind) isSplit() bool {
	return k == Split
}

func (k SymbolKind) isWild() bool {
	return k == Wild || k == WildScatter || k == WildBomb || k == WildShooter
}

func (k SymbolKind) isHero() bool {
	return k == Hero || k == HeroScatter
}

func (k SymbolKind) isScatter() bool {
	return k == Scatter || k == WildScatter || k == HeroScatter || k == ScatterShooter
}

func (k SymbolKind) isBomb() bool {
	return k == Bomb || k == WildBomb
}

func (k SymbolKind) isShooter() bool {
	return k == Shooter || k == WildShooter || k == ScatterShooter
}

func (k SymbolKind) isPrize() bool {
	return k == Prize
}

// Symbol represents a single symbol for a slot machine.
type Symbol struct {
	sticky         bool
	isSplit        bool
	isWild         bool
	isScatter      bool
	isHero         bool
	isBomb         bool
	isShooter      bool
	isPrize        bool
	varyMultiplier bool
	kind           SymbolKind
	minPayable     uint8
	id             utils.Index
	morphInto      utils.Index
	multiplier     float64
	name           string
	resource       string
	weights        []float64
	payouts        []float64
	scatterPayouts []float64
	wildFor        utils.Indexes
	clearPattern   []int
}

// NewSymbol instantiates a new symbol.
// The options can be used to change the symbol characteristics as required.
// Symbols are immutable once created, so they are safe to use across concurrent go-routines.
// The function will panic if an invalid symbol id (> MaxSymbolID) is given.
func NewSymbol(id utils.Index, opts ...SymbolOption) *Symbol {
	if id > MaxSymbolID {
		panic(consts.MsgInvalidSymbolID)
	}

	s := &Symbol{id: id, kind: Standard, multiplier: 1.0}
	for ix := range opts {
		opts[ix](s)
	}
	return s.init()
}

// ID returns the symbol id.
func (s *Symbol) ID() utils.Index {
	return s.id
}

// Name returns the symbol name.
func (s *Symbol) Name() string {
	return s.name
}

// Resource returns the symbol resource.
func (s *Symbol) Resource() string {
	return s.resource
}

// Kind returns the symbol kind.
func (s *Symbol) Kind() SymbolKind {
	return s.kind
}

// IsSticky returns true if the symbol is sticky.
func (s *Symbol) IsSticky() bool {
	return s.sticky
}

// IsSplit returns true if the symbol represents a split (semi-wild).
func (s *Symbol) IsSplit() bool {
	return s.isSplit
}

// IsWild returns true if the symbol represents a Wild.
func (s *Symbol) IsWild() bool {
	return s.isWild
}

// IsHero returns true if the symbol represents a Hero.
func (s *Symbol) IsHero() bool {
	return s.isHero
}

// IsScatter returns true if the symbol represents a Scatter.
func (s *Symbol) IsScatter() bool {
	return s.isScatter
}

// IsBomb returns true if the symbol represents a bomb.
func (s *Symbol) IsBomb() bool {
	return s.isBomb
}

// IsShooter returns true if the symbol represents a shooter.
func (s *Symbol) IsShooter() bool {
	return s.isShooter
}

// IsPrize returns true if the symbol represents a prize.
func (s *Symbol) IsPrize() bool {
	return s.isPrize
}

// VaryMultiplier returns true if the symbol can be associated with an arbitrary multiplier during spins.
func (s *Symbol) VaryMultiplier() bool {
	return s.varyMultiplier
}

// Weights returns the symbol weights.
func (s *Symbol) Weights() []float64 {
	return s.weights
}

// Payouts returns the symbol payline paytable.
func (s *Symbol) Payouts() []float64 {
	return s.payouts
}

// ScatterPayouts returns the symbol scatter paytable.
func (s *Symbol) ScatterPayouts() []float64 {
	return s.scatterPayouts
}

// Payable returns true if the count matches a positive payout for the symbol.
func (s *Symbol) Payable(count uint8) bool {
	return count >= s.minPayable
}

// MinPayable returns the count of symbols required to warrant a payout.
func (s *Symbol) MinPayable() uint8 {
	return s.minPayable
}

// Payout returns the payout for the given count of symbols.
func (s *Symbol) Payout(count uint8) float64 {
	if count == 0 {
		return 0.0
	}
	max := uint8(len(s.payouts))
	if max == 0 {
		return 0.0
	}
	if count > max {
		count = max
	}
	count--
	return s.payouts[count]
}

// WildFor returns true if the split symbol can substitute for the given symbol.
func (s *Symbol) WildFor(index utils.Index) bool {
	return s.wildFor.Contains(index)
}

// Multiplier returns the multiplier for the symbol.
func (s *Symbol) Multiplier() float64 {
	return s.multiplier
}

// IsEmpty implements the zjson.Encoder interface.
func (s *Symbol) IsEmpty() bool {
	return false
}

// EncodeFields implements the zjson.Encoder interface.
func (s *Symbol) EncodeFields(enc *zjson.Encoder) {
	enc.Uint16Field("id", uint16(s.id))
	enc.StringField("name", s.name)
	enc.StringField("resource", s.resource)
	enc.StringField("kind", s.kind.String())
	enc.Uint16FieldOpt("morphInto", uint16(s.morphInto))
	enc.FloatFieldOpt("multiplier", s.multiplier, 'f', 2)
	enc.IntBoolFieldOpt("sticky", s.sticky)
	enc.IntBoolFieldOpt("isSplit", s.isSplit)
	enc.IntBoolFieldOpt("isPrize", s.isPrize)
	enc.IntBoolFieldOpt("varyMultiplier", s.varyMultiplier)

	if len(s.wildFor) > 0 {
		enc.StartArrayField("wildFor")
		for ix := range s.wildFor {
			enc.Uint64(uint64(s.wildFor[ix]))
		}
	}

	if len(s.weights) > 0 {
		enc.StartArrayField("weights")
		for ix := range s.weights {
			enc.Float(s.weights[ix], 'f', 2)
		}
	}

	if len(s.payouts) > 0 {
		enc.StartArrayField("payouts")
		for ix := range s.payouts {
			enc.Float(s.payouts[ix], 'f', 2)
		}
	}

	if len(s.scatterPayouts) > 0 {
		enc.StartArrayField("scatterPayouts")
		for ix := range s.scatterPayouts {
			enc.Float(s.scatterPayouts[ix], 'f', 2)
		}
	}

	if len(s.clearPattern) > 0 {
		enc.StartArrayField("clearPattern")
		for ix := range s.clearPattern {
			enc.Int64(int64(s.clearPattern[ix]))
		}
	}
}

func (s *Symbol) init() *Symbol {
	s.isSplit = s.kind.isSplit()
	s.isWild = s.kind.isWild()
	s.isScatter = s.kind.isScatter()
	s.isHero = s.kind.isHero()
	s.isBomb = s.kind.isBomb()
	s.isShooter = s.kind.isShooter()
	s.isPrize = s.kind.isPrize()

	s.minPayable = math.MaxUint8
	for ix := range s.payouts {
		if s.payouts[ix] > 0 {
			s.minPayable = uint8(ix + 1)
			break
		}
	}

	return s
}

// SymbolOption is the function signature for Symbol option functions.
type SymbolOption = func(*Symbol)

// WithName is the symbol option to set the symbol name.
func WithName(name string) SymbolOption {
	return func(s *Symbol) {
		s.name = name
	}
}

// WithResource is the symbol option to set the symbol resource.
func WithResource(resource string) SymbolOption {
	return func(s *Symbol) {
		s.resource = resource
	}
}

// WithKind is the symbol option to set the symbol kind.
func WithKind(kind SymbolKind) SymbolOption {
	return func(s *Symbol) {
		s.kind = kind
	}
}

// IsSticky is the symbol option to set the sticky flag.
func IsSticky() SymbolOption {
	return func(s *Symbol) {
		s.sticky = true
	}
}

// WithWeights is the symbol option to set the symbol weights.
func WithWeights(weights ...float64) SymbolOption {
	return func(s *Symbol) {
		s.weights = append(s.weights, weights...)
	}
}

// WithPayouts is the symbol option to set the symbol paylines paytable.
func WithPayouts(payouts ...float64) SymbolOption {
	return func(s *Symbol) {
		s.payouts = append(s.payouts, payouts...)
	}
}

// WithScatterPayouts is the symbol option to set the symbol scatter paytable.
func WithScatterPayouts(payouts ...float64) SymbolOption {
	return func(s *Symbol) {
		s.scatterPayouts = append(s.scatterPayouts, payouts...)
	}
}

// WildFor is the symbol option to set the symbol indexes this symbol can substitute for.
// e.g. it is a "semi-wild" symbol that acts as a wild for the given symbol indexes only.
func WildFor(subs ...utils.Index) SymbolOption {
	return func(s *Symbol) {
		s.wildFor = append(s.wildFor, subs...)
		s.kind = Split
	}
}

// WithMultiplier is the symbol option to set the symbol multiplier.
func WithMultiplier(multiplier float64) SymbolOption {
	return func(s *Symbol) {
		s.multiplier = multiplier
	}
}

// VaryMultiplier is the symbol option to indicate the symbol can contain an arbitrary multiplier.
// This feature allows games to set and increase the multiplier for the symbol.
// The symbol instance itself is immutable, so the value must be stored in the game spin state.
// When the symbol lands in the grid, its fixed multiplier value can be used as the initializer for the mutable multiplier in the game state.
func VaryMultiplier() SymbolOption {
	return func(s *Symbol) {
		s.varyMultiplier = true
	}
}

// ClearPattern is the symbol option to set the pattern for a bomb symbol "explosion".
// E.g. horizontal only: {-rowCount, rowCount}, vertical only: {-1, +1},
// horizontal & vertical: {-rowCount, -1, 1, rowCount},
// 3 by 3 pattern: {-rowCount-1, -rowCount, -rowCount+1, -1, 1, rowCount-1, rowCount, rowCount+1}
// The position of the bomb itself is always cleared, so does not need to be specified in the pattern.
// Any positions that fall outside the grid are ignored.
// Note that the clearing positions cannot exceed the rowCount; so for rowCount==3, the maximum horizontal clearance
// is 1 row up or down. For rowCount==5, the maximum clearance is 2 rows up/down.
func ClearPattern(pattern ...int) SymbolOption {
	return func(s *Symbol) {
		s.clearPattern = append(s.clearPattern, pattern...)
	}
}

// MorphInto indicates that when this symbol lands in the grid, it may morph into the given symbol.
func MorphInto(symbol utils.Index) SymbolOption {
	return func(s *Symbol) {
		s.morphInto = symbol
	}
}

// Symbols is a convenience type for a slice of symbols.
type Symbols []*Symbol
