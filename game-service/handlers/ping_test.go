package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	t.Run("ping", func(t *testing.T) {
		app := fiber.New()
		require.NotNil(t, app)

		app.Get("/v1/ping", Ping)

		req := httptest.NewRequest(fiber.MethodGet, "/v1/ping", nil)
		require.NotNil(t, req)

		resp, err := app.Test(req, 100)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}
