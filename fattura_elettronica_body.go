package fatturapa

import (
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
)

const (
	TipoDocumentoDefault = "TD01"
)

type FatturaElettronicaBody struct {
	DatiGenerali    DatiGenerali
	DatiBeniServizi DatiBeniServizi
	DatiPagamento   DatiPagamento
}

type DatiGenerali struct {
	DatiGeneraliDocumento DatiGeneraliDocumento
}

type DatiGeneraliDocumento struct {
	TipoDocumento string
	Divisa        string
	Data          string
	Numero        string
	Causale       []string
}

type DatiOrdineAcquisto struct {
	RiferimentoNumeroLinea []string
	IdDocumento            string
	NumItem                string
}

type DatiBeniServizi struct {
	DettaglioLinee []DettaglioLinee
	DatiRiepilogo  []DatiRiepilogo
}

type DettaglioLinee struct {
	NumeroLinea    string
	Descrizione    string
	Quantita       string
	PrezzoUnitario string
	PrezzoTotale   string
	AliquotaIVA    string
}

type DatiRiepilogo struct {
	AliquotaIVA       string
	ImponibileImporto string
	Imposta           string
	EsigibilitaIVA    string
}

type DatiPagamento struct {
	CondizioniPagamento string
	DettaglioPagamento  []DettaglioPagamento
}

type DettaglioPagamento struct {
	ModalitaPagamento     string
	DataScadenzaPagamento string
	ImportoPagamento      string
}

func newFatturaElettronicaBody(inv bill.Invoice) (*FatturaElettronicaBody, error) {
	return &FatturaElettronicaBody{
		DatiGenerali: DatiGenerali{
			DatiGeneraliDocumento: DatiGeneraliDocumento{
				TipoDocumento: TipoDocumentoDefault,
				Divisa:        string(inv.Currency),
				Data:          inv.IssueDate.String(),
				Numero:        inv.Code,
				Causale:       extractInvoiceReasons(inv),
			},
		},
		DatiBeniServizi: DatiBeniServizi{
			DettaglioLinee: extractLines(inv),
			DatiRiepilogo:  extractTaxRates(inv),
		},
		DatiPagamento: DatiPagamento{
			CondizioniPagamento: "TP02", // TODO
			DettaglioPagamento: []DettaglioPagamento{
				{
					ModalitaPagamento: "MP05", // TODO
					ImportoPagamento:  inv.Totals.Due.String(),
				},
			},
		},
	}, nil
}

func extractInvoiceReasons(inv bill.Invoice) []string {
	// find inv.Notes with NoteKey as cbc.NoteKeyReason
	var reasons []string

	for _, note := range inv.Notes {
		if note.Key == cbc.NoteKeyReason {
			reasons = append(reasons, note.Text)
		}
	}

	return reasons
}

func extractLines(inv bill.Invoice) []DettaglioLinee {
	var lines []DettaglioLinee

	for _, line := range inv.Lines {
		desc := ""
		vatRate := ""

		for _, note := range line.Notes {
			if note.Key == cbc.NoteKeyGoods {
				desc += note.Text + "\n"
			}
		}

		for _, tax := range line.Taxes {
			if tax.Category == common.TaxCategoryVAT {
				vatRate = tax.Percent.String()
				break
			}
		}

		lines = append(lines, DettaglioLinee{
			NumeroLinea:    strconv.Itoa(line.Index),
			Descrizione:    desc,
			Quantita:       line.Quantity.String(),
			PrezzoUnitario: line.Item.Price.String(),
			PrezzoTotale:   line.Sum.String(),
			AliquotaIVA:    vatRate,
		})
	}

	return lines
}

func extractTaxRates(inv bill.Invoice) []DatiRiepilogo {
	rates := make(map[num.Percentage]num.Amount)

	for _, line := range inv.Lines {
		for _, tax := range line.Taxes {
			if tax.Category == common.TaxCategoryVAT {
				rates[tax.Percent] = rates[tax.Percent].Add(line.Sum)
			}
		}
	}

	var riepiloghi []DatiRiepilogo

	for rate, base := range rates {
		riepiloghi = append(riepiloghi, DatiRiepilogo{
			AliquotaIVA:       rate.String(),
			ImponibileImporto: base.String(),
			Imposta:           rate.Multiply(base).String(),
			EsigibilitaIVA:    "I", // TODO
		})
	}

	return riepiloghi
}
