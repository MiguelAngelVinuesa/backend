package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestNewGameState(t *testing.T) {
	t.Run("new game state", func(t *testing.T) {
		symbols := slots.NewSymbolSet(slots.NewSymbol(1), slots.NewSymbol(2), slots.NewSymbol(3), slots.NewSymbol(4))
		ss := slots.AcquireSymbolsState(symbols)

		gs := AcquireGameState(nil, ss, 150)
		require.NotNil(t, gs)
		defer gs.Release()

		enc := zjson.AcquireEncoder(256)
		defer enc.Release()

		enc.Object(gs)
		want := `{"roundSeq":1,"bet":150,"symbols":{"flagged":[0,0,0,0,0],"valid":[0,1,1,1,1]}}`
		assert.Equal(t, want, string(enc.Bytes()))
	})
}
