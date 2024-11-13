package catalog

import (
	"context"
	"log/slog"

	"github.com/shopspring/decimal"
	"github.com/spa5k/zeller_go/internal"
)

type Product struct {
	SKU   string
	Name  string
	Price decimal.Decimal
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
				{SKU: "ipd", Name: "Super iPad", Price: decimal.NewFromFloat(549.99)},
			},
			"mbp": {
				{SKU: "mbp", Name: "MacBook Pro", Price: decimal.NewFromFloat(1399.99)},
			},
			"atv": {
				{SKU: "atv", Name: "Apple TV", Price: decimal.NewFromFloat(109.50)},
			},
			"vga": {
				{SKU: "vga", Name: "VGA adapter", Price: decimal.NewFromFloat(30.00)},
			},
		},
	}
}

func (c *Catalog) GetProducts(ctx context.Context, sku string) ([]Product, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if sku == "" {
			logger.Error("SKU cannot be empty")
			return nil, internal.NewEmptySKUError("GetProducts")
		}
		products, ok := c.products[sku]
		if !ok || len(products) == 0 {
			logger.Error("Product with SKU not found", "sku", sku)
			return nil, internal.NewProductNotFoundError(sku)
		}
		return products, nil
	}
}

func (c *Catalog) Products() map[string][]Product {
	return c.products
}

func (c *Catalog) AddProduct(ctx context.Context, product Product) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		logger := internal.GetLogger(ctx)
		logger.Info("Adding product", "product", product)
		if product.SKU == "" {
			logger.Error("Product SKU cannot be empty")
			return internal.NewEmptySKUError("AddProduct")
		}
		c.products[product.SKU] = append(c.products[product.SKU], product)
		return nil
	}
}

func (c *Catalog) GetProduct(ctx context.Context, sku string) (Product, error) {
	select {
	case <-ctx.Done():
		return Product{}, ctx.Err()
	default:
		logger := internal.GetLogger(ctx)
		logger.Info("Getting product", "sku", sku)
		products, err := c.GetProducts(ctx, sku)
		if err != nil {
			return Product{}, err
		}
		logger.Info("Product found", "product", products[0])
		return products[0], nil
	}
}
