package fatturapa

import (
	"errors"
	"fmt"

	"github.com/invopop/gobl/addons/it/sdi"
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

const stampDutyCode = "SI"

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
	TipoDocumento          string
	Divisa                 string
	Data                   string
	Numero                 string
	DatiRitenuta           []*datiRitenuta
	DatiBollo              *datiBollo `xml:",omitempty"`
	ScontoMaggiorazione    []*scontoMaggiorazione
	ImportoTotaleDocumento string `xml:",omitempty"`
	Causale                []string
}

// datiBollo contains data about the stamp duty
type datiBollo struct {
	BolloVirtuale string
	ImportoBollo  string `xml:",omitempty"`
}

// scontoMaggiorazione contains data about price adjustments like discounts and
// charges.
type scontoMaggiorazione struct {
	Tipo        string `xml:"Tipo"`
	Percentuale string `xml:"Percentuale,omitempty"`
	Importo     string `xml:"Importo,omitempty"`
}

func newFatturaElettronicaBody(inv *bill.Invoice) (*fatturaElettronicaBody, error) {
	dbs := newDatiBeniServizi(inv)

	dp, err := newDatiPagamento(inv)
	if err != nil {
		return nil, err
	}

	dg, err := newDatiGenerali(inv)
	if err != nil {
		return nil, err
	}

	return &fatturaElettronicaBody{
		DatiGenerali:    dg,
		DatiBeniServizi: dbs,
		DatiPagamento:   dp,
	}, nil
}

func newDatiGenerali(inv *bill.Invoice) (*datiGenerali, error) {
	dr, err := extractRetainedTaxes(inv)
	if err != nil {
		return nil, err
	}

	codeTipoDocumento, err := findCodeTipoDocumento(inv)
	if err != nil {
		return nil, err
	}

	switch codeTipoDocumento {
	case "TD07", "TD08", "TD09":
		return nil, errors.New("simplified invoices are not currently supported")
	}

	code := inv.Code
	if inv.Series != "" {
		code = cbc.Code(fmt.Sprintf("%s-%s", inv.Series, inv.Code))
	}

	return &datiGenerali{
		DatiGeneraliDocumento: &datiGeneraliDocumento{
			TipoDocumento:          codeTipoDocumento,
			Divisa:                 string(inv.Currency),
			Data:                   inv.IssueDate.String(),
			Numero:                 code.String(),
			DatiRitenuta:           dr,
			DatiBollo:              newDatiBollo(inv.Charges),
			ImportoTotaleDocumento: formatAmount(&inv.Totals.Payable),
			ScontoMaggiorazione:    extractPriceAdjustments(inv),
			Causale:                extractInvoiceReasons(inv),
		},
	}, nil
}

func findCodeTipoDocumento(inv *bill.Invoice) (string, error) {
	if inv.Tax == nil {
		return "", fmt.Errorf("missing tax")
	}

	val, ok := inv.Tax.Ext[sdi.ExtKeyDocumentType]
	if !ok || val == "" {
		return "", fmt.Errorf("missing %s", sdi.ExtKeyDocumentType)
	}

	return val.String(), nil
}

func newDatiBollo(charges []*bill.Charge) *datiBollo {
	for _, charge := range charges {
		if charge.Key == it.ChargeKeyStampDuty {
			return &datiBollo{
				BolloVirtuale: stampDutyCode,
				ImportoBollo:  formatAmount(&charge.Amount),
			}
		}
	}

	return nil
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
