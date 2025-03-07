package fatturapa

import (
	"fmt"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

func goblBillTotalsAddRetainedTaxes(totals *bill.Totals, retainedTaxes []*RetainedTax) {
	if totals == nil || len(retainedTaxes) == 0 {
		return
	}

	// Process each retained tax
	for _, rt := range retainedTaxes {
		// Convert the retained tax to a category total
		catTotal := goblTaxCategoryTotalFromRetainedTax(rt)
		if catTotal != nil {
			// Initialize totals only when needed
			if totals.Taxes == nil {
				totals.Taxes = new(tax.Total)
			}
			if totals.Taxes.Categories == nil {
				totals.Taxes.Categories = make([]*tax.CategoryTotal, 0)
			}

			// Check if we already have this category
			existingCat := findCategoryTotal(totals.Taxes.Categories, catTotal.Code)
			if existingCat != nil {
				// Add the rate to the existing category
				existingCat.Rates = append(existingCat.Rates, catTotal.Rates...)
			} else {
				// Add the new category
				totals.Taxes.Categories = append(totals.Taxes.Categories, catTotal)
			}
		}
	}
}

// goblTaxCategoryTotalFromRetainedTax converts a RetainedTax to a tax.CategoryTotal
func goblTaxCategoryTotalFromRetainedTax(rt *RetainedTax) *tax.CategoryTotal {
	if rt == nil {
		return nil
	}

	// Convert Type to tax category code
	catCode, err := convertTaxTypeToTaxCategory(rt.Type)
	if err != nil {
		return nil
	}

	// Parse amount and rate
	amount, err1 := num.AmountFromString(rt.Amount)
	// FatturaPA stores the rate as a percentage without the % symbol so we add it so that the conversion works
	rate, err2 := num.PercentageFromString(rt.Rate + "%")
	if err1 != nil || err2 != nil {
		return nil
	}

	// Create rate total with extensions
	rateTotal := &tax.RateTotal{
		Percent: &rate,
		Amount:  amount,
		Ext:     make(tax.Extensions),
	}

	// Add the reason code to extensions if present
	if rt.Reason != "" {
		rateTotal.Ext[sdi.ExtKeyRetained] = cbc.Code(rt.Reason)
	}

	// Create and return the category total
	return &tax.CategoryTotal{
		Code:     catCode,
		Retained: true,
		Rates:    []*tax.RateTotal{rateTotal},
	}
}

// convertTipoRitenutaToTaxCategory converts a TipoRitenuta code to a tax category code
func convertTaxTypeToTaxCategory(tipoRitenuta string) (cbc.Code, error) {
	switch tipoRitenuta {
	case "RT01":
		return it.TaxCategoryIRPEF, nil
	case "RT02":
		return it.TaxCategoryIRES, nil
	case "RT03":
		return it.TaxCategoryINPS, nil
	case "RT04":
		return it.TaxCategoryENASARCO, nil
	case "RT05":
		return it.TaxCategoryENPAM, nil
	default:
		return "", fmt.Errorf("unknown TipoRitenuta code: %s", tipoRitenuta)
	}
}

// findCategoryTotal finds a category total by code in a slice of category totals
func findCategoryTotal(categories []*tax.CategoryTotal, code cbc.Code) *tax.CategoryTotal {
	for _, cat := range categories {
		if cat.Code == code {
			return cat
		}
	}
	return nil
}
