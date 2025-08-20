package fatturapa

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Map of document types to their corresponding tags
var documentTypeTags = map[string][]cbc.Key{
	// Standard invoices
	"TD01": {},
	// Advance or down payment
	"TD02": {tax.TagPartial},
	// Advance or down payment on freelance invoice
	"TD03": {tax.TagPartial, sdi.TagFreelance},
	// Credit notes
	"TD04": {},
	// Debit notes
	"TD05": {},
	// Freelancer invoice
	"TD06": {sdi.TagFreelance},
	// Simplified invoice
	"TD07": {tax.TagSimplified},
	// Simplified credit note
	"TD08": {tax.TagSimplified},
	// Simplified debit note
	"TD09": {tax.TagSimplified},
	// Reverse charge
	"TD16": {tax.TagSelfBilled, tax.TagReverseCharge},
	// Self-billed import
	"TD17": {tax.TagSelfBilled, sdi.TagImport},
	// Self-billed EU goods import
	"TD18": {tax.TagSelfBilled, sdi.TagImport, sdi.TagGoodsEU},
	// Self-billed goods import
	"TD19": {tax.TagSelfBilled, sdi.TagImport, sdi.TagGoods},
	// Self-billed regularization
	"TD20": {tax.TagSelfBilled, sdi.TagRegularization},
	// Self-billed ceiling exceeded
	"TD21": {tax.TagSelfBilled, sdi.TagCeilingExceeded},
	// Self-billed goods extracted
	"TD22": {tax.TagSelfBilled, sdi.TagGoodsExtracted},
	// Self-billed goods with tax
	"TD23": {tax.TagSelfBilled, sdi.TagGoodsWithTax},
	// Deferred invoice
	"TD24": {sdi.TagDeferred},
	// Deferred invoice third period
	"TD25": {sdi.TagDeferred, sdi.TagThirdPeriod},
	// Depreciable assets
	"TD26": {sdi.TagDepreciableAssets},
	// Self-billed for self consumption
	"TD27": {tax.TagSelfBilled},
	// Self-billed San Marino paper
	"TD28": {tax.TagSelfBilled, sdi.TagSanMarinoPaper},
}

// goblBillInvoiceAddBody adds a body to the GOBL invoice
func goblBillInvoiceAddBody(inv *bill.Invoice, body *Body) error {
	if inv == nil || body == nil {
		return nil
	}

	// Add general data
	if err := goblBillInvoiceAddGeneralData(inv, body.GeneralData); err != nil {
		return fmt.Errorf("adding general data: %w", err)
	}

	// Extract retained taxes from the general data
	var retainedTaxes []*RetainedTax
	if body.GeneralData != nil && body.GeneralData.Document != nil {
		retainedTaxes = body.GeneralData.Document.RetainedTaxes
	}

	// Add goods and services, passing the retained taxes
	if err := goblBillInvoiceAddGoodsServices(inv, body.GoodsServices, retainedTaxes); err != nil {
		return fmt.Errorf("adding goods and services: %w", err)
	}

	// Add payment data
	goblBillInvoiceAddPaymentsData(inv, body.PaymentsData)

	return nil
}

// goblBillInvoiceAddGeneralData adds general data to the GOBL invoice
func goblBillInvoiceAddGeneralData(inv *bill.Invoice, generalData *GeneralData) error {
	if inv == nil || generalData == nil {
		return nil
	}

	// Add general document data
	if err := goblBillInvoiceAddGeneralDocumentData(inv, generalData.Document); err != nil {
		return fmt.Errorf("adding general document data: %w", err)
	}

	// Add document references
	inv.Preceding = goblOrgDocumentRefsFromDocumentRefs(generalData.Preceding)

	// Create a new ordering object and populate it with document references
	ordering := &bill.Ordering{
		Purchases: goblOrgDocumentRefsFromDocumentRefs(generalData.Purchases),
		Contracts: goblOrgDocumentRefsFromDocumentRefs(generalData.Contracts),
		Tender:    goblOrgDocumentRefsFromDocumentRefs(generalData.Tender),
		Receiving: goblOrgDocumentRefsFromDocumentRefs(generalData.Receiving),
	}

	// Only set the ordering if at least one of the document reference arrays is not empty
	if len(ordering.Purchases) > 0 || len(ordering.Contracts) > 0 || len(ordering.Tender) > 0 || len(ordering.Receiving) > 0 {
		inv.Ordering = ordering
	}

	return nil
}

