package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// ConnectionKind is the kind of grid connections.
type ConnectionKind uint8

const (
	// Rectangular indicates the grid connections are horizontal and vertical only.
	Rectangular ConnectionKind = iota
	// Hexagonal indicates grid connections are hexagonal.
	Hexagonal
)

// ClusterPayouts contains the details for determining cluster payouts.
type ClusterPayouts struct {
	reels       uint8
	rows        uint8
	mask        []uint8
	kind        ConnectionKind
	connections []utils.UInt8s
}

// NewClusterPayouts instantiates a cluster payouts.
func NewClusterPayouts(reels, rows uint8, opts ...ClusterPayoutsOption) *ClusterPayouts {
	n := &ClusterPayouts{
		reels: reels,
		rows:  rows,
		kind:  Rectangular,
	}

	for ix := range opts {
		opts[ix](n)
	}

	if len(n.connections) == 0 {
		n.connections = PurgeConnections(n.connections, int(reels)*int(rows))
		if len(n.mask) == 0 || n.kind != Hexagonal {
			n.connectRectangular()
		} else {
			n.connectHexagonal()
		}
	}

	return n
}

// RemovePayouts removes all potential cluster payouts from the grid.
func (c *ClusterPayouts) RemovePayouts(spin *Spin) {
	n := c.Find(spin, nil)
	if n == 0 {
		return
	}

	symbols := spin.GetSymbols()
	maxID := int(symbols.maxID)
	var f float64
	for f < 0.01 {
		s := symbols.GetSymbol(utils.Index(maxID))
		if len(s.payouts) != 0 && !s.isScatter {
			for _, f = range s.payouts {
				if f > 0.0 {
					break
				}
			}
		}
		if f < 0.01 {
			maxID--
		}
	}

	prng := spin.prng
	reels, rows, mask := spin.reelCount, spin.rowCount, spin.mask

	unwanted := make(utils.Indexes, 0, 8)

	good := func(id utils.Index) bool {
		for _, id2 := range unwanted {
			if id == id2 {
				return false
			}
		}
		return true
	}

	for reel := 0; reel < reels; reel++ {
		rowM := int(mask[reel])
		offs := reel * rows
		for row := 0; row < rowM; row++ {
			conns := c.connections[offs]
			for ix := range conns {
				conn := int(conns[ix])
				if conn > offs && spin.indexes[conn] == spin.indexes[offs] {
					conns2 := c.connections[conn]
					for iy := range conns2 {
						if conn2 := int(conns2[iy]); conn2 < conn {
							unwanted = append(unwanted, spin.indexes[conn2])
						}
					}

					for {
						id := utils.Index(1 + prng.IntN(maxID))
						if s := symbols.GetSymbol(id); s != nil && good(id) {
							spin.indexes[conn] = id
							break
						}
					}

					unwanted = unwanted[:0]
				}
			}
			offs++
		}
	}

}

// Find searches for and reports cluster payouts in the given spin grid.
func (c *ClusterPayouts) Find(spin *Spin, res *results.Result) int {
	//log.Printf("\n\n RESETTING PAYOUTS FROM CLUSTERPAYOUTS.FIND CURRENT %v \n\n", spin.payouts)
	spin.resetPayouts()

	// fixed slices so they get allocated on the heap.
	// NOTE: If we get grids with a higher count this may need some rethinking (200 is ok too...).
	path := make([]uint8, 0, 100)
	payouts := make([]bool, 0, 100)[:spin.gridDef.gridSize]

	// take note of all the wilds!
	wilds := make([]bool, 0, 100)[:spin.gridDef.gridSize]
	for ix, id := range spin.indexes {
		if id > 0 {
			if symbol := spin.GetSymbol(id); symbol != nil && symbol.isWild {
				wilds[ix] = true
			}
		}
	}

	// the minimum count for payouts will help us speed things up.
	minimum := spin.symbols.minPayout

	// find the highest start offset that can still warrant a payout.
	maxStart := len(c.connections) - int(minimum) + int(spin.mask[c.reels-1]) - int(c.rows)

	// keep track of multiplier.
	haveMultipliers := len(spin.multipliers) == len(spin.indexes)
	var multiplier float64

	// we need some recursive power!
	var findPath func(uint8, uint8, utils.Index)
	findPath = func(start, current uint8, symbolID utils.Index) {
		for _, offset := range c.connections[current] {
			if payouts[offset] || (spin.payouts[offset] == 1 && !wilds[offset]) {
				continue
			}

			if spin.payouts[offset] == 0 && spin.indexes[offset] != symbolID && !wilds[offset] {
				continue
			}

			path = append(path, offset)
			payouts[offset] = true
			spin.payouts[offset]++

			if haveMultipliers {
				if m := float64(spin.multipliers[offset]); utils.ValidMultiplier(m) {
					multiplier *= m
				}
			}

			findPath(start, offset, symbolID)
		}
	}

	var payoutCount int
	for startOffset := range c.connections {
		// we only need to consider tiles that can still match the minimum count, that are not yet part of a payout,
		// and are connected to other tiles (e.g. non-rectangular grids have dangling tiles).
		if startOffset <= maxStart && spin.payouts[startOffset] == 0 && len(c.connections[startOffset]) > 0 {
			multiplier = 1.0
			symbolID := spin.indexes[startOffset]
			if symbol := spin.symbols.GetSymbol(symbolID); symbol != nil && !symbol.isWild {
				path = append(path, uint8(startOffset))
				payouts[startOffset] = true
				spin.payouts[startOffset] = 1

				// find all connections and record the payout if there are enough.
				findPath(uint8(startOffset), uint8(startOffset), symbolID)
				var p results.Payout
				if count := uint8(len(path)); count >= minimum {
					if count >= symbol.minPayable {
						p = ClusterPayout(symbol.Payout(count), utils.NewMultiplier(multiplier, spin.getMultiplier(0)), symbolID, count, path)
						payoutCount++
					}
				}

				if p == nil {
					for _, offset := range path {
						spin.payouts[offset]--
					}
				} else if res == nil {
					p.Release()
				} else {
					res.AddPayouts(p)
				}

				path = path[:0]
				clear(payouts)
			}
		}
	}
	return payoutCount
}

