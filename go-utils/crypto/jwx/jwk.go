package jwx

import (
	"context"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

// RegisterJWK registers a JWK url in the JWK cache and refreshes it.
func RegisterJWK(url string) error {
	if err := cache.Register(url, jwk.WithMinRefreshInterval(15*time.Minute)); err != nil {
		return err
	}
	_, err := cache.Refresh(ctx, url)
	return err
}

// GetJWK retrieves a cached JSK.
func GetJWK(url string) (jwk.Set, error) {
	return cache.Get(ctx, url)
}

// cache represents a JWK cache.
var (
	ctx   = context.Background()
	cache = jwk.NewCache(ctx)
)
