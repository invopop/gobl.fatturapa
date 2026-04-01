package fatturapa

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

// processRetainedTaxes processes retained taxes and adds them to the appropriate line items
// and fund contribution charges
func processRetainedTaxes(inv *bill.Invoice, lineDetails []*LineDetail, retainedTaxes []*RetainedTax, fundContributions []*FundContribution) error {
	if len(retainedTaxes) == 0 || len(inv.Lines) == 0 {
		return nil
	}

	// Collect only fund contribution charges that were marked with Ritenuta=SI
	// in the original XML. We match by fund type code against the invoice charges.
	retainedCharges := retainedFundContributionCharges(inv.Charges, fundContributions)

	// Process each retained tax
	for _, rt := range retainedTaxes {
		// Parse the retained tax rate and amount
		rtRate, err1 := num.PercentageFromString(strings.TrimSpace(rt.Rate) + "%")
		rtAmount, err2 := parseAmount(rt.Amount)
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid retained tax rate or amount: %s %s", rt.Rate, rt.Amount)
		}

		// Convert tax type to category code
		catCode, err := convertRetainedTaxType(rt.Type)
		if err != nil {
			return err
		}

		// Build the tax combo to apply
		taxCombo := &tax.Combo{
			Category: catCode,
			Percent:  &rtRate,
		}
		if rt.Reason != "" {
			taxCombo.Ext = tax.Extensions{
				sdi.ExtKeyRetained: cbc.Code(rt.Reason),
			}
		}

		// Try to match against a single line first (common case)
		matched := false
		for i, detail := range lineDetails {
			if detail.Retained == flagSI {
				line := inv.Lines[i]
				expectedAmount := rtRate.Of(*line.Total)

				if expectedAmount.Equals(rtAmount) {
					line.Taxes = append(line.Taxes, taxCombo)
					matched = true
					break
				}
			}
		}

		if matched {
			continue
		}

		// Try matching against the sum of all retained-flagged lines + fund contribution charges
		totalBase := num.MakeAmount(0, 2)
		var retainedLines []*bill.Line
		for i, detail := range lineDetails {
			if detail.Retained == flagSI {
				totalBase = totalBase.Add(*inv.Lines[i].Total)
				retainedLines = append(retainedLines, inv.Lines[i])
			}
		}
		for _, charge := range retainedCharges {
			totalBase = totalBase.Add(charge.Amount)
		}

		expectedTotal := rtRate.Of(totalBase)
		if expectedTotal.Equals(rtAmount) {
			// Apply retained tax to all matching lines and charges
			for _, line := range retainedLines {
				tc := *taxCombo
				tc.Ext = copyExtensions(taxCombo.Ext)
				line.Taxes = append(line.Taxes, &tc)
			}
			for _, charge := range retainedCharges {
				tc := *taxCombo
				tc.Ext = copyExtensions(taxCombo.Ext)
				charge.Taxes = append(charge.Taxes, &tc)
			}
			matched = true
		}

		if !matched {
			return fmt.Errorf("could not match retained tax: %s %s%% %s", rt.Type, rt.Rate, rt.Amount)
		}
	}

	return nil
}

// retainedFundContributionCharges returns the subset of invoice charges that
// correspond to fund contributions marked with Ritenuta=SI in the XML.
func retainedFundContributionCharges(charges []*bill.Charge, fcs []*FundContribution) []*bill.Charge {
	// Build a set of fund type codes that are marked as retained in the XML
	retainedTypes := make(map[cbc.Code]int)
	for _, fc := range fcs {
		if fc.Retained == flagSI {
			retainedTypes[cbc.Code(fc.Type)]++
		}
	}

	var out []*bill.Charge
	for _, charge := range charges {
		if !charge.Key.Has(sdi.KeyFundContribution) {
			continue
		}
		ft := charge.Ext[sdi.ExtKeyFundType]
		if count, ok := retainedTypes[ft]; ok && count > 0 {
			out = append(out, charge)
			retainedTypes[ft]--
		}
	}
	return out
}

// copyExtensions creates a shallow copy of tax extensions
func copyExtensions(ext tax.Extensions) tax.Extensions {
	if ext == nil {
		return nil
	}
	cp := make(tax.Extensions, len(ext))
	for k, v := range ext {
		cp[k] = v
	}
	return cp
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