// goblBillInvoiceAddGeneralDocumentData adds general document data to the GOBL invoice
func goblBillInvoiceAddGeneralDocumentData(inv *bill.Invoice, doc *GeneralDocumentData) error {
	if inv == nil || doc == nil {
		return nil
	}

	// Set invoice type based on document type
	goblBillInvoiceAddDocumentType(inv, doc.DocumentType)

	// Add currency
	inv.Currency = currency.Code(doc.Currency)

	// Add issue date
	date, err := parseDate(doc.IssueDate)
	if err != nil {
		return fmt.Errorf("adding issue date: %w", err)
	}
	inv.IssueDate = date

	// Add number
	// Check if the number contains a series (format: "SERIES-CODE")
	parseSeriesAndCode(doc.Number, &inv.Series, &inv.Code)

	// Add totals payable
	if doc.TotalAmount != "" {
		payable, err := num.AmountFromString(doc.TotalAmount)
		if err != nil {
			return fmt.Errorf("adding totals payable: %w", err)
		}
		if inv.Totals == nil {
			inv.Totals = new(bill.Totals)
		}
		inv.Totals.Payable = payable
	}

	// Add stamp duty
	goblBillInvoiceAddStampDuty(inv, doc.StampDuty)

	// Add price adjustments
	goblBillInvoiceAddPriceAdjustments(inv, doc.PriceAdjustments)

	// Add invoice reasons
	goblBillInvoiceAddReasons(inv, doc.Reasons)

	// Add retained taxes to totals
	// goblBillInvoiceAddRetainedTaxes(inv, doc.RetainedTaxes)

	return nil
}

// goblBillInvoiceAddDocumentType adds document type information to the GOBL invoice
func goblBillInvoiceAddDocumentType(inv *bill.Invoice, documentType string) {
	if inv == nil || documentType == "" {
		return
	}

	// Add tax extension key
	if inv.Tax == nil {
		inv.Tax = new(bill.Tax)
	}
	if inv.Tax.Ext == nil {
		inv.Tax.Ext = tax.Extensions{}
	}
	inv.Tax.Ext[sdi.ExtKeyDocumentType] = cbc.Code(documentType)

	// Set invoice type based on document type
	switch documentType {
	case "TD01", "TD02", "TD03", "TD06", "TD07", "TD16", "TD17", "TD18", "TD19", "TD20", "TD21", "TD22", "TD23", "TD24", "TD25", "TD26", "TD27", "TD28":
		inv.Type = bill.InvoiceTypeStandard
	case "TD04", "TD08":
		inv.Type = bill.InvoiceTypeCreditNote
	case "TD05", "TD09":
		inv.Type = bill.InvoiceTypeDebitNote
	default:
		// Default to standard if not recognized
		inv.Type = bill.InvoiceTypeStandard
	}

	// Get tags for the document type
	if tags, ok := documentTypeTags[documentType]; ok && len(tags) > 0 {
		// Get existing tags
		existingTags := inv.GetTags()

		// Add new tags
		for _, tag := range tags {
			if !tag.In(existingTags...) {
				existingTags = append(existingTags, tag)
			}
		}

		// Set the updated tags
		inv.SetTags(existingTags...)
	}
}

// goblBillInvoiceAddStampDuty adds stamp duty information from the FatturaPA document to the GOBL invoice
func goblBillInvoiceAddStampDuty(inv *bill.Invoice, stampDuty *StampDuty) {
	if inv == nil || stampDuty == nil || stampDuty.Amount == "" {
		return
	}

	amount, err := num.AmountFromString(stampDuty.Amount)
	if err == nil {
		if inv.Charges == nil {
			inv.Charges = make([]*bill.Charge, 0)
		}
		inv.Charges = append(inv.Charges, &bill.Charge{
			Key:    bill.ChargeKeyStampDuty,
			Amount: amount,
		})
	}
}

