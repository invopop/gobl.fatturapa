package fatturapa

import (
	"github.com/invopop/gobl/bill"
)

// fatturaElettronicaHeader contains all data related to the parties involved
// in the document.
type fatturaElettronicaHeader struct {
	DatiTrasmissione       *datiTrasmissione `xml:",omitempty"`
	CedentePrestatore      *supplier         `xml:",omitempty"`
	CessionarioCommittente *customer         `xml:",omitempty"`
}

func newFatturaElettronicaHeader(inv *bill.Invoice, datiTrasmissione *datiTrasmissione) (*fatturaElettronicaHeader, error) {
	supplier, err := newCedentePrestatore(inv)
	if err != nil {
		return nil, err
	}

	customer, err := newCessionarioCommittente(inv)
	if err != nil {
		return nil, err
	}

	return &fatturaElettronicaHeader{
		DatiTrasmissione:       datiTrasmissione,
		CedentePrestatore:      supplier,
		CessionarioCommittente: customer,
	}, nil
}
