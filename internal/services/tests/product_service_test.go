package tests

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "zesty-sips-api/internal/services"
)

func TestCreateProduct(t *testing.T) {
    productService := services.NewProductService()
    product, err := productService.CreateProduct("Test Product", 10.0)
    assert.NoError(t, err)
    assert.NotNil(t, product)
}