// connectRectangular calculates the grid connections for rectangular grids.
func (c *ClusterPayouts) connectRectangular() {
	reels, rows := c.reels, c.rows
	c.connections = c.connections[:reels*rows]

	var offset uint8
	for reel := uint8(0); reel < reels; reel++ {
		for row := uint8(0); row < rows; row++ {
			conn := utils.PurgeUInt8s(c.connections[offset], 6)

			if row > 0 {
				conn = append(conn, offset-1)
			}
			if reel > 0 {
				conn = append(conn, offset-rows)
			}
			if row < rows-1 {
				conn = append(conn, offset+1)
			}
			if reel < reels-1 {
				conn = append(conn, offset+rows)
			}

			c.connections[offset] = conn
			offset++
		}
	}
}

// connectHexagonal calculates the grid connections for hexagonal grids.
func (c *ClusterPayouts) connectHexagonal() {
	reels, rows := c.reels, c.rows
	c.connections = c.connections[:reels*rows]

	var offset uint8
	for reel := uint8(0); reel < reels; reel++ {
		currRows := c.mask[reel]
		for row := uint8(0); row < rows; row++ {
			conn := utils.PurgeUInt8s(c.connections[offset], 6)

			if row < currRows {
				if row > 0 {
					conn = append(conn, offset-1)
				}
				if row < currRows-1 {
					conn = append(conn, offset+1)
				}

				if reel > 0 {
					prevRows := c.mask[reel-1]
					if currRows > prevRows {
						if row > 0 {
							conn = append(conn, offset-rows-1)
						}
						if row < prevRows {
							conn = append(conn, offset-rows)
						}
					} else {
						if row < prevRows-1 {
							conn = append(conn, offset-rows)
						}
						conn = append(conn, offset-rows+1)
					}
				}

				if reel < reels-1 {
					nextRows := c.mask[reel+1]
					if currRows < nextRows {
						conn = append(conn, offset+rows)
						if row < nextRows-1 {
							conn = append(conn, offset+rows+1)
						}
					} else {
						if row > 0 {
							conn = append(conn, offset+rows-1)
						}
						if row < nextRows {
							conn = append(conn, offset+rows)
						}
					}
				}
			}

			c.connections[offset] = conn
			offset++
		}
	}
}

// ClusterPayoutsOption is the function signature for cluster payouts options.
type ClusterPayoutsOption func(c *ClusterPayouts)

// ClusterGridMask adds a grid mask and connection kind to the cluster payouts.
// Make sure to use the same mask as when instantiating the Slots structure!
func ClusterGridMask(mask []uint8, kind ConnectionKind) ClusterPayoutsOption {
	return func(c *ClusterPayouts) {
		c.mask = mask
		c.kind = kind
	}
}

// ClusterConnections fills the grid connections from the supplied slice.
func ClusterConnections(connections []utils.UInt8s) ClusterPayoutsOption {
	return func(c *ClusterPayouts) {
		c.connections = connections
	}
}

// PurgeConnections returns the input to zero length or a new slice if its capacity is less than requested.
func PurgeConnections(in []utils.UInt8s, capacity int) []utils.UInt8s {
	if cap(in) < capacity {
		return make([]utils.UInt8s, 0, capacity)
	}
	return in[:0]
}
