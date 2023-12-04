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

// Data related to the transmission of the invoice
type datiTrasmissione struct {
	IdTrasmittente      *taxID `xml:",omitempty"` // nolint:revive
	ProgressivoInvio    string `xml:",omitempty"`
	FormatoTrasmissione string `xml:",omitempty"`
	CodiceDestinatario  string
	PECDestinatario     string `xml:",omitempty"`
}

func (c *Converter) newDatiTrasmissione(inv *bill.Invoice, env *gobl.Envelope) *datiTrasmissione {
	dt := &datiTrasmissione{
		CodiceDestinatario: codiceDestinatario(inv.Customer),
		PECDestinatario:    pecDestinatario(inv.Customer),
	}

	// Do we need to add the transmitter info?
	if c.Config.Transmitter != nil {
		dt.IdTrasmittente = &taxID{
			IdPaese:  c.Config.Transmitter.CountryCode,
			IdCodice: c.Config.Transmitter.TaxID,
		}
		dt.ProgressivoInvio = env.Head.UUID.String()[:8]
		dt.FormatoTrasmissione = formatoTransmissione(inv.Customer)
	}

	return dt
}

func formatoTransmissione(cus *org.Party) string {
	if cus != nil {
		taxID := cus.TaxID
		if taxID.Country == l10n.IT && taxID.Type == it.TaxIdentityTypeGovernment {
			return formatoTrasmissioneFPA12
		}
	}

	return formatoTrasmissioneFPR12
}

func codiceDestinatario(cus *org.Party) string {
	if cus != nil {
		if cus.TaxID != nil && cus.TaxID.Country != l10n.IT {
			return defaultCodiceDestinatarioForeignBusiness
		}
		for _, inbox := range cus.Inboxes {
			if inbox.Key == it.KeyInboxSDICode {
				return inbox.Code
			}
		}
	}

	// When this is returned, we'll assume there is a PEC.
	// This is also valid for individuals.
	return defaultCodiceDestinatarioItalianBusiness
}

func pecDestinatario(cus *org.Party) string {
	if cus != nil {
		for _, inbox := range cus.Inboxes {
			if inbox.Key == it.KeyInboxSDIPEC {
				return inbox.Code
			}
		}
	}
	return ""
}
