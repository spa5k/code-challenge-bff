package checkout_test

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spa5k/zeller_go/internal/catalog"
	"github.com/spa5k/zeller_go/internal/checkout"
	"github.com/spa5k/zeller_go/internal/pricingrules"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
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

func TestCheckout_RandomBaskets(t *testing.T) {
	c := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"atv": &pricingrules.ThreeForTwoRule{SKU: "atv"},
		"ipd": &pricingrules.BulkDiscountRule{
			SKU:         "ipd",
			MinQuantity: 5,
			NewPrice:    499.99,
		},
	}

	skus := []string{"atv", "ipd", "mbp", "vga"}
	rand.Seed(uint64(time.Now().UnixNano()))

	for i := 0; i < 5; i++ {
		co := checkout.NewCheckout(pricingRules, c)
		basketSize := rand.Intn(20) // Random basket size between 0 and 19
		expectedTotal := 0.0
		itemCounts := make(map[string]int)

		// Build basket
		for j := 0; j < basketSize; j++ {
			sku := skus[rand.Intn(len(skus))]
			err := co.Scan(checkout.Item{SKU: sku})
			assert.NoError(t, err)
			itemCounts[sku]++
		}

		// Calculate expected total
		// Apply pricing rules manually for validation
		for sku, count := range itemCounts {
			product, _ := c.GetProduct(context.Background(), sku)
			switch sku {
			case "atv":
				freeItems := count / 3
				chargeable := count - freeItems
				expectedTotal += product.Price.Mul(decimal.NewFromInt(int64(chargeable))).InexactFloat64()
			case "ipd":
				if count >= 5 {
					expectedTotal += float64(count) * 499.99
				} else {
					expectedTotal += product.Price.Mul(decimal.NewFromInt(int64(count))).InexactFloat64()
				}
			default:
				expectedTotal += product.Price.Mul(decimal.NewFromInt(int64(count))).InexactFloat64()
			}
		}

		total, err := co.Total()
		assert.NoError(t, err)
		assert.Equal(t, expectedTotal, total, "Failed for random basket %d", i+1)
	}
}

func TestCheckout_BoundaryQuantities(t *testing.T) {
	c := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"atv": &pricingrules.ThreeForTwoRule{SKU: "atv"},
		"ipd": &pricingrules.BulkDiscountRule{
			SKU:         "ipd",
			MinQuantity: 5,
			NewPrice:    499.99,
		},
	}

	testCases := []struct {
		description string
		items       []string
	}{
		{
			"At threshold for bulk discount",
			[]string{"ipd", "ipd", "ipd", "ipd", "ipd"},
		},
		{
			"Just below threshold for bulk discount",
			[]string{"ipd", "ipd", "ipd", "ipd"},
		},
		{
			"At multiple of 3 for '3 for 2' deal",
			[]string{"atv", "atv", "atv"},
		},
		{
			"One over multiple of 3 for '3 for 2' deal",
			[]string{"atv", "atv", "atv", "atv"},
		},
	}

	for _, tc := range testCases {
		co := checkout.NewCheckout(pricingRules, c)
		for _, sku := range tc.items {
			err := co.Scan(checkout.Item{SKU: sku})
			assert.NoError(t, err)
		}

		total, err := co.Total()
		assert.NoError(t, err)

		// Manually calculate expected total
		expectedTotal := 0.0
		itemCounts := make(map[string]int)
		for _, sku := range tc.items {
			itemCounts[sku]++
		}

		for sku, count := range itemCounts {
			product, _ := c.GetProduct(context.Background(), sku)
			switch sku {
			case "atv":
				freeItems := count / 3
				chargeable := count - freeItems
				expectedTotal += product.Price.Mul(decimal.NewFromInt(int64(chargeable))).InexactFloat64()
			case "ipd":
				if count >= 5 {
					expectedTotal += float64(count) * 499.99
				} else {
					expectedTotal += product.Price.Mul(decimal.NewFromInt(int64(count))).InexactFloat64()
				}
			default:
				expectedTotal += product.Price.Mul(decimal.NewFromInt(int64(count))).InexactFloat64()
			}
		}

		assert.Equal(t, expectedTotal, total, tc.description)
	}
}

func TestCheckout_LargeQuantities(t *testing.T) {
	c := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"ipd": &pricingrules.BulkDiscountRule{
			SKU:         "ipd",
			MinQuantity: 5,
			NewPrice:    499.99,
		},
	}
	co := checkout.NewCheckout(pricingRules, c)

	quantity := 1000
	for i := 0; i < quantity; i++ {
		err := co.Scan(checkout.Item{SKU: "ipd"})
		assert.NoError(t, err)
	}

	total, err := co.Total()
	assert.NoError(t, err)

	expectedTotal := float64(quantity) * 499.99
	assert.Equal(t, expectedTotal, total)
}
