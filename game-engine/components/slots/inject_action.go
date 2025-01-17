package slots

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// InjectAction represents the action to inject a symbol or set of symbols.
// This operation may happen when the grid is already displayed, and the action will be animated in the FE.
type InjectAction struct {
	SpinAction
	single            bool
	singleFromEdge    bool
	singleEdgeSteps   uint8
	singleReels       utils.UInt8s
	singleMultipliers utils.WeightedGenerator
	cluster           bool
	clusterStart      uint8
	clusterMin        uint8
	clusterMax        uint8
	clusterOffsets    utils.UInt8s
}

// NewSymbolInject instantiates an inject action for the given symbol.
// Where the symbol can be injected may be limited by the given reels.
// Note that reels are 1-based!
func NewSymbolInject(symbol utils.Index, reels ...uint8) *InjectAction {
	a := &InjectAction{}
	a.init(Injection, SymbolsInjected, reflect.TypeOf(a).String())
	a.single = true
	a.symbol = symbol
	a.singleReels = reels
	return a.finalizer()
}

// NewSymbolInjectFromEdge instantiates an inject action for the given symbol.
// The insert can only happen on tiles that are the given number of steps from the edge of the grid.
// Note that tiles on the edge are 1 step from the edge!
func NewSymbolInjectFromEdge(symbol utils.Index, edgeSteps uint8) *InjectAction {
	a := &InjectAction{}
	a.init(Injection, SymbolsInjected, reflect.TypeOf(a).String())
	a.single = true
	a.singleFromEdge = true
	a.symbol = symbol
	a.singleEdgeSteps = edgeSteps
	return a.finalizer()
}

// NewClusterSymbolInject instantiates an action to inject a cluster with the given symbol.
// The starting position from where the cluster can be injected may be limited by the given offsets.
// The function panics if max < min.
func NewClusterSymbolInject(symbol utils.Index, min, max uint8, offsets ...uint8) *InjectAction {
	if max < min {
		panic("invalid cluster injection parameters")
	}

	a := &InjectAction{}
	a.init(Injection, SymbolsInjected, reflect.TypeOf(a).String())
	a.cluster = true
	a.symbol = symbol
	a.clusterMin = min
	a.clusterMax = max
	a.clusterOffsets = offsets
	return a.finalizer()
}

// NewClusterOffsetInject instantiates an action to inject a cluster with the given symbol.
// The cluster is injected starting with the given offset.
// The function panics if max < min.
func NewClusterOffsetInject(min, max uint8, offset uint8) *InjectAction {
	if max < min {
		panic("invalid cluster injection parameters")
	}

	a := &InjectAction{}
	a.init(Injection, SymbolsInjected, reflect.TypeOf(a).String())
	a.cluster = true
	a.clusterStart = offset
	a.clusterMin = min
	a.clusterMax = max
	return a.finalizer()
}

// WithMultipliers adds an array of weighted multipliers to choose from randomly.
// The function panics if the action is not marked to inject a single symbol.
func (a *InjectAction) WithMultipliers(multipliers utils.WeightedGenerator) *InjectAction {
	if !a.single {
		panic("invalid multipliers for inject action")
	}
	a.singleMultipliers = multipliers
	return a.finalizer()
}

// Triggered implements the SpinActioner.Triggered interface.
func (a *InjectAction) Triggered(spin *Spin) SpinActioner {
	switch {
	case a.single && a.singleFromEdge:
		return a.injectFromEdge(spin)
	case a.single:
		return a.injectSymbol(spin)
	case a.cluster && a.clusterStart > 0:
		return a.injectClusterAroundOffset(spin, a.clusterStart)
	case a.cluster:
		return a.injectClusterAtRandom(spin)
	}
	return nil
}

