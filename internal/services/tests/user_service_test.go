package tests

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "zesty-sips-api/internal/services"
)

func TestCreateUser(t *testing.T) {
    userService := services.NewUserService()
    user, err := userService.CreateUser("test@example.com", "password")
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
