package pool

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProducerError(t *testing.T) {
	t.Run("new producer error", func(t *testing.T) {
		defer func() {
			e := recover()
			require.NotNil(t, e)
		}()

		p := NewProducer(func() (Objecter, func()) { return nil, nil })
		require.Nil(t, p)
	})
}
