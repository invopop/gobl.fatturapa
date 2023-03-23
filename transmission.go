package fatturapa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

const (
	FormatoTrasmissioneFPA12 = "FPA12"
	FormatoTrasmissioneFPR12 = "FPR12"
)

// Invoices sent to Italian individuals or businesses can use 0000000 as the
// codice destinatario when it is not indicated explicitly.
// When the recipient is foreign, XXXXXXX is used.
const (
	DefaultCodiceDestinatarioItalianBusiness = "0000000"
	DefaultCodiceDestinatarioForeignBusiness = "XXXXXXX"
)

const InboxKeyCodiceDestinatario = "codice-destinatario"

// Data related to the transmitting subject
type DatiTrasmissione struct {
	IdTrasmittente      TaxID
	ProgressivoInvio    string
	FormatoTrasmissione string
	CodiceDestinatario  string
}

func newDatiTrasmissione(inv *bill.Invoice, c Client, uuid string) DatiTrasmissione {
	return DatiTrasmissione{
		IdTrasmittente: TaxID{
			IdPaese:  c.CountryCode,
			IdCodice: c.TaxID,
		},
		ProgressivoInvio:    uuid[:8],
		FormatoTrasmissione: formatoTransmissione(inv.Customer),
		CodiceDestinatario:  codiceDestinatario(inv.Customer),
	}
}

func formatoTransmissione(cus *org.Party) string {
	// TODO: check the tax ID type and return FPA12 if the invoice is issued
	// to a government entity. Return FPR12 otherwise.

	return FormatoTrasmissioneFPR12
}

func codiceDestinatario(cus *org.Party) string {
	if cus.TaxID.Country != l10n.IT {
		return DefaultCodiceDestinatarioForeignBusiness
	}

	for _, inbox := range cus.Inboxes {
		if inbox.Key == InboxKeyCodiceDestinatario {
			return inbox.Code
		}
	}

	return DefaultCodiceDestinatarioItalianBusiness
}
