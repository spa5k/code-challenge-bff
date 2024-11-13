package checkout

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/spa5k/zeller_go/internal/catalog"
	"github.com/spa5k/zeller_go/internal/pricingrules"
)

type Item struct {
	SKU string
}

type Checkout struct {
	pricingRules map[string]pricingrules.PricingRule
	items        []Item
	catalog      *catalog.Catalog
}

func NewCheckout(pricingRules map[string]pricingrules.PricingRule, catalog *catalog.Catalog) *Checkout {
	return &Checkout{
		pricingRules: pricingRules,
		catalog:      catalog,
	}
}

func (c *Checkout) Scan(item Item) error {
	if item.SKU == "" {
		return fmt.Errorf("Item SKU cannot be empty")
	}
	_, err := c.catalog.GetProduct(item.SKU)
	if err != nil {
		return err
	}
	c.items = append(c.items, item)
	return nil
}

func (c *Checkout) Total() (float64, error) {
	skuItems := make(map[string][]pricingrules.Item)
	for _, item := range c.items {
		skuItems[item.SKU] = append(skuItems[item.SKU], pricingrules.Item{SKU: item.SKU})
	}

	var total float64
	for sku, items := range skuItems {
		if rule, ok := c.pricingRules[sku]; ok {
			price, err := rule.Apply(items, c.catalog)
			if err != nil {
				return 0, err
			}
			total += price
		} else {
			product, err := c.catalog.GetProduct(sku)
			if err != nil {
				return 0, err
			}
			itemCount := decimal.NewFromInt(int64(len(items)))
			total += itemCount.Mul(product.Price).InexactFloat64()
		}
	}
	return total, nil
}
