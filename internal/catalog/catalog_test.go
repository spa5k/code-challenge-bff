package catalog_test

import (
	"testing"

	"github.com/spa5k/zeller_go/internal/catalog"

	"github.com/stretchr/testify/assert"
)

func TestGetProduct_Success(t *testing.T) {
	c := catalog.NewCatalog()
	product, err := c.GetProduct("ipd")
	assert.NoError(t, err)
	assert.Equal(t, "ipd", product.SKU)
	assert.Equal(t, "Super iPad", product.Name)
	assert.Equal(t, 549.99, product.Price)
}

func TestGetProduct_EmptySKU(t *testing.T) {
	c := catalog.NewCatalog()
	_, err := c.GetProduct("")
	assert.Error(t, err)
}

func TestGetProduct_ProductNotFound(t *testing.T) {
	c := catalog.NewCatalog()
	_, err := c.GetProduct("unknown")
	assert.Error(t, err)
}

func TestNewCatalog(t *testing.T) {
	c := catalog.NewCatalog()
	assert.NotNil(t, c)
	assert.NotEmpty(t, c.Products())
}
