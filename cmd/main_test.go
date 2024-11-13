package main_test

import (
	"testing"

	"github.com/spa5k/zeller_go/internal/catalog"
	"github.com/spa5k/zeller_go/internal/checkout"
	"github.com/spa5k/zeller_go/internal/pricingrules"
	"github.com/stretchr/testify/assert"
)

func TestScenario1(t *testing.T) {
	// Scenario 1: SKUs Scanned: atv, atv, atv  -> 3 for 2 rule applies -> 2 * $109.50 = $219.00
	catalog := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"atv": &pricingrules.ThreeForTwoRule{SKU: "atv"},
		"ipd": &pricingrules.BulkDiscountRule{SKU: "ipd", MinQuantity: 4, NewPrice: 499.99},
	}

	co := checkout.NewCheckout(pricingRules, catalog)
	err := co.Scan(checkout.Item{SKU: "atv"})
	assert.NoError(t, err)
	err = co.Scan(checkout.Item{SKU: "atv"})
	assert.NoError(t, err)
	err = co.Scan(checkout.Item{SKU: "atv"})
	assert.NoError(t, err)

	total, err := co.Total()
	assert.NoError(t, err)
	expectedTotal := 219.00
	assert.Equal(t, expectedTotal, total)
}

func TestScenario2(t *testing.T) {
	// Scenario 2: SKUs Scanned: ipd, ipd, ipd, ipd
	// Buy 4 iPads so that the bulk discount rule applies: 4 * $499.99 = $1999.96
	catalog := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"atv": &pricingrules.ThreeForTwoRule{SKU: "atv"},
		"ipd": &pricingrules.BulkDiscountRule{SKU: "ipd", MinQuantity: 4, NewPrice: 499.99},
	}

	co := checkout.NewCheckout(pricingRules, catalog)
	for _, sku := range []string{"ipd", "ipd", "ipd", "ipd"} {
		err := co.Scan(checkout.Item{SKU: sku})
		assert.NoError(t, err)
	}

	total, err := co.Total()
	assert.NoError(t, err)
	expectedTotal := 1999.96
	assert.Equal(t, expectedTotal, total)
}

func TestInvalidSKU(t *testing.T) {
	// Edge Case Scenario: Invalid SKU
	catalog := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"atv": &pricingrules.ThreeForTwoRule{SKU: "atv"},
		"ipd": &pricingrules.BulkDiscountRule{SKU: "ipd", MinQuantity: 4, NewPrice: 499.99},
	}

	co := checkout.NewCheckout(pricingRules, catalog)
	err := co.Scan(checkout.Item{SKU: "unknown"})
	assert.Error(t, err)
	assert.EqualError(t, err, "Product with SKU 'unknown' not found")
}

func TestScenario3(t *testing.T) {
	// Example scenario - SKUs Scanned: atv, atv, atv, vga Total expected: $249.00
	catalog := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"atv": &pricingrules.ThreeForTwoRule{SKU: "atv"},
		"ipd": &pricingrules.BulkDiscountRule{SKU: "ipd", MinQuantity: 4, NewPrice: 499.99},
	}

	co := checkout.NewCheckout(pricingRules, catalog)
	for _, sku := range []string{"atv", "atv", "atv", "vga"} {
		err := co.Scan(checkout.Item{SKU: sku})
		assert.NoError(t, err)
	}

	total, err := co.Total()
	assert.NoError(t, err)
	expectedTotal := 249.00
	assert.Equal(t, expectedTotal, total)
}

func TestScenario4(t *testing.T) {
	// Example scenario - SKUs Scanned: atv, ipd, ipd, atv, ipd, ipd, ipd Total expected: $2718.95
	catalog := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"atv": &pricingrules.ThreeForTwoRule{SKU: "atv"},
		"ipd": &pricingrules.BulkDiscountRule{SKU: "ipd", MinQuantity: 4, NewPrice: 499.99},
	}

	co := checkout.NewCheckout(pricingRules, catalog)
	for _, sku := range []string{"atv", "ipd", "ipd", "atv", "ipd", "ipd", "ipd"} {
		err := co.Scan(checkout.Item{SKU: sku})
		assert.NoError(t, err)
	}

	total, err := co.Total()
	assert.NoError(t, err)
	expectedTotal := 2718.95
	assert.Equal(t, expectedTotal, total)
}

