package fatturapa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
)

const (
	ScontoMaggiorazioneTypeDiscount = "SC" // sconto
	ScontoMaggiorazioneTypeCharge   = "MG" // maggiorazione
)

const (
	CondizioniPagamentoInstallments = "TP01" // pagamenti in rate
	CondizioniPagamentoFull         = "TP02" // pagamento completo
	CondizioniPagamentoAdvance      = "TP03" // anticipo
)

// FatturaElettronicaBody contains all invoice data apart from the parties
// involved, which are contained in FatturaElettronicaHeader.
type FatturaElettronicaBody struct {
	DatiGenerali    *DatiGenerali
	DatiBeniServizi *DatiBeniServizi
	DatiPagamento   *DatiPagamento `xml:",omitempty"`
}

// DatiGenerali contains general data about the invoice such as retained taxes,
// invoice number, invoice date, document type, etc.
type DatiGenerali struct {
	DatiGeneraliDocumento *DatiGeneraliDocumento
}

type DatiGeneraliDocumento struct {
	TipoDocumento       string
	Divisa              string
	Data                string
	Numero              string
	Causale             []string
	DatiRitenuta        []*DatiRitenuta
	ScontoMaggiorazione []*ScontoMaggiorazione
}

type ScontoMaggiorazione struct {
	Tipo        string
	Percentuale string
	Importo     string
}

type DatiPagamento struct {
	CondizioniPagamento string
	DettaglioPagamento  []*DettaglioPagamento
}

type DettaglioPagamento struct {
	ModalitaPagamento     string
	DataScadenzaPagamento string `xml:",omitempty"`
	ImportoPagamento      string
}

func newFatturaElettronicaBody(inv *bill.Invoice) (*FatturaElettronicaBody, error) {
	return &FatturaElettronicaBody{
		DatiGenerali: &DatiGenerali{
			DatiGeneraliDocumento: &DatiGeneraliDocumento{
				TipoDocumento:       findCodeTipoDocumento(inv),
				Divisa:              string(inv.Currency),
				Data:                inv.IssueDate.String(),
				Numero:              inv.Code,
				Causale:             extractInvoiceReasons(inv),
				DatiRitenuta:        extractRetainedTaxes(inv),
				ScontoMaggiorazione: extractPriceAdjustments(inv),
			},
		},
		DatiBeniServizi: newDatiBeniServizi(inv),
		// GOBL does not yet support Italian codes for payment methods.
		//
		// DatiPagamento: DatiPagamento{
		// 	CondizioniPagamento: determinePaymentConditions(inv.Payment),
		// 	DettaglioPagamento: []DettaglioPagamento{
		// 		{
		// 			ModalitaPagamento: "MP05", // TODO
		// 			ImportoPagamento:  inv.Totals.Due.String(),
		// 		},
		// 	},
		// },
	}, nil
}

func extractInvoiceReasons(inv *bill.Invoice) []string {
	// find inv.Notes with NoteKey as cbc.NoteKeyReason
	var reasons []string

	for _, note := range inv.Notes {
		if note.Key == cbc.NoteKeyReason {
			reasons = append(reasons, note.Text)
		}
	}

	return reasons
}

func extractPriceAdjustments(inv *bill.Invoice) []*ScontoMaggiorazione {
	var scontiMaggiorazioni []*ScontoMaggiorazione

	for _, discount := range inv.Discounts {
		scontiMaggiorazioni = append(scontiMaggiorazioni, &ScontoMaggiorazione{
			Tipo:        ScontoMaggiorazioneTypeDiscount,
			Percentuale: discount.Percent.String(),
			Importo:     discount.Amount.String(),
		})
	}

	for _, charge := range inv.Charges {
		scontiMaggiorazioni = append(scontiMaggiorazioni, &ScontoMaggiorazione{
			Tipo:        ScontoMaggiorazioneTypeCharge,
			Percentuale: charge.Percent.String(),
			Importo:     charge.Amount.String(),
		})
	}

	return scontiMaggiorazioni
}

// func determinePaymentConditions(payment *bill.Payment) string {
// 	switch {
// 	case payment.Terms == nil:
// 		return CondizioniPagamentoFull
// 	case len(payment.Terms.DueDates) > 1:
// 		return CondizioniPagamentoInstallments
// 	case payment.Terms.Key == pay.TermKeyAdvance:
// 		return CondizioniPagamentoAdvance
// 	default:
// 		return CondizioniPagamentoFull
// 	}
// }
