package handlers

import (
	"football-team-management/test"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingHandlerImpl_Handle(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		handler := NewPingHandlerImpl()
		router := test.Router("/api/ping", handler.Handle, http.MethodGet)

		response := test.MakeRequest(router, http.MethodGet, "/api/ping", nil)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "pong", response.Body.String())
	})
}
