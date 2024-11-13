package internal

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// ErrProductNotFound represents an error when a product cannot be found by SKU
type ErrProductNotFound struct {
	SKU string
}

func NewProductNotFoundError(sku string) ErrProductNotFound {
	return ErrProductNotFound{
		SKU: sku,
	}
}

func (e ErrProductNotFound) Error() string {
	return fmt.Sprintf("product not found: %s", e.SKU)
}

// ErrEmptySKU represents an error when an empty SKU is provided
type ErrEmptySKU struct {
	Operation string
}

func NewEmptySKUError(operation string) ErrEmptySKU {
	return ErrEmptySKU{
		Operation: operation,
	}
}

func (e ErrEmptySKU) Error() string {
	return fmt.Sprintf("empty SKU provided: %s", e.Operation)
}

// ErrNegativePrice represents an error when a product has a negative price
type ErrNegativePrice struct {
	SKU   string
	Price decimal.Decimal
}

func NewNegativePriceError(sku string, price decimal.Decimal) ErrNegativePrice {
	return ErrNegativePrice{
		SKU:   sku,
		Price: price,
	}
}

func (e ErrNegativePrice) Error() string {
	return fmt.Sprintf("negative price not allowed for product %s: %s", e.SKU, e.Price)
}

// ErrInvalidProduct represents a generic product validation error
type ErrInvalidProduct struct {
	SKU    string
	Reason string
}

func NewInvalidProductError(sku, reason string) ErrInvalidProduct {
	return ErrInvalidProduct{
		SKU:    sku,
		Reason: reason,
	}
}

func (e ErrInvalidProduct) Error() string {
	return fmt.Sprintf("invalid product %s: %s", e.SKU, e.Reason)
}