func (a *InjectAction) injectFromEdge(spin *Spin) SpinActioner {
	def, prng := spin.gridDef, spin.prng

	options := make(utils.UInt8s, 0, 100)
	for offset := range def.stepsOffGrid {
		if def.stepsOffGrid[offset] == a.singleEdgeSteps {
			options = append(options, uint8(offset))
		}
	}

	l := len(options)
	if l == 0 {
		return nil
	}

	var deadlock int
	for {
		ix := prng.IntN(l)
		offset := options[ix]
		if !spin.sticky[offset] {
			return a.injectSymbolAtOffset(spin, int(offset), a.symbol)
		}
		if deadlock++; deadlock >= 3*l {
			return nil
		}
	}
}

func (a *InjectAction) injectSymbol(spin *Spin) SpinActioner {
	def := spin.gridDef
	l := def.gridSize
	offset := spin.prng.IntN(l)

	valid := func() bool {
		return !spin.sticky[offset]
	}

	if len(a.singleReels) > 0 {
		valid = func() bool {
			if spin.sticky[offset] {
				return false
			}

			reel, _ := def.ReelRowFromOffset(offset)
			for ix := range a.singleReels {
				if reel == int(a.singleReels[ix])-1 {
					return true
				}
			}
			return false
		}
	}

	var deadlock int
	for {
		if def.IsValidOffset(offset) && valid() {
			return a.injectSymbolAtOffset(spin, offset, a.symbol)
		}

		if deadlock++; deadlock >= l {
			return nil
		}

		offset += def.offsetStep
		if offset >= l {
			offset -= l
		}
	}
}

func (a *InjectAction) injectSymbolAtOffset(spin *Spin, offset int, symbol utils.Index) SpinActioner {
	spin.indexes[offset] = symbol
	spin.injections[offset] = symbol

	if a.singleMultipliers != nil {
		m := a.singleMultipliers.RandomIndex(spin.prng)
		if utils.ValidMultiplier(float64(m)) {
			if l := len(spin.indexes); len(spin.multipliers) != l {
				spin.multipliers = utils.PurgeUInt16s(spin.multipliers, l)[:l]
			}
			spin.multipliers[offset] = uint16(m)
		}
	}

	return a
}

func (a *InjectAction) injectClusterAtRandom(spin *Spin) SpinActioner {
	if l := len(a.clusterOffsets); l > 0 {
		ix := spin.prng.IntN(l)
		return a.injectClusterAroundOffset(spin, a.clusterOffsets[ix])
	}

	offset := spin.prng.IntN(spin.gridDef.gridSize)
	return a.injectClusterAroundOffset(spin, uint8(offset))
}

