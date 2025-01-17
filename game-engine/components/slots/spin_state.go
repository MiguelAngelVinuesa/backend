package slots

import (
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// AcquireSpinState instantiates a new spin state from the memory pool.
func AcquireSpinState(spin *Spin) *SpinState {
	s := spinStateProducer.Acquire().(*SpinState)
	if spin == nil {
		return s
	}

	l := len(spin.indexes)
	s.stickySymbol = spin.stickySymbol
	s.freeSpins = spin.freeSpins
	s.indexes = utils.PurgeIndexes(s.indexes, l)[:l]
	s.sticky = utils.PurgeBools(s.sticky, l)[:l]

	copy(s.indexes, spin.indexes)
	copy(s.sticky, spin.sticky)

	return s
}

// FreeSpins returns the number of free spins.
func (s *SpinState) FreeSpins() uint64 {
	return s.freeSpins
}

// SetIndexes updates the saved grid with the given slice and resets the sticky flags.
func (s *SpinState) SetIndexes(indexes utils.Indexes) {
	copy(s.indexes, indexes)
	for ix := range s.sticky {
		s.sticky[ix] = false
	}
}

// SetStickySymbol updates the sticky state using a player selected symbol.
func (s *SpinState) SetStickySymbol(symbol utils.Index) {
	if symbol == s.stickySymbol {
		return
	}

	s.stickySymbol = symbol
	for ix, id := range s.indexes {
		s.sticky[ix] = id == symbol
	}
}

// StartGrid returns the start grid used to perform the spin.
func (s *SpinState) StartGrid() int {
	return s.startGrid
}

// SetStartGrid sets the start grid used to perform the spin.
func (s *SpinState) SetStartGrid(grid int) {
	s.startGrid = grid
}

// EncodeFields implements the zjson encoder interface.
func (s *SpinState) EncodeFields(enc *zjson.Encoder) {
	enc.Uint16FieldOpt("stickySymbol", uint16(s.stickySymbol))
	enc.IntFieldOpt("startGrid", s.startGrid)
	enc.Uint64FieldOpt("freeSpins", s.freeSpins)

	enc.StartArrayField("indexes") // can't be empty
	for ix := range s.indexes {
		enc.Uint64(uint64(s.indexes[ix]))
	}
	enc.EndArray()

	if len(s.sticky) > 0 {
		enc.StartArrayField("sticky")
		for ix := range s.sticky {
			enc.IntBool(s.sticky[ix])
		}
		enc.EndArray()
	}
}

// DecodeField implements the zjson decoder interface.
func (s *SpinState) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool
	var i int64
	var u uint64

	if string(key) == "stickySymbol" {
		if i, ok = dec.Int64(); ok && i != 0 && i < math.MaxUint16 {
			s.stickySymbol = utils.Index(i)
		}
	} else if string(key) == "startGrid" {
		if i, ok = dec.Int64(); ok && i > 0 && i < math.MaxInt {
			s.startGrid = int(i)
		}
	} else if string(key) == "freeSpins" {
		if u, ok = dec.Uint64(); ok && i >= 0 && i <= 100 {
			s.freeSpins = u
		}
	} else if string(key) == "indexes" {
		ok = dec.Array(s.decodeIndexes)
	} else if string(key) == "sticky" {
		ok = dec.Array(s.decodeSticky)
	}

	if ok {
		return nil
	}
	return dec.Error()
}

func (s *SpinState) decodeIndexes(dec *zjson.Decoder) error {
	if i, ok := dec.Uint16(); ok {
		s.indexes = append(s.indexes, utils.Index(i))
		return nil
	}
	return dec.Error()
}

func (s *SpinState) decodeSticky(dec *zjson.Decoder) error {
	if i, ok := dec.IntBool(); ok {
		s.sticky = append(s.sticky, i)
		return nil
	}
	return dec.Error()
}

// SpinState is used to keep the state of a spin alive across the session.
// Keep fields ordered by ascending SizeOf().
type SpinState struct {
	stickySymbol utils.Index
	startGrid    int
	freeSpins    uint64
	indexes      utils.Indexes
	sticky       []bool
	pool.Object
}

// spinStateProducer is the memory pool for spin states.
// Make sure to initialize all slices appropriately!
var spinStateProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	s := &SpinState{
		stickySymbol: utils.MaxIndex,
		indexes:      make(utils.Indexes, 0, 16),
		sticky:       make([]bool, 0, 16),
	}
	return s, s.reset
})

// reset clears the spin state.
func (s *SpinState) reset() {
	if s != nil {
		s.stickySymbol = utils.MaxIndex
		s.startGrid = 0
		s.freeSpins = 0
		s.indexes = s.indexes[:0]
		s.sticky = s.sticky[:0]
	}
}
