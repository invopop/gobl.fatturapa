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

// processRetainedTaxes processes retained taxes and distributes them across applicable line items and fund contributions
func processRetainedTaxes(inv *bill.Invoice, lineDetails []*LineDetail, retainedTaxes []*RetainedTax, fundContributions []*FundContribution) error {
	if len(retainedTaxes) == 0 {
		return nil
	}

	// If there's only one retained tax, distribute it proportionally
	if len(retainedTaxes) == 1 {
		return distributeSingleRetainedTax(inv, lineDetails, retainedTaxes[0], fundContributions)
	}

	// For multiple retained taxes, try to match each one to specific items
	return matchMultipleRetainedTaxes(inv, lineDetails, retainedTaxes, fundContributions)
}

// distributeSingleRetainedTax distributes a single retained tax across all applicable items
func distributeSingleRetainedTax(inv *bill.Invoice, lineDetails []*LineDetail, retainedTax *RetainedTax, fundContributions []*FundContribution) error {
	// Parse retained tax
	rtRate, err := num.PercentageFromString(retainedTax.Rate + "%")
	if err != nil {
		return fmt.Errorf("invalid retained tax rate: %s", retainedTax.Rate)
	}
	rtAmount, err := num.AmountFromString(retainedTax.Amount)
	if err != nil {
		return fmt.Errorf("invalid retained tax amount: %s", retainedTax.Amount)
	}
	catCode, err := convertRetainedTaxType(retainedTax.Type)
	if err != nil {
		return err
	}

	// Add at least 2 decimal places to amount for precision
	rtAmount = rtAmount.RescaleUp(2)

	// Collect retainable items and calculate total base
	var retainableLines []*bill.Line
	var retainableCharges []*bill.Charge
	totalBase := num.MakeAmount(0, 2) // Initialize with 2 decimal places for monetary precision

	// Add lines with Ritenuta="SI"
	for i, detail := range lineDetails {
		if detail.Retained == "SI" && i < len(inv.Lines) {
			line := inv.Lines[i]
			retainableLines = append(retainableLines, line)
			// Ensure line total has same precision before adding
			lineTotal := line.Total.RescaleUp(2)
			totalBase = totalBase.Add(lineTotal)
		}
	}

	// Add fund contributions with Ritenuta="SI"
	chargeIndex := 0
	for _, fc := range fundContributions {
		if fc.Retained == "SI" {
			for j, charge := range inv.Charges {
				if charge.Key.Has(sdi.KeyFundContribution) && j >= chargeIndex {
					retainableCharges = append(retainableCharges, charge)
					// Ensure charge amount has same precision before adding
					chargeAmount := charge.Amount.RescaleUp(2)
					totalBase = totalBase.Add(chargeAmount)
					chargeIndex = j + 1
					break
				}
			}
		}
	}

	if len(retainableLines) == 0 && len(retainableCharges) == 0 {
		return fmt.Errorf("retained taxes found but no items marked with Ritenuta='SI'")
	}

	// Check if this should be treated as exact matching instead of proportional
	expectedTotal := rtRate.Of(totalBase)
	tolerance := num.MakeAmount(1, 2) // 0.01 tolerance for precision differences

	diff := rtAmount.Sub(expectedTotal).Abs()

	// If the declared amount matches the total calculation exactly, apply the rate to all items
	if diff.Compare(tolerance) <= 0 {
		// Apply same rate to all items - the total will match
		for _, line := range retainableLines {
			if err := addRetainedTaxToLine(line, catCode, rtRate, retainedTax.Reason); err != nil {
				return fmt.Errorf("adding retained tax to line %d: %w", line.Index, err)
			}
		}
		for _, charge := range retainableCharges {
			if err := addRetainedTaxToCharge(charge, catCode, rtRate, retainedTax.Reason); err != nil {
				return fmt.Errorf("adding retained tax to charge: %w", err)
			}
		}
		return nil
	}

	// Try to match with individual lines first
	for _, line := range retainableLines {
		expectedAmount := rtRate.Of(*line.Total)
		if rtAmount.Sub(expectedAmount).Abs().Compare(tolerance) <= 0 {
			// Exact match - apply to this line only
			return addRetainedTaxToLine(line, catCode, rtRate, retainedTax.Reason)
		}
	}

	// Try to match with fund contributions
	for _, charge := range retainableCharges {
		expectedAmount := rtRate.Of(charge.Amount)
		if rtAmount.Sub(expectedAmount).Abs().Compare(tolerance) <= 0 {
			// Exact match - apply to this charge only
			return addRetainedTaxToCharge(charge, catCode, rtRate, retainedTax.Reason)
		}
	}

	// If no exact match found, return error
	return fmt.Errorf("cannot match retained tax %s (rate %s%%, amount %s) to any specific item or total",
		retainedTax.Type, retainedTax.Rate, retainedTax.Amount)
}