func (a *InjectAction) injectClusterAroundOffset(spin *Spin, offset uint8) SpinActioner {
	// determine the cluster size to inject.
	count := int(a.clusterMin)
	if a.clusterMax > a.clusterMin {
		count += spin.prng.IntN(int(a.clusterMax - a.clusterMin + 1))
	}

	def := spin.gridDef
	added := make([]uint8, 0, 100)

	// check if the center is a wild!
	symbols := spin.GetSymbols()
	symbolID := spin.indexes[offset]
	symbol := symbols.GetSymbol(symbolID)

	for symbol == nil || symbol.kind != Standard {
		// if the offset is occupied by a special symbol, we need to determine what symbol to use for the cluster!
		ix := count
		if max := len(symbols.bestWildSym); ix >= max {
			ix = max - 1
		}
		symbolID = symbols.bestWildSym[ix]

		// we need a new anchor point to start building the cluster!
		// we'll iterate thropugh all neighbors until we find a suitable candidate.
		// it must have the symbol we're looking for, or some other standard symbol.
		// if we can't find a suitable neighbor, we'll fall back into picking a random tile.
		var deadlock int
		for symbol == nil || symbol.kind != Standard {
			list := def.withoutSelf[offset]
			offset = list.RandomNeighbor(spin.prng)
			if !spin.sticky[offset] {
				if spin.indexes[offset] == symbolID {
					break
				}
				if symbol = symbols.GetSymbol(spin.indexes[offset]); symbol != nil && symbol.kind == Standard {
					break
				}
			}

			if deadlock++; deadlock >= 10 {
				// we hit the brick wall!
				// there's no valid neighboring tiles, so let's just find a random tile!
				// it must have the symbol we're looking for, or some other standard symbol.
				// if we can't find a random tile, we've run out of options, and cannot build a cluster!
				offset = uint8(spin.prng.IntN(def.gridSize))
				step := uint8(def.offsetStep)
				size := uint8(def.gridSize)
				for {
					if def.IsValidOffset(int(offset)) {
						if spin.indexes[offset] == symbolID {
							break
						}
						if !spin.sticky[offset] {
							if symbol = symbols.GetSymbol(spin.indexes[offset]); symbol != nil && symbol.kind == Standard {
								break
							}
						}
					}

					// we've optimized for speed here. we bail out if too many attempts have failed.
					// the solutiuon may not hit every tile in the grid, but it's the best effort.
					// if there are such few options available, it's unlikely we can create a valid cluster anyway!
					if deadlock++; deadlock >= 50 {
						return nil
					}

					// round-robin through the grid with the chosen optimal prime from the grid definition.
					if offset += step; offset > size {
						offset -= size
					}
				}
			}
		}

		if spin.indexes[offset] != symbolID {
			spin.indexes[offset] = symbolID
			spin.injections[offset] = symbolID
			added = append(added, offset)
			count--
		}
	}

	// keep going until we hit our cluster count, or there are no more possibilities.
	var injected int
	for count > 0 {
		// inject the selected symbol in every neighboring tile where possible.
		// if we hit the requested symbol, we can keep its offset as another anchor point.
		// if the counter hits zero, we're also done.
		list := def.withoutSelf[offset]
		for offset = range list {
			if !spin.sticky[offset] {
				if spin.indexes[offset] == symbolID {
					added = append(added, offset)
				} else {
					spin.indexes[offset] = symbolID
					spin.injections[offset] = symbolID
					added = append(added, offset)
					injected++
					if count--; count <= 0 {
						break
					}
				}
			}
		}

		// continue from another anchor point if needed.
		// if there aren't any, we've run out of options, so we've got to bail!
		if count > 0 {
			if len(added) > 0 {
				offset = added[0]
				added = added[1:]
			} else {
				if injected > 0 {
					return a
				}
				return nil
			}
		}
	}

	return a
}

func (a *InjectAction) finalizer() *InjectAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	switch {
	case a.single:
		if a.singleFromEdge {
			b.WriteString(",singleFromEdge=true")
		} else {
			b.WriteString(",single=true")
		}
		b.WriteString(",symbol=")
		b.WriteString(strconv.Itoa(int(a.symbol)))
		if a.singleFromEdge {
			b.WriteString(",steps=")
			b.WriteString(strconv.Itoa(int(a.singleEdgeSteps)))
		} else {
			if len(a.singleReels) > 0 {
				b.WriteString(",reels=")
				j, _ := json.Marshal(a.singleReels)
				b.Write(j)
			}
		}

	case a.cluster:
		if a.clusterStart > 0 {
			b.WriteString(",clusterOffset=true")
			b.WriteString(",start=")
			b.WriteString(strconv.Itoa(int(a.clusterStart)))
		} else {
			b.WriteString(",cluster=true")
		}
		b.WriteString(",symbol=")
		b.WriteString(strconv.Itoa(int(a.symbol)))
		b.WriteString(",min=")
		b.WriteString(strconv.Itoa(int(a.clusterMin)))
		b.WriteString(",max=")
		b.WriteString(strconv.Itoa(int(a.clusterMax)))
		if len(a.clusterOffsets) > 0 {
			b.WriteString(",offsets=")
			j, _ := json.Marshal(a.clusterOffsets)
			b.Write(j)
		}
	}

	a.config = b.String()
	return a
}
