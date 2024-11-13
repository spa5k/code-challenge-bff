package pricingrules_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"

	"github.com/spa5k/zeller_go/internal/catalog"
	"github.com/spa5k/zeller_go/internal/pricingrules"
)

func TestThreeForTwoRule_Apply_Success(t *testing.T) {
	rule := &pricingrules.ThreeForTwoRule{SKU: "atv"}
	c := catalog.NewCatalog()

	items := []pricingrules.Item{
		{SKU: "atv"},
		{SKU: "atv"},
		{SKU: "atv"},
	}

	total, err := rule.Apply(items, c)
	assert.NoError(t, err)

	expectedTotal := 2 * 109.50
	assert.Equal(t, expectedTotal, total)
}

func TestThreeForTwoRule_Apply_NoItems(t *testing.T) {
	rule := &pricingrules.ThreeForTwoRule{SKU: "atv"}
	c := catalog.NewCatalog()

	items := []pricingrules.Item{}

	total, err := rule.Apply(items, c)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, total)
}

func TestThreeForTwoRule_ProductNotFound(t *testing.T) {
	rule := &pricingrules.ThreeForTwoRule{SKU: "unknown"}
	c := catalog.NewCatalog()

	items := []pricingrules.Item{
		{SKU: "unknown"},
		{SKU: "unknown"},
		{SKU: "unknown"},
	}

	_, err := rule.Apply(items, c)
	assert.Error(t, err)
}

func TestBulkDiscountRule_Apply_DiscountApplied(t *testing.T) {
	rule := &pricingrules.BulkDiscountRule{
		SKU:         "ipd",
		MinQuantity: 5,
		NewPrice:    499.99,
	}
	c := catalog.NewCatalog()

	items := []pricingrules.Item{
		{SKU: "ipd"},
		{SKU: "ipd"},
		{SKU: "ipd"},
		{SKU: "ipd"},
		{SKU: "ipd"},
	}

	total, err := rule.Apply(items, c)
	assert.NoError(t, err)

	expectedTotal := 5 * 499.99
	assert.Equal(t, expectedTotal, total)
}

func TestBulkDiscountRule_Apply_NoDiscount(t *testing.T) {
	rule := &pricingrules.BulkDiscountRule{
		SKU:         "ipd",
		MinQuantity: 5,
		NewPrice:    499.99,
	}
	c := catalog.NewCatalog()

	items := []pricingrules.Item{
		{SKU: "ipd"},
		{SKU: "ipd"},
		{SKU: "ipd"},
	}

	total, err := rule.Apply(items, c)
	assert.NoError(t, err)

	expectedTotal := 3 * 549.99
	assert.Equal(t, expectedTotal, total)
}

func TestBulkDiscountRule_ProductNotFound(t *testing.T) {
	rule := &pricingrules.BulkDiscountRule{
		SKU:         "unknown",
		MinQuantity: 5,
		NewPrice:    499.99,
	}
	c := catalog.NewCatalog()

	items := []pricingrules.Item{
		{SKU: "unknown"},
		{SKU: "unknown"},
		{SKU: "unknown"},
		{SKU: "unknown"},
		{SKU: "unknown"},
	}

	_, err := rule.Apply(items, c)
	assert.Error(t, err)
}

func TestBulkDiscountRule_Apply_NoItems(t *testing.T) {
	rule := &pricingrules.BulkDiscountRule{
		SKU:         "ipd",
		MinQuantity: 5,
		NewPrice:    499.99,
	}
	c := catalog.NewCatalog()

	items := []pricingrules.Item{}

	total, err := rule.Apply(items, c)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, total)
}

func TestThreeForTwoRule_Apply_VaryingQuantities(t *testing.T) {
	rule := &pricingrules.ThreeForTwoRule{SKU: "atv"}
	c := catalog.NewCatalog()

	rand.Seed(uint64(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		quantity := rand.Intn(20) // Random quantity between 0 and 19
		items := make([]pricingrules.Item, quantity)
		for j := 0; j < quantity; j++ {
			items[j] = pricingrules.Item{SKU: "atv"}
		}

		total, err := rule.Apply(items, c)
		assert.NoError(t, err)

		expectedChargeable := quantity - (quantity / 3)
		expectedTotal := float64(expectedChargeable) * 109.50
		assert.Equal(t, expectedTotal, total, "Failed for quantity %d", quantity)
	}
}

func TestBulkDiscountRule_Apply_AtThreshold(t *testing.T) {
	rule := &pricingrules.BulkDiscountRule{
		SKU:         "ipd",
		MinQuantity: 5,
		NewPrice:    499.99,
	}
	c := catalog.NewCatalog()

	// Test quantities around the threshold
	for quantity := 4; quantity <= 6; quantity++ {
		items := make([]pricingrules.Item, quantity)
		for i := 0; i < quantity; i++ {
			items[i] = pricingrules.Item{SKU: "ipd"}
		}

		total, err := rule.Apply(items, c)
		assert.NoError(t, err)

		var expectedTotal float64
		if quantity >= rule.MinQuantity {
			expectedTotal = float64(quantity) * rule.NewPrice
		} else {
			expectedTotal = float64(quantity) * 549.99
		}
		assert.Equal(t, expectedTotal, total, "Failed for quantity %d", quantity)
	}
}

func TestBulkDiscountRule_Apply_LargeQuantities(t *testing.T) {
	rule := &pricingrules.BulkDiscountRule{
		SKU:         "ipd",
		MinQuantity: 5,
		NewPrice:    499.99,
	}
	c := catalog.NewCatalog()

	quantities := []int{10, 50, 100, 1000}
	for _, quantity := range quantities {
		items := make([]pricingrules.Item, quantity)
		for i := 0; i < quantity; i++ {
			items[i] = pricingrules.Item{SKU: "ipd"}
		}

		total, err := rule.Apply(items, c)
		assert.NoError(t, err)

		expectedTotal := float64(quantity) * rule.NewPrice
		assert.Equal(t, expectedTotal, total, "Failed for quantity %d", quantity)
	}
}
