package fatturapa

import (
	"github.com/invopop/gobl/bill"
)

func goblBillInvoiceAddHeader(inv *bill.Invoice, header *Header) {
	if inv == nil || header == nil {
		return
	}

	inv.Supplier = goblOrgPartyFromSupplier(header.Supplier)
	inv.Customer = goblOrgPartyFromCustomer(header.Customer)

	// Need to do after customer is set
	goblBillInvoiceAddTransmission(inv, header.TransmissionData)
}
