package fatturapa

import (
	"github.com/invopop/gobl/bill"
)

// fatturaElettronicaHeader contains all data related to the parties involved
// in the document.
type fatturaElettronicaHeader struct {
	DatiTrasmissione *datiTrasmissione `xml:"DatiTrasmissione,omitempty"`
	Supplier         *Supplier         `xml:"CedentePrestatore,omitempty"`
	Customer         *Customer         `xml:"CessionarioCommittente,omitempty"`
}

func newFatturaElettronicaHeader(inv *bill.Invoice, datiTrasmissione *datiTrasmissione) *fatturaElettronicaHeader {
	return &fatturaElettronicaHeader{
		DatiTrasmissione: datiTrasmissione,
		Supplier:         newSupplier(inv.Supplier),
		Customer:         newCustomer(inv.Customer),
	}
}
