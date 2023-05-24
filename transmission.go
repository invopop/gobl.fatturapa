package fatturapa

import (
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/it"
)

const (
	formatoTrasmissioneFPA12 = "FPA12"
	formatoTrasmissioneFPR12 = "FPR12"
)

// Invoices sent to Italian individuals or businesses can use 0000000 as the
// codice destinatario when it is not indicated explicitly.
// When the recipient is foreign, XXXXXXX is used.
const (
	defaultCodiceDestinatarioItalianBusiness = "0000000"
	defaultCodiceDestinatarioForeignBusiness = "XXXXXXX"
)

const inboxKeyCodiceDestinatario = "codice-destinatario"

// Data related to the transmission of the invoice
type datiTrasmissione struct {
	IdTrasmittente      taxID  `xml:",omitempty"` // nolint:revive
	ProgressivoInvio    string `xml:",omitempty"`
	FormatoTrasmissione string `xml:",omitempty"`
	CodiceDestinatario  string
}

func (c *Converter) newDatiTrasmissione(inv *bill.Invoice, env *gobl.Envelope) *datiTrasmissione {
	if c.Config.Transmitter == nil {
		return &datiTrasmissione{
			CodiceDestinatario: codiceDestinatario(inv.Customer),
		}
	}

	return &datiTrasmissione{
		IdTrasmittente: taxID{
			IdPaese:  c.Config.Transmitter.CountryCode,
			IdCodice: c.Config.Transmitter.TaxID,
		},
		ProgressivoInvio:    env.Head.UUID.String()[:8],
		FormatoTrasmissione: formatoTransmissione(inv.Customer),
		CodiceDestinatario:  codiceDestinatario(inv.Customer),
	}
}

func formatoTransmissione(cus *org.Party) string {
	taxID := cus.TaxID

	if taxID.Country == l10n.IT && taxID.Type == it.TaxIdentityTypeGovernment {
		return formatoTrasmissioneFPA12
	}

	return formatoTrasmissioneFPR12
}

func codiceDestinatario(cus *org.Party) string {
	if cus.TaxID.Country != l10n.IT {
		return defaultCodiceDestinatarioForeignBusiness
	}

	for _, inbox := range cus.Inboxes {
		if inbox.Key == inboxKeyCodiceDestinatario {
			return inbox.Code
		}
	}

	return defaultCodiceDestinatarioItalianBusiness
}
