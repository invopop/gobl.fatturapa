package fatturapa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/it"
)

const (
	scontoMaggiorazioneTypeDiscount = "SC" // sconto
	scontoMaggiorazioneTypeCharge   = "MG" // maggiorazione
)

const (
	condizioniPagamentoInstallments = "TP01" // pagamenti in rate
	condizioniPagamentoFull         = "TP02" // pagamento completo
	condizioniPagamentoAdvance      = "TP03" // anticipo
)

// fatturaElettronicaBody contains all invoice data apart from the parties
// involved, which are contained in FatturaElettronicaHeader.
type fatturaElettronicaBody struct {
	DatiGenerali    *datiGenerali
	DatiBeniServizi *datiBeniServizi
	DatiPagamento   *datiPagamento `xml:",omitempty"`
}

// datiGenerali contains general data about the invoice such as retained taxes,
// invoice number, invoice date, document type, etc.
type datiGenerali struct {
	DatiGeneraliDocumento *datiGeneraliDocumento
}

type datiGeneraliDocumento struct {
	TipoDocumento       string
	Divisa              string
	Data                string
	Numero              string
	Causale             []string
	DatiRitenuta        []*datiRitenuta
	ScontoMaggiorazione []*scontoMaggiorazione
}

// scontoMaggiorazione contains data about price adjustments like discounts and
// charges.
type scontoMaggiorazione struct {
	Tipo        string
	Percentuale string
	Importo     string
}

func newFatturaElettronicaBody(inv *bill.Invoice) (*fatturaElettronicaBody, error) {
	dbs, err := newDatiBeniServizi(inv)
	if err != nil {
		return nil, err
	}

	dp, err := newDatiPagamento(inv)
	if err != nil {
		return nil, err
	}

	return &fatturaElettronicaBody{
		DatiGenerali: &datiGenerali{
			DatiGeneraliDocumento: &datiGeneraliDocumento{
				TipoDocumento:       findCodeTipoDocumento(inv),
				Divisa:              string(inv.Currency),
				Data:                inv.IssueDate.String(),
				Numero:              inv.Code,
				Causale:             extractInvoiceReasons(inv),
				DatiRitenuta:        extractRetainedTaxes(inv),
				ScontoMaggiorazione: extractPriceAdjustments(inv),
			},
		},
		DatiBeniServizi: dbs,
		DatiPagamento:   dp,
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

func extractPriceAdjustments(inv *bill.Invoice) []*scontoMaggiorazione {
	var scontiMaggiorazioni []*scontoMaggiorazione

	for _, discount := range inv.Discounts {
		scontiMaggiorazioni = append(scontiMaggiorazioni, &scontoMaggiorazione{
			Tipo:        scontoMaggiorazioneTypeDiscount,
			Percentuale: formatPercentage(discount.Percent),
			Importo:     formatAmount(&discount.Amount),
		})
	}

	for _, charge := range inv.Charges {
		scontiMaggiorazioni = append(scontiMaggiorazioni, &scontoMaggiorazione{
			Tipo:        scontoMaggiorazioneTypeCharge,
			Percentuale: formatPercentage(charge.Percent),
			Importo:     formatAmount(&charge.Amount),
		})
	}

	return scontiMaggiorazioni
}

func findCodeTipoDocumento(inv *bill.Invoice) string {
	ss := inv.ScenarioSummary()

	return ss.Meta[it.KeyFatturaPATipoDocumento]
}
