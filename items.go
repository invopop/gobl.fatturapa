package fatturapa

import (
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

type DatiBeniServizi struct {
	DettaglioLinee []*DettaglioLinee
	DatiRiepilogo  []*DatiRiepilogo
}

type DettaglioLinee struct {
	NumeroLinea         string
	Descrizione         string
	Quantita            string
	PrezzoUnitario      string
	PrezzoTotale        string
	AliquotaIVA         string
	Natura              string                 `xml:",omitempty"`
	ScontoMaggiorazione []*ScontoMaggiorazione `xml:",omitempty"`
}

type DatiRiepilogo struct {
	AliquotaIVA       string
	ImponibileImporto string
	Imposta           string
	EsigibilitaIVA    string `xml:",omitempty"`
}

func newDatiBeniServizi(inv *bill.Invoice) *DatiBeniServizi {
	return &DatiBeniServizi{
		DettaglioLinee: newDettaglioLinee(inv),
		DatiRiepilogo:  newDatiRiepilogo(inv),
	}
}

func newDettaglioLinee(inv *bill.Invoice) []*DettaglioLinee {
	var dl []*DettaglioLinee

	for _, line := range inv.Lines {
		vatRate := ""

		for _, tax := range line.Taxes {
			if tax.Category == common.TaxCategoryVAT {
				vatRate = formatPercentage(tax.Percent)
				break
			}
		}

		dl = append(dl, &DettaglioLinee{
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

func newDatiRiepilogo(inv *bill.Invoice) []*DatiRiepilogo {
	var dr []*DatiRiepilogo
	var vatRates []*tax.RateTotal

	for _, cat := range inv.Totals.Taxes.Categories {
		if cat.Code == common.TaxCategoryVAT {
			vatRates = cat.Rates
		}
	}

	for _, rate := range vatRates {
		dr = append(dr, &DatiRiepilogo{
			AliquotaIVA:       formatPercentage(&rate.Percent),
			ImponibileImporto: formatAmount(&rate.Base),
			Imposta:           formatAmount(&rate.Amount),
		})
	}

	return dr
}

func extractLinePriceAdjustments(line *bill.Line) []*ScontoMaggiorazione {
	var scontiMaggiorazioni []*ScontoMaggiorazione

	for _, discount := range line.Discounts {
		scontiMaggiorazioni = append(scontiMaggiorazioni, &ScontoMaggiorazione{
			Tipo:        ScontoMaggiorazioneTypeDiscount,
			Percentuale: formatPercentage(discount.Percent),
			Importo:     formatAmount(&discount.Amount),
		})
	}

	for _, charge := range line.Charges {
		scontiMaggiorazioni = append(scontiMaggiorazioni, &ScontoMaggiorazione{
			Tipo:        ScontoMaggiorazioneTypeCharge,
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
