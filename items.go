package fatturapa

import (
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

// datiBeniServizi contains all data related to the goods and services sold.
type datiBeniServizi struct {
	DettaglioLinee []*dettaglioLinee
	DatiRiepilogo  []*datiRiepilogo
}

// dettaglioLinee contains line data such as description, quantity, price, etc.
type dettaglioLinee struct {
	NumeroLinea         string
	Descrizione         string
	Quantita            string
	PrezzoUnitario      string
	ScontoMaggiorazione []*scontoMaggiorazione `xml:",omitempty"`
	PrezzoTotale        string
	AliquotaIVA         string
	Natura              string `xml:",omitempty"`
}

// datiRiepilogo contains tax summary data such as tax rate, tax amount, etc.
type datiRiepilogo struct {
	AliquotaIVA          string
	Natura               string `xml:",omitempty"`
	ImponibileImporto    string
	Imposta              string
	RiferimentoNormativo string `xml:",omitempty"`
}

func newDatiBeniServizi(inv *bill.Invoice) *datiBeniServizi {
	return &datiBeniServizi{
		DettaglioLinee: generateLineDetails(inv),
		DatiRiepilogo:  generateTaxSummary(inv),
	}
}

func generateLineDetails(inv *bill.Invoice) []*dettaglioLinee {
	var dl []*dettaglioLinee

	for _, line := range inv.Lines {
		d := &dettaglioLinee{
			NumeroLinea:         strconv.Itoa(line.Index),
			Descrizione:         line.Item.Name,
			Quantita:            formatAmount(&line.Quantity),
			PrezzoUnitario:      formatAmount(&line.Item.Price),
			PrezzoTotale:        formatAmount(&line.Sum),
			ScontoMaggiorazione: extractLinePriceAdjustments(line),
		}
		if len(line.Taxes) > 0 {
			vatTax := line.Taxes.Get(tax.CategoryVAT)
			if vatTax != nil {
				d.AliquotaIVA = formatPercentage(vatTax.Percent)
				d.Natura = vatTax.Ext[it.ExtKeySDINature].String()
			}
		}

		dl = append(dl, d)
	}

	return dl
}

func generateTaxSummary(inv *bill.Invoice) []*datiRiepilogo {
	var dr []*datiRiepilogo
	var vatRateTotals []*tax.RateTotal

	for _, cat := range inv.Totals.Taxes.Categories {
		if cat.Code == tax.CategoryVAT {
			vatRateTotals = cat.Rates
			break
		}
	}

	for _, rateTotal := range vatRateTotals {
		dr = append(dr, &datiRiepilogo{
			AliquotaIVA:          formatPercentage(rateTotal.Percent),
			Natura:               rateTotal.Ext[it.ExtKeySDINature].String(),
			ImponibileImporto:    formatAmount(&rateTotal.Base),
			Imposta:              formatAmount(&rateTotal.Amount),
			RiferimentoNormativo: findRiferimentoNormativo(rateTotal),
		})
	}

	return dr
}

func extractLinePriceAdjustments(line *bill.Line) []*scontoMaggiorazione {
	var scontiMaggiorazioni []*scontoMaggiorazione

	for _, discount := range line.Discounts {
		scontiMaggiorazioni = append(scontiMaggiorazioni, &scontoMaggiorazione{
			Tipo:        scontoMaggiorazioneTypeDiscount,
			Percentuale: formatPercentage(discount.Percent),
			Importo:     formatAmount(&discount.Amount),
		})
	}

	for _, charge := range line.Charges {
		scontiMaggiorazioni = append(scontiMaggiorazioni, &scontoMaggiorazione{
			Tipo:        scontoMaggiorazioneTypeCharge,
			Percentuale: formatPercentage(charge.Percent),
			Importo:     formatAmount(&charge.Amount),
		})
	}

	return scontiMaggiorazioni
}

func findRiferimentoNormativo(rateTotal *tax.RateTotal) string {
	def := regime.ExtensionDef(it.ExtKeySDINature)

	nature := rateTotal.Ext[it.ExtKeySDINature]
	for _, c := range def.Codes {
		if c.Code == nature {
			return c.Name[i18n.IT]
		}
	}

	return ""
}
