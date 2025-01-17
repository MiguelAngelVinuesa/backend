package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// JumpParams defines the parameters for grid jumping.
type JumpParams struct {
	direction GridDirection // possible jump direction(s).
	onSymbols bool          // allow jumping onto existing symbols.
	clone     bool          // clone symbol instead of jumping.
	refill    bool          // refill original tile of jumping symbol.
	minJump   uint8         // minimum jump size.
	maxJump   uint8         // maximum jump size.
	offGrid   float64       // chance to jump off grid.
}

// GridJumps represents a slice of grid jumps.
type GridJumps []gridJump

// TestGridJump tests if the symbol on the given offset can jump, and if so, adds the jump to the slice.
// The function returns the modified input slice.
func (j GridJumps) TestGridJump(params *JumpParams, offset uint8, spin *Spin) GridJumps {
	n := gridJump{
		from:   offset,
		params: params,
		grid:   spin.gridDef,
		spin:   spin,
	}
	if n.test(j) {
		return append(j, n)
	}
	return j
}

// Jump executes all jumps in the slice.
func (j GridJumps) Jump() bool {
	for ix := range j {
		j[ix].jump()
	}
	return len(j) > 0
}

// offsetUsed returns true if the given offset is already the target for a jump in the slice.
func (j GridJumps) offsetUsed(offset uint8) bool {
	for ix := range j {
		if j[ix].to == offset {
			return true
		}
	}
	return false
}

// gridJump represents a symbol jump within the spin grid or off the grid.
type gridJump struct {
	offGrid bool            // the jump lands outside the grid.
	from    uint8           // the offset to jump from.
	to      uint8           // the offset to jump to.
	dir     GridDirection   // the direction of the jump.
	params  *JumpParams     // jump parameters.
	grid    *GridDefinition // grid definition.
	spin    *Spin           // the actual spin the symbol jump is allocated for.
}

// jump executes the symbol jump.
func (j *gridJump) jump() {
	spin, from := j.spin, j.from
	multipliers := len(spin.multipliers) > 0

	if !j.offGrid {
		to := j.to
		spin.indexes[to] = spin.indexes[from]
		spin.sticky[to] = spin.sticky[from]
		spin.jumps[from] = to + 1
		if multipliers {
			spin.multipliers[to] = spin.multipliers[from]
		}
	} else {
		spin.jumps[from] = 255
	}

	if !j.params.clone {
		spin.indexes[from] = 0
		spin.sticky[from] = false
		if multipliers {
			spin.multipliers[from] = 0
		}

		if j.params.refill {
			reel := int(from) / spin.rowCount
			spin.reels[reel].Spin(spin.prng, spin.indexes[from:from+1])
		}
	}
}

func (j *gridJump) test(jumps GridJumps) bool {
	if j.params.offGrid > 0 && j.grid.IsOnTheEdge(int(j.from)) {
		// may jump off grid.
		if float64(j.spin.prng.IntN(10000)) < j.params.offGrid*100 {
			j.offGrid = true
			return true
		}
	}

	// get list of neighbors.
	offs := make(utils.UInt8s, 0, MaxNeighbors)
	dirs := make([]GridDirection, 0, MaxNeighbors)
	for k, v := range j.grid.neighbors[j.from] {
		if (j.params.onSymbols || j.spin.indexes[k] == utils.NullIndex) &&
			j.directionOK(j.params.direction, v) &&
			!jumps.offsetUsed(k) {
			offs = append(offs, k)
			dirs = append(dirs, v)
		}
	}

	var ix int
	switch len(offs) {
	case 0:
		// nowhere to go.
		return false
	case 1:
		// only 1 option, so easy pickings.
		ix = 0
	default:
		// select a key at random.
		ix = j.spin.prng.IntN(len(offs))
	}

	j.to = offs[ix]
	j.dir = dirs[ix]
	return true
}

func (j *gridJump) directionOK(want, got GridDirection) bool {
	switch want {
	case GridAny:
		return true
	case GridHorizontal:
		if got.IsHorizontal() {
			return true
		}
	case GridVertical:
		if got.IsVertical() {
			return true
		}
	case GridDiagonal:
		if got.IsDiagonal() {
			return true
		}
	default:
		if got == want {
			return true
		}
	}
	return false
}