func TestEdgeCase_EmptySKU(t *testing.T) {
	// Scanning an item with an empty SKU
	catalog := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{}

	co := checkout.NewCheckout(pricingRules, catalog)
	err := co.Scan(checkout.Item{SKU: ""})
	assert.Error(t, err)
	assert.EqualError(t, err, "Item SKU cannot be empty")
}

func TestEdgeCase_ZeroQuantityPricingRule(t *testing.T) {
	// Pricing rule with zero minimum quantity
	catalog := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"ipd": &pricingrules.BulkDiscountRule{SKU: "ipd", MinQuantity: 0, NewPrice: 400.00},
	}

	co := checkout.NewCheckout(pricingRules, catalog)
	err := co.Scan(checkout.Item{SKU: "ipd"})
	assert.NoError(t, err)
	total, err := co.Total()
	assert.NoError(t, err)

	// Since MinQuantity is 0, the discount should always apply
	expectedTotal := 400.00
	assert.Equal(t, expectedTotal, total)
}

func TestEdgeCase_HighQuantity(t *testing.T) {
	// Scanning a very large number of items
	catalog := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"ipd": &pricingrules.BulkDiscountRule{SKU: "ipd", MinQuantity: 5, NewPrice: 499.99},
	}

	co := checkout.NewCheckout(pricingRules, catalog)
	highQuantity := 10000
	for i := 0; i < highQuantity; i++ {
		err := co.Scan(checkout.Item{SKU: "ipd"})
		assert.NoError(t, err)
	}
	total, err := co.Total()
	assert.NoError(t, err)

	expectedTotal := float64(highQuantity) * 499.99
	assert.Equal(t, expectedTotal, total)
}

func TestEdgeCase_ScanAfterTotal(t *testing.T) {
	// Scanning items after calling Total()
	catalog := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{}

	co := checkout.NewCheckout(pricingRules, catalog)
	err := co.Scan(checkout.Item{SKU: "mbp"})
	assert.NoError(t, err)
	total1, err := co.Total()
	assert.NoError(t, err)

	err = co.Scan(checkout.Item{SKU: "vga"})
	assert.NoError(t, err)
	total2, err := co.Total()
	assert.NoError(t, err)

	// Ensure that the new item is included in the total
	expectedTotal1 := 1399.99
	expectedTotal2 := 1399.99 + 30.00

	assert.Equal(t, expectedTotal1, total1)
	assert.Equal(t, expectedTotal2, total2)
}

func TestEdgeCase_DiscountedPriceHigherThanOriginal(t *testing.T) {
	// Pricing rule that sets a new price higher than the original price
	catalog := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"ipd": &pricingrules.BulkDiscountRule{SKU: "ipd", MinQuantity: 5, NewPrice: 600.00},
	}

	co := checkout.NewCheckout(pricingRules, catalog)
	for i := 0; i < 5; i++ {
		co.Scan(checkout.Item{SKU: "ipd"})
	}
	total, err := co.Total()
	assert.NoError(t, err)

	// Since the discounted price is higher, the system should decide whether to apply it or not
	// For this test, we'll assume the system applies the higher price
	expectedTotal := 5 * 600.00
	assert.Equal(t, expectedTotal, total)
}

func TestEdgeCase_PricingRuleWithNegativePrice(t *testing.T) {
	// Pricing rule that sets a negative price
	catalog := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"vga": &pricingrules.BulkDiscountRule{SKU: "vga", MinQuantity: 1, NewPrice: -10.00},
	}

	co := checkout.NewCheckout(pricingRules, catalog)
	co.Scan(checkout.Item{SKU: "vga"})
	total, err := co.Total()
	assert.NoError(t, err)

	expectedTotal := -10.00
	assert.Equal(t, expectedTotal, total)
}
