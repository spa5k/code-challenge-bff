package main

import (
	"log"
	"log/slog"

	"github.com/spa5k/zeller_go/internal"
	"github.com/spa5k/zeller_go/internal/catalog"
	"github.com/spa5k/zeller_go/internal/checkout"
	"github.com/spa5k/zeller_go/internal/pricingrules"
)

func main() {
	logger := internal.NewLogger()
	catalog := catalog.NewCatalog()
	pricingRules := map[string]pricingrules.PricingRule{
		"atv": &pricingrules.ThreeForTwoRule{SKU: "atv"},
		"ipd": &pricingrules.BulkDiscountRule{SKU: "ipd", MinQuantity: 4, NewPrice: 499.99},
	}

	// Scenario 1: SKUs Scanned: atv, atv, atv  -> 3 for 2 rule applies -> 2 * $109.50 = $219.00
	co1 := checkout.NewCheckout(pricingRules, catalog)
	err := co1.Scan(checkout.Item{SKU: "atv"})
	if err != nil {
		log.Fatal(err)
	}
	err = co1.Scan(checkout.Item{SKU: "atv"})
	if err != nil {
		log.Fatal(err)
	}
	err = co1.Scan(checkout.Item{SKU: "atv"})
	if err != nil {
		log.Fatal(err)
	}
	total1, err := co1.Total()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Scenario 1 Total price", "total", total1)

	// Buy 4 iPads so that the bulk discount rule applies, each ipad is $499.99 * 4 = $1999.96
	// Scenario 2: SKUs Scanned: ipd, ipd, ipd, ipd
	co2 := checkout.NewCheckout(pricingRules, catalog)
	for _, sku := range []string{"ipd", "ipd", "ipd", "ipd"} {
		err := co2.Scan(checkout.Item{SKU: sku})
		if err != nil {
			log.Fatal(err)
		}
	}
	total2, err := co2.Total()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Scenario 2 Total price", "total", total2)

	// Edge Case Scenario: Invalid SKU
	co3 := checkout.NewCheckout(pricingRules, catalog)
	err = co3.Scan(checkout.Item{SKU: "unknown"})
	if err != nil {
		slog.Error("Error scanning item", "error", err)
	}
	logger.Info("Scenario 3 Total price", "total", 0)

	// Example scenario - SKUs Scanned: atv, atv, atv, vga Total expected: $249.00
	co4 := checkout.NewCheckout(pricingRules, catalog)
	for _, sku := range []string{"atv", "atv", "atv", "vga"} {
		err := co4.Scan(checkout.Item{SKU: sku})
		if err != nil {
			log.Fatal(err)
		}
	}
	total4, err := co4.Total()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Scenario 4 Total price", "total", total4)

	// Example scenario - SKUs Scanned: atv, ipd, ipd, atv, ipd, ipd, ipd Total expected: $2718.95
	co5 := checkout.NewCheckout(pricingRules, catalog)
	for _, sku := range []string{"atv", "ipd", "ipd", "atv", "ipd", "ipd", "ipd"} {
		err := co5.Scan(checkout.Item{SKU: sku})
		if err != nil {
			log.Fatal(err)
		}
	}
	total5, err := co5.Total()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Scenario 5 Total price", "total", total5)
}
