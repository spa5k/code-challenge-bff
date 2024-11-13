package catalog_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/spa5k/zeller_go/internal/catalog"
	"github.com/stretchr/testify/assert"
)

func TestGetProducts_Success(t *testing.T) {
	c := catalog.NewCatalog()
	products, err := c.GetProducts("ipd")
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	product := products[0]
	assert.Equal(t, "ipd", product.SKU)
	assert.Equal(t, "Super iPad", product.Name)
	assert.Equal(t, 549.99, product.Price)
}

func TestGetProducts_EmptySKU(t *testing.T) {
	c := catalog.NewCatalog()
	_, err := c.GetProducts("")
	assert.Error(t, err)
	assert.EqualError(t, err, "SKU cannot be empty")
}

func TestGetProducts_ProductNotFound(t *testing.T) {
	c := catalog.NewCatalog()
	_, err := c.GetProducts("unknown")
	assert.Error(t, err)
	assert.EqualError(t, err, "Product with SKU 'unknown' not found")
}

func TestNewCatalog(t *testing.T) {
	c := catalog.NewCatalog()
	assert.NotNil(t, c)
	assert.NotEmpty(t, c.Products())
}

func TestGetProducts_NegativePrice(t *testing.T) {
	c := catalog.NewCatalog()
	err := c.AddProduct(catalog.Product{SKU: "negprice", Name: "Negative Price Product", Price: -99.99})
	assert.NoError(t, err)
	products, err := c.GetProducts("negprice")
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, -99.99, products[0].Price)
}

func TestGetProducts_ZeroPrice(t *testing.T) {
	c := catalog.NewCatalog()
	err := c.AddProduct(catalog.Product{SKU: "zeroprice", Name: "Zero Price Product", Price: 0.0})
	assert.NoError(t, err)
	products, err := c.GetProducts("zeroprice")
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, 0.0, products[0].Price)
}

func TestGetProducts_LargeCatalog(t *testing.T) {
	c := catalog.NewCatalog()
	numProducts := 1000
	for i := 0; i < numProducts; i++ {
		sku := fmt.Sprintf("sku%04d", i) // Ensure the format string is in quotes
		err := c.AddProduct(catalog.Product{SKU: sku, Name: "Product " + sku, Price: float64(i)})
		assert.NoError(t, err)
	}
	initialProductCount := 4 // Number of initial products
	expectedProductCount := numProducts + initialProductCount
	assert.Len(t, c.Products(), expectedProductCount)
}

func TestGetProducts_MaxFloatPrice(t *testing.T) {
	c := catalog.NewCatalog()
	err := c.AddProduct(catalog.Product{SKU: "maxfloat", Name: "Max Float Price Product", Price: math.MaxFloat64})
	assert.NoError(t, err)
	products, err := c.GetProducts("maxfloat")
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, math.MaxFloat64, products[0].Price)
}

func TestAddProduct_EmptySKU(t *testing.T) {
	c := catalog.NewCatalog()
	err := c.AddProduct(catalog.Product{SKU: "", Name: "Empty SKU Product", Price: 100.0})
	assert.Error(t, err)
	assert.EqualError(t, err, "Product SKU cannot be empty")
}

func TestAddProduct_EmptyName(t *testing.T) {
	c := catalog.NewCatalog()
	err := c.AddProduct(catalog.Product{SKU: "emptyname", Name: "", Price: 100.0})
	// Assuming that the catalog allows empty product names
	assert.NoError(t, err)
	products, err := c.GetProducts("emptyname")
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "", products[0].Name)
}

func TestAddProduct_DuplicateSKU(t *testing.T) {
	c := catalog.NewCatalog()
	err := c.AddProduct(catalog.Product{SKU: "duplicate", Name: "Duplicate SKU Product", Price: 100.0})
	assert.NoError(t, err)
	err = c.AddProduct(catalog.Product{SKU: "duplicate", Name: "Duplicate SKU Product 2", Price: 200.0})
	assert.NoError(t, err)

	products, err := c.GetProducts("duplicate")
	assert.NoError(t, err)
	assert.Len(t, products, 2)

	// Verify both products are stored
	assert.Equal(t, "Duplicate SKU Product", products[0].Name)
	assert.Equal(t, 100.0, products[0].Price)

	assert.Equal(t, "Duplicate SKU Product 2", products[1].Name)
	assert.Equal(t, 200.0, products[1].Price)
}

func TestGetProducts_MultipleProductsPerSKU(t *testing.T) {
	c := catalog.NewCatalog()
	// Add multiple products with the same SKU
	err := c.AddProduct(catalog.Product{SKU: "ipd", Name: "Super iPad Pro", Price: 649.99})
	assert.NoError(t, err)

	products, err := c.GetProducts("ipd")
	assert.NoError(t, err)
	assert.Len(t, products, 2)

	// Verify both products are retrieved
	assert.Equal(t, "Super iPad", products[0].Name)
	assert.Equal(t, 549.99, products[0].Price)

	assert.Equal(t, "Super iPad Pro", products[1].Name)
	assert.Equal(t, 649.99, products[1].Price)
}
