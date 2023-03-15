package fatturapa

import (
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

type DatiBeniServizi struct {
	DettaglioLinee []DettaglioLinee
	DatiRiepilogo  []DatiRiepilogo
}

type DettaglioLinee struct {
	NumeroLinea         string
	Descrizione         string
	Quantita            string
	PrezzoUnitario      string
	PrezzoTotale        string
	AliquotaIVA         string
	Natura              string                `xml:",omitempty"`
	ScontoMaggiorazione []ScontoMaggiorazione `xml:",omitempty"`
}

type DatiRiepilogo struct {
	AliquotaIVA       string
	ImponibileImporto string
	Imposta           string
	EsigibilitaIVA    string `xml:",omitempty"`
}

func newDatiBeniServizi(inv *bill.Invoice) DatiBeniServizi {
	return DatiBeniServizi{
		DettaglioLinee: newDettaglioLinee(inv),
		DatiRiepilogo:  newDatiRiepilogo(inv),
	}
}

func newDettaglioLinee(inv *bill.Invoice) []DettaglioLinee {
	var dl []DettaglioLinee

	for _, line := range inv.Lines {
		vatRate := ""

		for _, tax := range line.Taxes {
			if tax.Category == common.TaxCategoryVAT {
				vatRate = tax.Percent.String()
				break
			}
		}

		dl = append(dl, DettaglioLinee{
			NumeroLinea:         strconv.Itoa(line.Index),
			Descrizione:         line.Item.Name,
			Quantita:            line.Quantity.String(),
			PrezzoUnitario:      line.Item.Price.String(),
			PrezzoTotale:        line.Sum.String(),
			AliquotaIVA:         vatRate,
			Natura:              findCodeNaturaZeroVat(line),
			ScontoMaggiorazione: extractLinePriceAdjustments(line),
		})
	}

	return dl
}

func newDatiRiepilogo(inv *bill.Invoice) []DatiRiepilogo {
	var dr []DatiRiepilogo
	var vatRates []*tax.RateTotal

	for _, cat := range inv.Totals.Taxes.Categories {
		if cat.Code == common.TaxCategoryVAT {
			vatRates = cat.Rates
		}
	}

	for _, rate := range vatRates {
		dr = append(dr, DatiRiepilogo{
			AliquotaIVA:       rate.Percent.String(),
			ImponibileImporto: rate.Base.String(),
			Imposta:           rate.Amount.String(),
		})
	}

	return dr
}

func extractLinePriceAdjustments(line *bill.Line) []ScontoMaggiorazione {
	var scontiMaggiorazioni []ScontoMaggiorazione

	for _, discount := range line.Discounts {
		scontiMaggiorazioni = append(scontiMaggiorazioni, ScontoMaggiorazione{
			Tipo:        ScontoMaggiorazioneTypeDiscount,
			Percentuale: discount.Percent.String(),
			Importo:     discount.Amount.String(),
		})
	}

	for _, charge := range line.Charges {
		scontiMaggiorazioni = append(scontiMaggiorazioni, ScontoMaggiorazione{
			Tipo:        ScontoMaggiorazioneTypeCharge,
			Percentuale: charge.Percent.String(),
			Importo:     charge.Amount.String(),
		})
	}

	return scontiMaggiorazioni
}
