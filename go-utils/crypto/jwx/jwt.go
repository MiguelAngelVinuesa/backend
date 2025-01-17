package jwx

import (
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// ParseJWT parses and verifies a JWT using the given JWK url.
func ParseJWT(data []byte, url string) (jwt.Token, error) {
	token, err := GetJWK(url)
	if err != nil {
		return nil, err
	}

	set, err2 := jwk.PublicSetOf(token)
	if err2 != nil {
		return nil, err2
	}

	return jwt.Parse(data, jwt.WithKeySet(set))
}
