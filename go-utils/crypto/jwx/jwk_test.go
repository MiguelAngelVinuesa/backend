package jwx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const cognito = "https://cognito-idp.eu-central-1.amazonaws.com/eu-central-1_RzO6963xD/.well-known/jwks.json"

func TestRegister(t *testing.T) {
	t.Run("register JWK", func(t *testing.T) {
		err := RegisterJWK(cognito)
		require.Nil(t, err)

		set, err2 := GetJWK(cognito)
		require.Nil(t, err2)
		require.NotNil(t, set)
	})
}
