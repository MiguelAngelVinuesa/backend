package cards

import "sync"

// Table represents a table that the players of a game are placed at.
// It is safe to use across multiple go-routines simultaneously.
type Table struct {
	mutex   sync.RWMutex
	deck    *Deck
	players []string
	hands   []*Hand
	playing bool
}

var tablePool = sync.Pool{
	New: func() interface{} {
		return &Table{players: make([]string, 0)}
	},
}

// NewTable instantiates a new table from the memory pool.
func NewTable(deck *Deck, players ...string) *Table {
	t := tablePool.Get().(*Table)
	t.deck = deck
	t.players = append(t.players, players...)
	for range players {
		t.hands = append(t.hands, NewHand())
	}
	return t
}

// Release puts a table back into the memory pool.
func (t *Table) Release() {
	if t != nil {
		for _, h := range t.hands {
			h.Release()
		}
		t.deck = nil
		t.players = t.players[:0]
		t.hands = t.hands[:0]
		t.playing = false
		tablePool.Put(t)
	}
}

// AddPlayers adds one or more players to the table.
// When adding players they are placed at the next available seat.
// The first player is assigned seat 0, the second seat 1, and so on.
// Players can only be added when the game is not in progress!
func (t *Table) AddPlayers(players ...string) *Table {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.playing {
		return t
	}

	t.players = append(t.players, players...)
	for range players {
		t.hands = append(t.hands, NewHand())
	}
	return t
}

// IsPlaying returns whether the game at the table is in progress.
func (t *Table) IsPlaying() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.playing
}

// PlayerCount returns the number of players at the table.
func (t *Table) PlayerCount() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return len(t.players)
}

// Player returns a player's current hand based on their seat at the table (0-based).
func (t *Table) Player(seat int) *Hand {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if seat < 0 || seat >= len(t.hands) {
		return nil
	}
	return t.hands[seat]
}

// PlayerByID returns a player's seat and current hand based on their ID.
// If the ID is not found, it returns the invalid seat (-1) and a nil hand.
func (t *Table) PlayerByID(player string) (int, *Hand) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	for ix, p := range t.players {
		if p == player {
			return ix, t.hands[ix]
		}
	}
	return -1, nil
}

// Remaining returns the number of remaining cards in the deck.
func (t *Table) Remaining() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.deck.Remaining()
}

// Draw draws count cards into the hand of each player at the table.
// If count is less than one it is adjusted to one.
func (t *Table) Draw(count int) *Table {
	if count < 1 {
		count = 1
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.playing = true

	for ix := 0; ix < count; ix++ {
		for _, hand := range t.hands {
			c := t.deck.Draw()
			if c == nil {
				return t
			}
			hand.Add(c)
		}
	}

	return t
}

// PlayerDraw draws count cards into the hand of the player based on their seat at the table (0-based).
// If seat contains an invalid value, nothing happens.
// If count is less than one it is adjusted to one.
func (t *Table) PlayerDraw(seat, count int) *Table {
	if count < 1 {
		count = 1
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	if seat < 0 || seat >= len(t.hands) {
		return t
	}

	t.playing = true

	for iy := 0; iy < count; iy++ {
		c := t.deck.Draw()
		if c == nil {
			return t
		}
		t.hands[seat].Add(c)
	}

	return t
}

// Burn burns one card from the deck.
func (t *Table) Burn() *Table {
	t.mutex.Lock()
	t.deck.Burn()
	t.playing = true
	t.mutex.Unlock()
	return t
}

// Shuffle resets the deck and shuffles it.
// Hands are also cleared.
func (t *Table) Shuffle() *Table {
	t.mutex.Lock()
	for _, hand := range t.hands {
		hand.Clear()
	}
	t.deck.Shuffle()
	t.playing = false
	t.mutex.Unlock()
	return t
}
