# Checkout System for Zeller Computer Store

This project implements a flexible checkout system for Zeller's new computer store, incorporating specific pricing rules and designed to accommodate future changes easily.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Project Structure](#project-structure)
- [How It Works](#how-it-works)
- [Usage](#usage)
  - [Requirements](#requirements)
  - [Installation](#installation)
  - [Running the Application](#running-the-application)
- [Testing](#testing)
- [Pricing Rules](#pricing-rules)
- [Extending the System](#extending-the-system)
- [License](#license)

## Overview

Zeller is launching a new computer store with a flexible checkout system that can accommodate changing pricing rules. The system allows scanning items in any order and calculates the total price based on the current pricing rules.

## Features

- **Product Catalog**: Manages products with their SKU, name, and price.
- **Flexible Pricing Rules**: Easily add or modify pricing rules without changing the core checkout logic.
- **Promotions**:
  - **3 for 2 Deal on Apple TVs**: Buy 3 Apple TVs and pay for only 2.
  - **Bulk Discount on Super iPads**: Price drops to $499.99 each when buying 5 or more Super iPads.
- **Edge Case Handling**: Robust error handling for invalid SKUs, empty inputs, and other edge cases.
- **Unit Tests**: Comprehensive tests using the `testify` framework for easy assertions.

## Project Structure

The project follows Go's standard project layout:

```
- cmd/
  - main.go
- internal/
  - catalog/
    - catalog.go
    - catalog_test.go
  - checkout/
    - checkout.go
    - checkout_test.go
  - pricingrules/
    - pricingrules.go
    - pricingrules_test.go
- go.mod
- README.md
```

- **cmd/**: Contains the entry point of the application.
- **internal/**: Contains the internal packages:
  - **catalog/**: Manages the product catalog.
  - **checkout/**: Handles scanning items and calculating totals.
  - **pricingrules/**: Implements flexible pricing rules.

## How It Works

The checkout system consists of several components working together:

1. **Catalog**: Stores information about products, including SKU, name, and price.
2. **Pricing Rules**: Defines promotions and discounts applied during checkout.
3. **Checkout**: Scans items, applies pricing rules, and calculates the total price.

### Logic Flow

- **Scanning Items**: Items are scanned using their SKU.
- **Applying Pricing Rules**: For each unique SKU, the system checks if there are any applicable pricing rules and applies them.
- **Calculating Total**: The total price is calculated by summing up the prices of all items after applying the pricing rules.

## Usage

### Requirements

- Go 1.16 or higher

### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/checkout-system.git
   cd checkout-system
   ```

2. **Initialize Go Modules**

   ```bash
   go mod tidy
   ```

### Running the Application

1. **Run the Application**

   ```bash
   make run
   ```

2. **Expected Output**

   ```
   Total price: $249.00
   Total price: $2718.95
   Error scanning item: Product with SKU 'unknown' not found
   ```

   - The application runs two scenarios demonstrating the pricing rules and handles an edge case with an invalid SKU.

## Testing

The project includes comprehensive unit tests using the `testify` framework.

### Running Tests

```bash
make test
```

### Test Output

```
ok  	yourmodule/internal/catalog	    0.XXXs
ok  	yourmodule/internal/pricingrules	0.XXXs
ok  	yourmodule/internal/checkout	    0.XXXs
```

### Test Coverage

To check test coverage:

```bash
make test-cover
```

## Pricing Rules

### Current Promotions

1. **3 for 2 Deal on Apple TVs (`atv`):**

   - **Description**: Buy 3 Apple TVs and pay for only 2.
   - **Implementation**: `ThreeForTwoRule` in `pricingrules/`.

2. **Bulk Discount on Super iPads (`ipd`):**

   - **Description**: Price drops to $499.99 each when buying 5 or more.
   - **Implementation**: `BulkDiscountRule` in `pricingrules/`.

### Adding New Pricing Rules

To add new pricing rules:

1. **Create a New Rule Struct**: Implement the `PricingRule` interface.

   ```go
   type NewPricingRule struct {
       // Rule-specific fields
   }

   func (r *NewPricingRule) Apply(items []Item, catalog *catalog.Catalog) (float64, error) {
       // Rule logic
   }
   ```

2. **Register the Rule**: Add the new rule to the `pricingRules` map when initializing the checkout.

   ```go
   pricingRules := map[string]pricingrules.PricingRule{
       "newsku": &pricingrules.NewPricingRule{/* initialization */},
   }
   ```

## Extending the System

### Adding New Products

1. **Update the Catalog**

   - Add the new product to the `products` map in `catalog.go`.

   ```go
   "newsku": {SKU: "newsku", Name: "New Product", Price: 99.99},
   ```

### Handling Edge Cases

- **Invalid SKUs**: The system returns an error if an invalid SKU is scanned.
- **Empty Inputs**: The system handles empty SKUs and returns appropriate errors.
- **Zero Items**: Calculating the total with zero items returns zero without error.