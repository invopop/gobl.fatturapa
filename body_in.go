package fatturapa

import (
	"strings"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// goblBillInvoiceAddBody adds a body to the GOBL invoice
func goblBillInvoiceAddBody(inv *bill.Invoice, body *Body) error {
	if inv == nil || body == nil {
		return nil
	}

	if err := goblBillInvoiceAddGeneralData(inv, body.GeneralData); err != nil {
		return err
	}

	if err := goblBillInvoiceAddGoodsServices(inv, body.GoodsServices); err != nil {
		return err
	}

	goblBillInvoiceAddPaymentsData(inv, body.PaymentsData)

	return nil
}

// goblBillInvoiceAddGeneralData adds general data to the GOBL invoice
func goblBillInvoiceAddGeneralData(inv *bill.Invoice, generalData *GeneralData) error {
	if inv == nil || generalData == nil {
		return nil
	}

	// Add document data
	if err := goblBillInvoiceAddGeneralDocumentData(inv, generalData.Document); err != nil {
		return err
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

	// Add tax extension key
	if inv.Tax == nil {
		inv.Tax = new(bill.Tax)
	}
	if inv.Tax.Ext == nil {
		inv.Tax.Ext = tax.Extensions{}
	}
	inv.Tax.Ext[sdi.ExtKeyDocumentType] = cbc.Code(doc.DocumentType)

	// Add currency
	inv.Currency = currency.Code(doc.Currency)

	// Add issue date
	date, err := parseDate(doc.IssueDate)
	if err != nil {
		return err
	}
	inv.IssueDate = date

	// Add number
	// Check if the number contains a series (format: "SERIES-CODE")
	parts := strings.Split(doc.Number, "-")
	if len(parts) > 1 {
		inv.Series = cbc.Code(parts[0])
		inv.Code = cbc.Code(parts[1])
	} else {
		inv.Code = cbc.Code(doc.Number)
	}

	// Add totals payable
	if inv.Totals == nil {
		inv.Totals = new(bill.Totals)
	}
	payable, err := num.AmountFromString(doc.TotalAmount)
	if err != nil {
		return err
	}
	inv.Totals.Payable = payable

	// Add stamp duty
	goblBillInvoiceAddStampDuty(inv, doc.StampDuty)

	// Add price adjustments
	goblBillInvoiceAddPriceAdjustments(inv, doc.PriceAdjustments)

	// Add invoice reasons
	goblBillInvoiceAddReasons(inv, doc.Reasons)

	// Add retained taxes
	goblBillTotalsAddRetainedTaxes(inv.Totals, doc.RetainedTaxes)

	return nil
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

	orgRef := &org.DocumentRef{
		Code: cbc.Code(ref.Code),
	}
	date, err := parseDate(ref.IssueDate)
	if err != nil {
		return nil
	}
	orgRef.IssueDate = &date

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
