package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
    "zesty-sips-api/internal/api/handlers"
)

func TestGetUser(t *testing.T) {
    req, err := http.NewRequest("GET", "/users/1", nil)
    assert.NoError(t, err)

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handlers.GetUser)

    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
}
