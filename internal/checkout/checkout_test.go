package checkout_test

import (
	"testing"

	"github.com/spa5k/zeller_go/internal/catalog"
	"github.com/spa5k/zeller_go/internal/checkout"
	"github.com/spa5k/zeller_go/internal/pricingrules"
	"github.com/stretchr/testify/assert"
)

func TestCheckout_Scenario1(t *testing.T) {
	c := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"atv": &pricingrules.ThreeForTwoRule{SKU: "atv"},
		"ipd": &pricingrules.BulkDiscountRule{
			SKU:         "ipd",
			MinQuantity: 5,
			NewPrice:    499.99,
		},
	}
	co := checkout.NewCheckout(pricingRules, c)

	items := []string{"atv", "atv", "atv", "vga"}

	for _, sku := range items {
		err := co.Scan(checkout.Item{SKU: sku})
		assert.NoError(t, err)
	}

	total, err := co.Total()
	assert.NoError(t, err)

	expectedTotal := 249.00
	assert.Equal(t, expectedTotal, total)
}

func TestCheckout_Scenario2(t *testing.T) {
	c := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"atv": &pricingrules.ThreeForTwoRule{SKU: "atv"},
		"ipd": &pricingrules.BulkDiscountRule{
			SKU:         "ipd",
			MinQuantity: 5,
			NewPrice:    499.99,
		},
	}
	co := checkout.NewCheckout(pricingRules, c)

	items := []string{"atv", "ipd", "ipd", "atv", "ipd", "ipd", "ipd"}

	for _, sku := range items {
		err := co.Scan(checkout.Item{SKU: sku})
		assert.NoError(t, err)
	}

	total, err := co.Total()
	assert.NoError(t, err)

	expectedTotal := 2718.95
	assert.Equal(t, expectedTotal, total)
}

func TestCheckout_InvalidSKU(t *testing.T) {
	c := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{}
	co := checkout.NewCheckout(pricingRules, c)

	err := co.Scan(checkout.Item{SKU: "unknown"})
	assert.Error(t, err)
}

func TestCheckout_EmptySKU(t *testing.T) {
	c := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{}
	co := checkout.NewCheckout(pricingRules, c)

	err := co.Scan(checkout.Item{SKU: ""})
	assert.Error(t, err)
}

func TestCheckout_NoItemsScanned(t *testing.T) {
	c := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{}
	co := checkout.NewCheckout(pricingRules, c)

	total, err := co.Total()
	assert.NoError(t, err)

	expectedTotal := 0.0
	assert.Equal(t, expectedTotal, total)
}

func TestCheckout_MixedSKUs_NoPricingRules(t *testing.T) {
	c := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{}
	co := checkout.NewCheckout(pricingRules, c)

	items := []string{"ipd", "mbp", "vga"}

	for _, sku := range items {
		err := co.Scan(checkout.Item{SKU: sku})
		assert.NoError(t, err)
	}

	total, err := co.Total()
	assert.NoError(t, err)

	expectedTotal := 549.99 + 1399.99 + 30.00
	assert.Equal(t, expectedTotal, total)
}
