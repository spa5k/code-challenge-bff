package catalog

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spa5k/zeller_go/internal"
)

type Product struct {
	SKU   string
	Name  string
	Price float64
}

type Catalog struct {
	products map[string][]Product
}

var logger *slog.Logger

func init() {
	logger = internal.NewLogger()
}

func NewCatalog() *Catalog {
	return &Catalog{
		products: map[string][]Product{
			"ipd": {
				{SKU: "ipd", Name: "Super iPad", Price: 549.99},
			},
			"mbp": {
				{SKU: "mbp", Name: "MacBook Pro", Price: 1399.99},
			},
			"atv": {
				{SKU: "atv", Name: "Apple TV", Price: 109.50},
			},
			"vga": {
				{SKU: "vga", Name: "VGA adapter", Price: 30.00},
			},
		},
	}
}

func (c *Catalog) GetProducts(sku string) ([]Product, error) {
	if sku == "" {
		logger.Error("SKU cannot be empty")
		return nil, fmt.Errorf("SKU cannot be empty")
	}
	products, ok := c.products[sku]
	if !ok || len(products) == 0 {
		logger.Error("Product with SKU not found", "sku", sku)
		return nil, fmt.Errorf("Product with SKU '%s' not found", sku)
	}
	return products, nil
}

func (c *Catalog) Products() map[string][]Product {
	return c.products
}

func (c *Catalog) AddProduct(product Product) error {
	logger := internal.GetLogger(context.Background())
	logger.Info("Adding product", "product", product)
	if product.SKU == "" {
		logger.Error("Product SKU cannot be empty")
		return fmt.Errorf("Product SKU cannot be empty")
	}
	c.products[product.SKU] = append(c.products[product.SKU], product)
	return nil
}

func (c *Catalog) GetProduct(sku string) (Product, error) {
	logger := internal.GetLogger(context.Background())
	logger.Info("Getting product", "sku", sku)
	products, err := c.GetProducts(sku)
	if err != nil {
		return Product{}, err
	}
	logger.Info("Product found", "product", products[0])
	return products[0], nil
}
