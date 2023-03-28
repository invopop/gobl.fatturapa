package fatturapa

import (
	"github.com/invopop/gobl/bill"
)

// FatturaElettronicaHeader contains all data related to the parties involved
// in the document.
type FatturaElettronicaHeader struct {
	DatiTrasmissione       DatiTrasmissione       `xml:",omitempty"`
	CedentePrestatore      CedentePrestatore      `xml:",omitempty"`
	CessionarioCommittente CessionarioCommittente `xml:",omitempty"`
}

func newFatturaElettronicaHeader(inv *bill.Invoice, c *Client, uuid string) (*FatturaElettronicaHeader, error) {
	supplier, err := newCedentePrestatore(inv)
	if err != nil {
		return nil, err
	}

	customer, err := newCessionarioCommittente(inv)
	if err != nil {
		return nil, err
	}

	return &FatturaElettronicaHeader{
		DatiTrasmissione:       newDatiTrasmissione(inv, c, uuid),
		CedentePrestatore:      *supplier,
		CessionarioCommittente: *customer,
	}, nil
}
