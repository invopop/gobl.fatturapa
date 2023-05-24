package fatturapa

import (
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/common"
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
	PrezzoTotale        string
	AliquotaIVA         string
	Natura              string                 `xml:",omitempty"`
	ScontoMaggiorazione []*scontoMaggiorazione `xml:",omitempty"`
}

// datiRiepilogo contains summary data such as total amount, total VAT, etc.
type datiRiepilogo struct {
	AliquotaIVA       string
	ImponibileImporto string
	Imposta           string
	EsigibilitaIVA    string `xml:",omitempty"`
}

func newDatiBeniServizi(inv *bill.Invoice) *datiBeniServizi {
	return &datiBeniServizi{
		DettaglioLinee: newDettaglioLinee(inv),
		DatiRiepilogo:  newDatiRiepilogo(inv),
	}
}

func newDettaglioLinee(inv *bill.Invoice) []*dettaglioLinee {
	var dl []*dettaglioLinee

	for _, line := range inv.Lines {
		vatRate := ""

		for _, tax := range line.Taxes {
			if tax.Category == common.TaxCategoryVAT {
				vatRate = formatPercentage(tax.Percent)
				break
			}
		}

		dl = append(dl, &dettaglioLinee{
			NumeroLinea:         strconv.Itoa(line.Index),
			Descrizione:         line.Item.Name,
			Quantita:            formatAmount(&line.Quantity),
			PrezzoUnitario:      formatAmount(&line.Item.Price),
			PrezzoTotale:        formatAmount(&line.Sum),
			AliquotaIVA:         vatRate,
			Natura:              findCodeNaturaZeroVat(line),
			ScontoMaggiorazione: extractLinePriceAdjustments(line),
		})
	}

	return dl
}

func newDatiRiepilogo(inv *bill.Invoice) []*datiRiepilogo {
	var dr []*datiRiepilogo
	var vatRates []*tax.RateTotal

	for _, cat := range inv.Totals.Taxes.Categories {
		if cat.Code == common.TaxCategoryVAT {
			vatRates = cat.Rates
		}
	}

	for _, rate := range vatRates {
		dr = append(dr, &datiRiepilogo{
			AliquotaIVA:       formatPercentage(&rate.Percent),
			ImponibileImporto: formatAmount(&rate.Base),
			Imposta:           formatAmount(&rate.Amount),
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

func findCodeNaturaZeroVat(line *bill.Line) string {
	var tagKeys []cbc.Key

	for _, tax := range line.Taxes {
		if tax.Category == common.TaxCategoryVAT {
			tagKeys = tax.Tags
		}
	}

	if len(tagKeys) == 0 {
		return ""
	}

	taxCategoryVat := regime.Category(common.TaxCategoryVAT)

	if taxCategoryVat == nil {
		return ""
	}

	tagKey := tagKeys[0]

	for _, tag := range taxCategoryVat.Tags {
		if tag.Key == tagKey {
			return tag.Meta[it.KeyFatturaPANatura]
		}
	}

	return ""
}
