package pricingrules

import (
	"github.com/shopspring/decimal"
	"github.com/spa5k/zeller_go/internal/catalog"
)

type Item struct {
	SKU string
}

type PricingRule interface {
	Apply(items []Item, catalog *catalog.Catalog) (float64, error)
}

// ThreeForTwoRule applies a "3 for 2" deal on a specific SKU
type ThreeForTwoRule struct {
	SKU string
}

func (r *ThreeForTwoRule) Apply(items []Item, catalog *catalog.Catalog) (float64, error) {
	if len(items) == 0 {
		return 0, nil
	}
	product, err := catalog.GetProduct(r.SKU)
	if err != nil {
		return 0, err
	}
	count := len(items)
	freeItems := count / 3
	chargeableCount := count - freeItems
	totalPrice := product.Price.Mul(decimal.NewFromInt(int64(chargeableCount))).InexactFloat64()
	return totalPrice, nil
}

// BulkDiscountRule applies a bulk discount when a minimum quantity is purchased
type BulkDiscountRule struct {
	SKU         string
	MinQuantity int
	NewPrice    float64
}

func (r *BulkDiscountRule) Apply(items []Item, catalog *catalog.Catalog) (float64, error) {
	if len(items) == 0 {
		return 0, nil
	}
	product, err := catalog.GetProduct(r.SKU)
	if err != nil {
		return 0, err
	}
	count := len(items)
	if count >= r.MinQuantity {
		return float64(count) * r.NewPrice, nil
	}
	return product.Price.Mul(decimal.NewFromInt(int64(count))).InexactFloat64(), nil
}
