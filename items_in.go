package fatturapa

import (
	"strconv"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// goblBillInvoiceAddGoodsServices adds goods and services from the FatturaPA document to the GOBL invoice
func goblBillInvoiceAddGoodsServices(inv *bill.Invoice, goodsServices *GoodsServices) error {
	if inv == nil || goodsServices == nil {
		return nil
	}

	// Add line details
	if err := goblBillInvoiceAddLineDetails(inv, goodsServices.LineDetails); err != nil {
		return err
	}

	// Add tax summary
	goblBillInvoiceAddTaxSummary(inv, goodsServices.TaxSummary)

	return nil
}

// goblBillInvoiceAddLineDetails adds line details from the FatturaPA document to the GOBL invoice
func goblBillInvoiceAddLineDetails(inv *bill.Invoice, lineDetails []*LineDetail) error {
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
		// returns errors as they are required fields
		index, err := strconv.Atoi(detail.LineNumber)
		if err != nil {
			return err
		}
		quantity, err := num.AmountFromString(detail.Quantity)
		if err != nil {
			return err
		}
		unitPrice, err := num.AmountFromString(detail.UnitPrice)
		if err != nil {
			return err
		}
		totalPrice, err := num.AmountFromString(detail.TotalPrice)
		if err != nil {
			return err
		}

		// Create a new line
		line := &bill.Line{
			Index:    index,
			Quantity: quantity,
			Total:    totalPrice,
			Item: &org.Item{
				Name:  detail.Description,
				Price: unitPrice,
			},
		}

		// Add price adjustments
		goblBillLineAddPriceAdjustments(line, detail.PriceAdjustments)

		// Add tax information
		if detail.TaxRate != "" || detail.TaxNature != "" {
			// Create tax
			taxCombo := &tax.Combo{
				Category: tax.CategoryVAT,
			}

			// Add tax rate if it's not zero
			// FatturaPA stores the tax rate as a percentage without the % symbol so we add it so that the conversion works
			taxRate, _ := num.PercentageFromString(detail.TaxRate + "%")
			if taxRate != num.PercentageZero {
				taxCombo.Percent = &taxRate
			}

			// Add exempt extension if nature is provided
			if detail.TaxNature != "" {
				if taxCombo.Ext == nil {
					taxCombo.Ext = tax.Extensions{}
				}
				taxCombo.Ext[sdi.ExtKeyExempt] = cbc.Code(detail.TaxNature)
			}

			// Add tax to line
			line.Taxes = append(line.Taxes, taxCombo)
		}

		// Add line to invoice
		inv.Lines = append(inv.Lines, line)
	}

	return nil
}

// goblBillLineAddPriceAdjustments adds price adjustments to a line
func goblBillLineAddPriceAdjustments(line *bill.Line, adjustments []*PriceAdjustment) {
	if line == nil || len(adjustments) == 0 {
		return
	}

	for _, adj := range adjustments {
		amount, err1 := num.AmountFromString(adj.Amount)
		percent, err2 := num.PercentageFromString(adj.Percent)

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

// goblBillInvoiceAddTaxSummary adds tax summary from the FatturaPA document to the GOBL invoice
func goblBillInvoiceAddTaxSummary(inv *bill.Invoice, taxSummaries []*TaxSummary) {
	if inv == nil || len(taxSummaries) == 0 {
		return
	}

	// Initialize totals if needed
	if inv.Totals == nil {
		inv.Totals = new(bill.Totals)
	}
	if inv.Totals.Taxes == nil {
		inv.Totals.Taxes = new(tax.Total)
	}

	// Find or create VAT category
	var vatCategory *tax.CategoryTotal
	for _, cat := range inv.Totals.Taxes.Categories {
		if cat.Code == tax.CategoryVAT {
			vatCategory = cat
			break
		}
	}

	if vatCategory == nil {
		vatCategory = &tax.CategoryTotal{
			Code:  tax.CategoryVAT,
			Rates: make([]*tax.RateTotal, 0),
		}
		inv.Totals.Taxes.Categories = append(inv.Totals.Taxes.Categories, vatCategory)
	}

	// Process each tax summary
	for _, summary := range taxSummaries {
		// Parse tax rate, taxable amount, and tax amount
		taxRate, err1 := num.PercentageFromString(summary.TaxRate)
		taxableAmount, err2 := num.AmountFromString(summary.TaxableAmount)
		taxAmount, err3 := num.AmountFromString(summary.TaxAmount)

		if err1 != nil || err2 != nil || err3 != nil {
			// Skip if any of the values are invalid
			continue
		}

		// Create rate total
		rateTotal := &tax.RateTotal{
			Percent: &taxRate,
			Base:    taxableAmount,
			Amount:  taxAmount,
			Ext:     make(tax.Extensions),
		}

		// Add exempt extension if nature is provided
		if summary.TaxNature != "" {
			rateTotal.Ext[sdi.ExtKeyExempt] = cbc.Code(summary.TaxNature)
		}

		// Add VAT liability if provided
		if summary.TaxLiability != "" {
			rateTotal.Ext[sdi.ExtKeyVATLiability] = cbc.Code(summary.TaxLiability)
		}

		// Add rate total to VAT category
		vatCategory.Rates = append(vatCategory.Rates, rateTotal)
	}
}
