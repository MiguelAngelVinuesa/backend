package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// PayDirection indicates the direction(s) in which paylines will warrant a payout.
type PayDirection uint8

const (
	// PayLTR means left-to-right (the default).
	PayLTR PayDirection = iota + 1
	// PayRTL means right-to-left.
	PayRTL
	// PayBoth means both ways (LTR and RTL).
	PayBoth
	PayCluster
	PayScatter
)

// String implements the Stringer interface.
func (d PayDirection) String() string {
	switch d {
	case PayLTR:
		return "ltr"
	case PayRTL:
		return "rtl"
	case PayBoth:
		return "both"
	case PayCluster:
		return "cluster"
	case PayScatter:
		return "scatter"
	default:
		return "[unknown]"
	}
}

// Payline represents a payline for a slot machine.
type Payline struct {
	id      uint8
	rows    utils.UInt8s
	offsets []int
}

// NewPayline instantiates a new payline.
// id should be > 0, as the zero id is used for the "AllPaylines" feature.
func NewPayline(id, rowCount uint8, rows ...uint8) *Payline {
	p := &Payline{
		id:   id,
		rows: rows,
	}

	max := len(rows)
	p.offsets = make([]int, max)

	var offset int
	for reel, row := range rows {
		p.offsets[reel] = offset + int(row)
		offset += int(rowCount)
	}

	return p
}

// ID returns the id of the payline.
func (p *Payline) ID() uint8 {
	return p.id
}

// RowMap returns the row map of the payline.
func (p *Payline) RowMap() utils.UInt8s {
	return p.rows
}

// IsEmpty implements the zjson.Encoder interface.
func (p *Payline) IsEmpty() bool {
	return false
}

// EncodeFields implements the zjson.Encoder interface.
func (p *Payline) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("id", p.id)
	enc.StartArrayField("rows")
	for ix := range p.rows {
		enc.Uint64(uint64(p.rows[ix]))
	}
	enc.EndArray()
	enc.StartArrayField("offsets")
	for ix := range p.offsets {
		enc.Uint64(uint64(p.offsets[ix]))
	}
	enc.EndArray()
}

// Paylines is a slice of paylines.
type Paylines []*Payline

// PurgePaylines returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgePaylines(input Paylines, capacity int) Paylines {
	if cap(input) < capacity {
		return make(Paylines, 0, capacity)
	}
	return input[:0]
}
