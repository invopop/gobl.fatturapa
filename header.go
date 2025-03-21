package fatturapa

import "github.com/invopop/gobl/bill"

// Header contains all data related to the parties involved in the document.
type Header struct {
	TransmissionData *TransmissionData `xml:"DatiTrasmissione,omitempty"`
	Supplier         *Supplier         `xml:"CedentePrestatore,omitempty"`
	Customer         *Customer         `xml:"CessionarioCommittente,omitempty"`
}

func newHeader(inv *bill.Invoice, TransmissionData *TransmissionData) *Header {
	return &Header{
		TransmissionData: TransmissionData,
		Supplier:         newSupplier(inv.Supplier),
		Customer:         newCustomer(inv.Customer),
	}
}
