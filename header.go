package fatturapa

import (
	"github.com/invopop/gobl/bill"
)

const (
	FormatoTrasmissione        = "FPA12"
	InboxKeyCodiceDestinatario = "codice-destinario"
)

// FatturaElettronicaHeader contains all data related to the parties involved
// in the document.
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

func newFatturaElettronicaHeader(inv *bill.Invoice) (*FatturaElettronicaHeader, error) {
	supplier, err := newCedentePrestatore(inv)
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

func newDatiTrasmissione(inv *bill.Invoice) DatiTrasmissione {
	cd := ""

	for _, inbox := range inv.Customer.Inboxes {
		if inbox.Key == InboxKeyCodiceDestinatario {
			cd = inbox.Code
		}
	}

	return DatiTrasmissione{
		IdTrasmittente: TaxID{
			IdPaese:  inv.Supplier.TaxID.Country.String(),
			IdCodice: inv.Supplier.TaxID.Code.String(),
		},
		ProgressivoInvio:    inv.Code,
		FormatoTrasmissione: "TODO",
		CodiceDestinatario:  cd,
	}
}
