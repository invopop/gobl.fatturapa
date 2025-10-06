package fatturapa

import (
	"fmt"
	"strconv"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// goblBillInvoiceAddGoodsServices adds goods and services from the FatturaPA document to the GOBL invoice
func goblBillInvoiceAddGoodsServices(inv *bill.Invoice, goodsServices *GoodsServices, retainedTaxes []*RetainedTax) error {
	if inv == nil || goodsServices == nil {
		return nil
	}

	// Add line details, passing the retained taxes and tax summaries
	if err := goblBillInvoiceAddLineDetails(inv, goodsServices.LineDetails, retainedTaxes, goodsServices.TaxSummary); err != nil {
		return fmt.Errorf("adding line details: %w", err)
	}

	return nil
}

// goblBillInvoiceAddLineDetails adds line details from the FatturaPA document to the GOBL invoice
func goblBillInvoiceAddLineDetails(inv *bill.Invoice, lineDetails []*LineDetail, retainedTaxes []*RetainedTax, taxSummaries []*TaxSummary) error {
	if inv == nil || len(lineDetails) == 0 {
		return nil
	}

	// Initialize lines if needed
	if inv.Lines == nil {
		inv.Lines = make([]*bill.Line, 0)
	}

	// Process each line detail
	for _, detail := range lineDetails {
		// Parse line number, quantity, unit price, and total price
		index, err := strconv.Atoi(detail.LineNumber)
		if err != nil {
			return fmt.Errorf("parsing line number: %w", err)
		}

		// Parse total price and unit price first
		unitPrice, err := num.AmountFromString(detail.UnitPrice)
		if err != nil {
			return fmt.Errorf("parsing unit price: %w", err)
		}
		totalPrice, err := num.AmountFromString(detail.TotalPrice)
		if err != nil {
			return fmt.Errorf("parsing total price: %w", err)
		}

		// Handle quantity parsing/calculation
		var quantity num.Amount
		if detail.Quantity != "" {
			quantity, err = num.AmountFromString(detail.Quantity)
			if err != nil {
				return fmt.Errorf("parsing quantity: %w", err)
			}
		} else {
			quantity = totalPrice.Divide(unitPrice)
		}

		// Create a new line
		line := &bill.Line{
			Index:    index,
			Quantity: quantity,
			Total:    &totalPrice,
			Item: &org.Item{
				Name:  detail.Description,
				Price: &unitPrice,
			},
		}

		// Add unit
		if detail.Unit != "" {
			line.Item.Unit = org.Unit(detail.Unit)
		}

		// Add price adjustments
		goblBillLineAddPriceAdjustments(line, detail.PriceAdjustments)

		// Add tax information
		if detail.TaxRate != "" || detail.TaxNature != "" {
			// Create tax
			taxCombo := &tax.Combo{
				Category: tax.CategoryVAT,
				Ext:      tax.Extensions{},
			}

			// Add tax rate if it's not zero
			taxRate, _ := num.PercentageFromString(detail.TaxRate + "%")
			taxCombo.Percent = &taxRate

			// Add exempt extension if nature is provided
			if detail.TaxNature != "" {
				taxCombo.Ext[sdi.ExtKeyExempt] = cbc.Code(detail.TaxNature)
				taxCombo.Percent = nil // Clear percent if exempt
			}

			// Add tax to line
			line.Taxes = append(line.Taxes, taxCombo)
		}

		// Add line to invoice
		inv.Lines = append(inv.Lines, line)
	}

	// Process retained taxes
	if err := processRetainedTaxes(inv, lineDetails, retainedTaxes); err != nil {
		return fmt.Errorf("processing retained taxes: %w", err)
	}

	// Match tax summary liability information with line items
	if len(inv.Lines) > 0 && len(taxSummaries) > 0 {
		goblBillLinesAddTaxSummary(inv.Lines, lineDetails, taxSummaries)
	}

	return nil
}

// goblBillLinesAddTaxSummary matches tax summary liability information with line items
func goblBillLinesAddTaxSummary(lines []*bill.Line, lineDetails []*LineDetail, taxSummaries []*TaxSummary) {
	// Process each line
	for i, line := range lines {
		// Skip if line has no taxes
		if len(line.Taxes) == 0 {
			continue
		}

		// Get the VAT tax for this line
		var vatTax *tax.Combo
		for _, t := range line.Taxes {
			if t.Category == tax.CategoryVAT {
				vatTax = t
				break
			}
		}

		// Skip if no VAT tax found
		if vatTax == nil {
			continue
		}

		// Get the line detail for this line
		detail := lineDetails[i]

		// Find a matching tax summary
		for _, summary := range taxSummaries {
			// Match by tax rate or nature
			rateMatches := detail.TaxRate == summary.TaxRate
			natureMatches := detail.TaxNature != "" && detail.TaxNature == summary.TaxNature

			if (rateMatches || natureMatches) && summary.TaxLiability != "" {
				// Add the tax liability to the VAT tax extensions
				if vatTax.Ext == nil {
					vatTax.Ext = tax.Extensions{}
				}
				vatTax.Ext[sdi.ExtKeyVATLiability] = cbc.Code(summary.TaxLiability)
				break
			}
		}
	}
}

// goblBillLineAddPriceAdjustments adds price adjustments to a line
func goblBillLineAddPriceAdjustments(line *bill.Line, adjustments []*PriceAdjustment) {
	if line == nil || len(adjustments) == 0 {
		return
	}

	for _, adj := range adjustments {
		amount, err1 := num.AmountFromString(adj.Amount)
		// FatturaPA stores the percentage as a string without the % symbol so we add it so that the conversion works
		percent, err2 := num.PercentageFromString(adj.Percent + "%")

		if err1 != nil && err2 != nil {
			// Skip if both amount and percent are invalid
			continue
		}

		// Check if percentage is zero, done so that if percentage is zero, it is not added to the invoice
		var percentPtr *num.Percentage
		if err2 == nil && percent != num.PercentageZero {
			percentPtr = &percent
		}

		// If the quantity is not 1, we need to multiply the amount by the quantity
		// because FatturaPA stores the discount/charge per unit
		if line.Quantity.Value() != 1 {
			amount = amount.Multiply(line.Quantity)
		}

		if adj.Type == scontoMaggiorazioneTypeDiscount {
			if line.Discounts == nil {
				line.Discounts = make([]*bill.LineDiscount, 0)
			}
			line.Discounts = append(line.Discounts, &bill.LineDiscount{
				Amount:  amount,
				Percent: percentPtr,
			})
		} else if adj.Type == scontoMaggiorazioneTypeCharge {
			if line.Charges == nil {
				line.Charges = make([]*bill.LineCharge, 0)
			}
			line.Charges = append(line.Charges, &bill.LineCharge{
				Amount:  amount,
				Percent: percentPtr,
			})
		}
	}
}
