package slots

import (
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// AcquireTile instantiates a new grid tile.
func AcquireTile(offset uint8, symbol utils.Index, sticky uint8, multiplier uint16) *Tile {
	t := tileProducer.Acquire().(*Tile)
	t.offset = offset
	t.symbol = symbol
	t.sticky = sticky
	t.multiplier = multiplier
	return t
}

// AcquireJumpedTile instantiates a new grid tile.
func AcquireJumpedTile(from, to uint8, symbol utils.Index, sticky uint8, multiplier uint16) *Tile {
	t := tileProducer.Acquire().(*Tile)
	t.jumped = true
	t.offsetFrom = from
	t.offset = to
	t.symbol = symbol
	t.sticky = sticky
	t.multiplier = multiplier
	return t
}

// AcquireTileFromJSON instantiates a new grid tile from the given json data.
func AcquireTileFromJSON(dec *zjson.Decoder) (*Tile, bool) {
	t := tileProducer.Acquire().(*Tile)
	if dec.Object(t) {
		return t, true
	}
	t.Release()
	return nil, false
}

// IsJump returns true if the tile represents a "jumped" symbol.
func (t *Tile) IsJump() bool {
	return t.jumped
}

// Offset returns the grid offset of the tile.
func (t *Tile) Offset() uint8 {
	return t.offset
}

// Jump returns the from and to grid offset for a "jumped" symbol.
func (t *Tile) Jump() (from uint8, to uint8) {
	return t.offsetFrom, t.offset
}

// Symbol returns the symbol on the tile.
func (t *Tile) Symbol() utils.Index {
	return t.symbol
}

// Sticky indicates if the tile is sticky.
func (t *Tile) Sticky() uint8 {
	return t.sticky
}

// Multiplier returns the optional multiplier for the symbol.
func (t *Tile) Multiplier() uint16 {
	return t.multiplier
}

// EncodeFields implements the zjson.Encoder.EncodeFields interface.
func (t *Tile) EncodeFields(enc *zjson.Encoder) {
	if t.jumped {
		enc.Uint8Field("from", t.offsetFrom)
		enc.Uint8Field("to", t.offset)
	} else {
		enc.Uint8Field("offset", t.offset)
	}

	enc.Uint16Field("symbol", uint16(t.symbol))
	enc.Uint8FieldOpt("sticky", t.sticky)
	enc.Uint16FieldOpt("multiplier", t.multiplier)
}

// DecodeField implements the zjson.Decoder.DecodeField interface.
func (t *Tile) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool
	var i16 uint16

	if string(key) == "offset" {
		t.offset, ok = dec.Uint8()
		t.jumped = false
	} else if string(key) == "from" {
		t.offsetFrom, ok = dec.Uint8()
		t.jumped = true
	} else if string(key) == "to" {
		t.offset, ok = dec.Uint8()
		t.jumped = true
	} else if string(key) == "symbol" {
		if i16, ok = dec.Uint16(); ok {
			t.symbol = utils.Index(i16)
		}
	} else if string(key) == "sticky" {
		t.sticky, ok = dec.Uint8()
	} else if string(key) == "multiplier" {
		t.multiplier, ok = dec.Uint16()
	} else {
		return fmt.Errorf("Tile.DecodeField invalid field: %s", string(key))
	}

	if ok {
		return nil
	}
	return dec.Error()
}

// Tile contains the details of a grid tile or a "jumped" symbol.
type Tile struct {
	jumped     bool        // true indicates a "jumped" symbol.
	offset     uint8       // grid offset for the tile.
	offsetFrom uint8       // grid offset the symbol jumped from (if jumped == true).
	sticky     uint8       // sticky indicator.
	symbol     utils.Index // the symbol on the tile.
	multiplier uint16      // optional multiplier.
	pool.Object
}

// tileProducer is the memory pool for tiles.
var tileProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &Tile{}
	return e, e.reset
})

// reset clears the tile.
func (t *Tile) reset() {
	if t != nil {
		t.jumped = false
		t.offset = 0
		t.offsetFrom = 0
		t.sticky = 0
		t.symbol = 0
		t.multiplier = 0
	}
}

// DeepEqual is used internally for unit tests.
func (t *Tile) DeepEqual(other *Tile) bool {
	return t.jumped == other.jumped &&
		t.offset == other.offset &&
		t.offsetFrom == other.offsetFrom &&
		t.symbol == other.symbol &&
		t.sticky == other.sticky &&
		t.multiplier == other.multiplier
}

// Tiles is a convenience type for a slice of tiles.
type Tiles []*Tile

// ReleaseTiles releases the tiles in the slice and returns the emptied slice.
func ReleaseTiles(tiles Tiles) Tiles {
	if tiles == nil {
		return nil
	}
	for ix := range tiles {
		tiles[ix].Release()
		tiles[ix] = nil
	}
	return tiles[:0]
}

// DeepEqual is used internally for unit tests.
func (t Tiles) DeepEqual(other Tiles) bool {
	if len(t) != len(other) {
		return false
	}
	for ix := range t {
		if !t[ix].DeepEqual(other[ix]) {
			return false
		}
	}
	return true
}
