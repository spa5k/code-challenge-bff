package pricingrules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
