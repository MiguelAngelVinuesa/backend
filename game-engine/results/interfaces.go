package results

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// BonusRunner is the interface to be implemented to allow a game to be used as a bonus game inside another game.
type BonusRunner interface {
	pool.Objecter

	// RequireParams must indicate if the bonus game requires any input from the player.
	// If the return value is true, the caller must wait for user input from the front-end, before calling Run().
	RequireParams() bool

	// Run must execute the bonus game and return the game result.
	// The input result is the result from the initial game that warranted the bonus game.
	// The params are the parameters, if any, as input by the player.
	// The caller is responsible for returning both the original result and the new result to the memory pool after use.
	Run(input *Result, params ...interface{}) (int, *Result)
}
