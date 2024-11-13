package catalog

import "fmt"

type Product struct {
	SKU   string
	Name  string
	Price float64
}

type Catalog struct {
	products map[string]Product
}

func NewCatalog() *Catalog {
	return &Catalog{
		products: map[string]Product{
			"ipd": {SKU: "ipd", Name: "Super iPad", Price: 549.99},
			"mbp": {SKU: "mbp", Name: "MacBook Pro", Price: 1399.99},
			"atv": {SKU: "atv", Name: "Apple TV", Price: 109.50},
			"vga": {SKU: "vga", Name: "VGA adapter", Price: 30.00},
		},
	}
}

func (c *Catalog) GetProduct(sku string) (Product, error) {
	if sku == "" {
		return Product{}, fmt.Errorf("SKU cannot be empty")
	}
	p, ok := c.products[sku]
	if !ok {
		return Product{}, fmt.Errorf("Product with SKU '%s' not found", sku)
	}
	return p, nil
}

func (c *Catalog) Products() map[string]Product {
	return c.products
}
