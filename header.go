package fatturapa

import (
	"github.com/invopop/gobl/bill"
)

const (
	FormatoTrasmissione = "FPA12"
)

type FatturaElettronicaHeader struct {
	DatiTrasmissione       DatiTrasmissione       `xml:",omitempty"`
	CedentePrestatore      CedentePrestatore      `xml:",omitempty"`
	CessionarioCommittente CessionarioCommittente `xml:",omitempty"`
}

// Data related to the transmitting subject
type DatiTrasmissione struct {
	IdTrasmittente      TaxID
	ProgressivoInvio    string
	FormatoTrasmissione string
	CodiceDestinatario  string
}

func newFatturaElettronicaHeader(inv bill.Invoice) (*FatturaElettronicaHeader, error) {
	supplier, err := newCedentePrestatore(inv.Supplier)
	if err != nil {
		return nil, err
	}

	customer, err := newCessionarioCommittente(inv.Customer)
	if err != nil {
		return nil, err
	}

	return &FatturaElettronicaHeader{
		DatiTrasmissione:       newDatiTrasmissione(inv),
		CedentePrestatore:      *supplier,
		CessionarioCommittente: *customer,
	}, nil
}

func newDatiTrasmissione(inv bill.Invoice) DatiTrasmissione {
	return DatiTrasmissione{
		IdTrasmittente: TaxID{
			IdPaese:  inv.Supplier.TaxID.Country.String(),
			IdCodice: inv.Supplier.TaxID.Code.String(),
		},
		ProgressivoInvio:    "TODO",
		FormatoTrasmissione: "TODO",
		CodiceDestinatario:  "TODO",
	}
}
