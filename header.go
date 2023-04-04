package fatturapa

import (
	"github.com/invopop/gobl/bill"
)

// FatturaElettronicaHeader contains all data related to the parties involved
// in the document.
type FatturaElettronicaHeader struct {
	DatiTrasmissione       *DatiTrasmissione `xml:",omitempty"`
	CedentePrestatore      *Party            `xml:",omitempty"`
	CessionarioCommittente *Party            `xml:",omitempty"`
}

func newFatturaElettronicaHeader(inv *bill.Invoice, datiTrasmissione *DatiTrasmissione) (*FatturaElettronicaHeader, error) {
	supplier, err := newCedentePrestatore(inv)
	if err != nil {
		return nil, err
	}

	customer, err := newCessionarioCommittente(inv)
	if err != nil {
		return nil, err
	}

	return &FatturaElettronicaHeader{
		DatiTrasmissione:       datiTrasmissione,
		CedentePrestatore:      supplier,
		CessionarioCommittente: customer,
	}, nil
}