// matchMultipleRetainedTaxes attempts to match multiple retained taxes to specific items using heuristics
func matchMultipleRetainedTaxes(inv *bill.Invoice, lineDetails []*LineDetail, retainedTaxes []*RetainedTax, fundContributions []*FundContribution) error {
	// Collect all retainable items
	var lines []*bill.Line
	var charges []*bill.Charge
	usedLineIndices := make(map[int]bool)
	usedChargeIndices := make(map[int]bool)

	// Collect lines with Ritenuta="SI"
	for i, detail := range lineDetails {
		if detail.Retained == "SI" && i < len(inv.Lines) {
			lines = append(lines, inv.Lines[i])
		}
	}

	// Collect fund contributions with Ritenuta="SI"
	chargeIndex := 0
	for _, fc := range fundContributions {
		if fc.Retained == "SI" {
			for j, charge := range inv.Charges {
				if charge.Key.Has(sdi.KeyFundContribution) && j >= chargeIndex {
					charges = append(charges, charge)
					chargeIndex = j + 1
					break
				}
			}
		}
	}

	// Try to match each retained tax to a specific item
	tolerance := num.MakeAmount(1, 2) // 0.01 tolerance for precision differences

	for _, rt := range retainedTaxes {
		rtRate, err := num.PercentageFromString(rt.Rate + "%")
		if err != nil {
			return fmt.Errorf("invalid retained tax rate: %s", rt.Rate)
		}
		rtAmount, err := num.AmountFromString(rt.Amount)
		if err != nil {
			return fmt.Errorf("invalid retained tax amount: %s", rt.Amount)
		}
		// Ensure proper precision for retained tax amount
		rtAmount = rtAmount.RescaleUp(2)

		catCode, err := convertRetainedTaxType(rt.Type)
		if err != nil {
			return err
		}

		// First, check if this retained tax could match multiple items (ambiguous case)
		var potentialLineMatches []int
		var potentialChargeMatches []int

		// Check all available lines for potential matches
		for i, line := range lines {
			if usedLineIndices[i] {
				continue
			}
			expectedAmount := rtRate.Of(*line.Total)
			diff := rtAmount.Sub(expectedAmount).Abs()
			if diff.Compare(tolerance) <= 0 {
				potentialLineMatches = append(potentialLineMatches, i)
			}
		}

		// Check all available charges for potential matches
		for i, charge := range charges {
			if usedChargeIndices[i] {
				continue
			}
			expectedAmount := rtRate.Of(charge.Amount)
			diff := rtAmount.Sub(expectedAmount).Abs()
			if diff.Compare(tolerance) <= 0 {
				potentialChargeMatches = append(potentialChargeMatches, i)
			}
		}

		// Check for ambiguous matching (multiple potential matches)
		totalMatches := len(potentialLineMatches) + len(potentialChargeMatches)
		if totalMatches > 1 {
			return fmt.Errorf("retained tax %s (rate %s%%, amount %s) matches %d items - cannot determine which item it applies to",
				rt.Type, rt.Rate, rt.Amount, totalMatches)
		}

		matched := false

		// Apply to the single matching line if found
		if len(potentialLineMatches) == 1 {
			i := potentialLineMatches[0]
			line := lines[i]
			if err := addRetainedTaxToLine(line, catCode, rtRate, rt.Reason); err != nil {
				return fmt.Errorf("adding retained tax to line %d: %w", line.Index, err)
			}
			usedLineIndices[i] = true
			matched = true
		}

		// Apply to the single matching charge if found
		if len(potentialChargeMatches) == 1 && !matched {
			i := potentialChargeMatches[0]
			charge := charges[i]
			if err := addRetainedTaxToCharge(charge, catCode, rtRate, rt.Reason); err != nil {
				return fmt.Errorf("adding retained tax to charge: %w", err)
			}
			usedChargeIndices[i] = true
			matched = true
		}

		if !matched {
			return fmt.Errorf("cannot match retained tax %s (rate %s%%, amount %s) to any specific item",
				rt.Type, rt.Rate, rt.Amount)
		}
	}

	return nil
}

// addRetainedTaxToLine adds a retained tax to a bill line
func addRetainedTaxToLine(line *bill.Line, catCode cbc.Code, rate num.Percentage, reason string) error {
	// Check for existing retained tax of the same category
	for _, existing := range line.Taxes {
		if existing.Category == catCode {
			return fmt.Errorf("line already has retained tax of category %s", catCode)
		}
	}

	// Create tax combo
	taxCombo := &tax.Combo{
		Category: catCode,
		Percent:  &rate,
		Ext:      tax.Extensions{},
	}

	if reason != "" {
		taxCombo.Ext[sdi.ExtKeyRetained] = cbc.Code(reason)
	}

	line.Taxes = append(line.Taxes, taxCombo)
	return nil
}

// addRetainedTaxToCharge adds a retained tax to a bill charge
func addRetainedTaxToCharge(charge *bill.Charge, catCode cbc.Code, rate num.Percentage, reason string) error {
	// Check for existing retained tax of the same category
	for _, existing := range charge.Taxes {
		if existing.Category == catCode {
			return fmt.Errorf("charge already has retained tax of category %s", catCode)
		}
	}

	// Create tax combo
	taxCombo := &tax.Combo{
		Category: catCode,
		Percent:  &rate,
		Ext:      tax.Extensions{},
	}

	if reason != "" {
		taxCombo.Ext[sdi.ExtKeyRetained] = cbc.Code(reason)
	}

	charge.Taxes = append(charge.Taxes, taxCombo)
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
	case "RT06":
		return it.TaxCategoryOTHER, nil
	default:
		return "", fmt.Errorf("unknown TipoRitenuta code: %s", tipoRitenuta)
	}
}
