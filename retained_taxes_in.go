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

// processRetainedTaxes processes retained taxes and adds them to the appropriate line items
func processRetainedTaxes(inv *bill.Invoice, lineDetails []*LineDetail, retainedTaxes []*RetainedTax) error {
	if len(retainedTaxes) == 0 || len(inv.Lines) == 0 {
		return nil
	}

	// Process each retained tax
	for _, rt := range retainedTaxes {
		// Parse the retained tax rate and amount
		rtRate, err1 := num.PercentageFromString(rt.Rate + "%")
		rtAmount, err2 := num.AmountFromString(rt.Amount)
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid retained tax rate or amount: %s %s", rt.Rate, rt.Amount)
		}

		// Convert tax type to category code
		catCode, err := convertRetainedTaxType(rt.Type)
		if err != nil {
			return err
		}

		// Find a line with Retained="SI" that matches this retained tax
		matched := false
		for i, detail := range lineDetails {
			// Only consider lines marked with Ritenuta="SI"
			if detail.Retained == "SI" {
				line := inv.Lines[i]
				expectedAmount := rtRate.Of(line.Total)

				// Check if this matches the retained tax amount exactly
				if expectedAmount.Equals(rtAmount) {
					// Create and add the retained tax to the line
					taxCombo := &tax.Combo{
						Category: catCode,
						Percent:  &rtRate,
						Ext:      tax.Extensions{},
					}

					// Add the reason code to extensions if present
					if rt.Reason != "" {
						taxCombo.Ext[sdi.ExtKeyRetained] = cbc.Code(rt.Reason)
					}

					line.Taxes = append(line.Taxes, taxCombo)
					matched = true
					break
				}
			}
		}

		// If we couldn't match this retained tax, return an error
		if !matched {
			return fmt.Errorf("could not match retained tax: %s %s%% %s", rt.Type, rt.Rate, rt.Amount)
		}
	}

	return nil
}

// convertRetainedTaxType converts a TipoRitenuta code to a tax category code
func convertRetainedTaxType(tipoRitenuta string) (cbc.Code, error) {
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