// goblBillInvoiceAddPriceAdjustments adds price adjustments (discounts and charges) from the FatturaPA document to the GOBL invoice
func goblBillInvoiceAddPriceAdjustments(inv *bill.Invoice, adjustments []*PriceAdjustment) {
	if inv == nil || len(adjustments) == 0 {
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

		if adj.Type == scontoMaggiorazioneTypeDiscount {
			if inv.Discounts == nil {
				inv.Discounts = make([]*bill.Discount, 0)
			}
			inv.Discounts = append(inv.Discounts, &bill.Discount{
				Amount:  amount,
				Percent: percentPtr,
			})
		} else if adj.Type == scontoMaggiorazioneTypeCharge {
			if inv.Charges == nil {
				inv.Charges = make([]*bill.Charge, 0)
			}
			inv.Charges = append(inv.Charges, &bill.Charge{
				Amount:  amount,
				Percent: percentPtr,
			})
		}
	}
}

// goblBillInvoiceAddReasons adds invoice reasons from the FatturaPA document to the GOBL invoice
func goblBillInvoiceAddReasons(inv *bill.Invoice, reasons []string) {
	if inv == nil || len(reasons) == 0 {
		return
	}

	for _, reason := range reasons {
		note := &org.Note{
			Key:  org.NoteKeyReason,
			Text: reason,
		}
		inv.Notes = append(inv.Notes, note)
	}
}

// goblOrgDocumentRefsFromDocumentRefs converts a slice of DocumentRef to a slice of org.DocumentRef
func goblOrgDocumentRefsFromDocumentRefs(refs []*DocumentRef) []*org.DocumentRef {
	if len(refs) == 0 {
		return nil
	}

	result := make([]*org.DocumentRef, 0, len(refs))
	for _, ref := range refs {
		if orgRef := goblOrgDocumentRefFromDocumentRef(ref); orgRef != nil {
			result = append(result, orgRef)
		}
	}
	return result
}

// goblOrgDocumentRefFromDocumentRef converts a DocumentRef to an org.DocumentRef
func goblOrgDocumentRefFromDocumentRef(ref *DocumentRef) *org.DocumentRef {
	if ref == nil {
		return nil
	}

	orgRef := &org.DocumentRef{}

	// Parse series and code
	parseSeriesAndCode(ref.Code, &orgRef.Series, &orgRef.Code)

	// Add issue date
	if ref.IssueDate != "" {
		date, err := parseDate(ref.IssueDate)
		if err != nil {
			return nil
		}
		orgRef.IssueDate = &date
	}

	// Add identities
	if ref.LineCode != "" {
		orgRef.Identities = append(orgRef.Identities, &org.Identity{
			Key:  org.IdentityKeyItem,
			Code: cbc.Code(ref.LineCode),
		})
	}

	if ref.OrderCode != "" {
		orgRef.Identities = append(orgRef.Identities, &org.Identity{
			Key:  org.IdentityKeyOrder,
			Code: cbc.Code(ref.OrderCode),
		})
	}

	if ref.CIGCode != "" {
		orgRef.Identities = append(orgRef.Identities, &org.Identity{
			Type: sdi.IdentityTypeCIG,
			Code: cbc.Code(ref.CIGCode),
		})
	}

	if ref.CUPCode != "" {
		orgRef.Identities = append(orgRef.Identities, &org.Identity{
			Type: sdi.IdentityTypeCUP,
			Code: cbc.Code(ref.CUPCode),
		})
	}

	return orgRef
}

// parseSeriesAndCode parses a document number into series and code components
// If the number contains a hyphen (format: "SERIES-CODE"), it will split the string
// and set the series and code accordingly. Otherwise, it will set only the code.
func parseSeriesAndCode(number string, series *cbc.Code, code *cbc.Code) {
	if code != nil {
		*code = cbc.Code(number)
	}

	parts := strings.Split(number, "-")
	if len(parts) > 1 && series != nil && code != nil {
		*series = cbc.Code(parts[0])
		*code = cbc.Code(parts[1])
	}
}

// adjustTotals compares the totals of the invoice with the totals of the FatturaPA document and adds rounding if necessary
func adjustTotals(inv *bill.Invoice, doc *GeneralDocumentData) error {
	if inv == nil || doc == nil {
		return nil
	}
	if doc.TotalAmount != "" {
		ft, err := num.AmountFromString(doc.TotalAmount)
		if err != nil {
			return err
		}

		// Calculate to get totals
		if err = inv.Calculate(); err != nil {
			return err
		}

		if inv.Totals == nil {
			return nil
		}

		r := ft.Subtract(inv.Totals.Payable)
		if r.Compare(num.AmountZero) != 0 {
			inv.Totals.Rounding = &r
			fmt.Println(r.String())
		}
	}

	return nil
}
