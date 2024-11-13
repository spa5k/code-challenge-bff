package internal

import "fmt"

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
