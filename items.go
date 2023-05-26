package fatturapa

import (
	"fmt"
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
	ScontoMaggiorazione []*scontoMaggiorazione `xml:",omitempty"`
	PrezzoTotale        string
	AliquotaIVA         string
	Natura              string `xml:",omitempty"`
}

// datiRiepilogo contains summary data such as total amount, total VAT, etc.
type datiRiepilogo struct {
	AliquotaIVA       string
	ImponibileImporto string
	Imposta           string
	EsigibilitaIVA    string `xml:",omitempty"`
}

func newDatiBeniServizi(inv *bill.Invoice) (*datiBeniServizi, error) {
	dl, err := newDettaglioLinee(inv)
	if err != nil {
		return nil, err
	}

	return &datiBeniServizi{
		DettaglioLinee: dl,
		DatiRiepilogo:  newDatiRiepilogo(inv),
	}, nil
}

func newDettaglioLinee(inv *bill.Invoice) ([]*dettaglioLinee, error) {
	var dl []*dettaglioLinee

	for _, line := range inv.Lines {
		vatRate := ""

		for _, tax := range line.Taxes {
			if tax.Category == common.TaxCategoryVAT {
				vatRate = formatPercentage(tax.Percent)
				break
			}
		}

		codeNatura, err := findCodeNatura(line)
		if err != nil {
			return nil, err
		}

		dl = append(dl, &dettaglioLinee{
			NumeroLinea:         strconv.Itoa(line.Index),
			Descrizione:         line.Item.Name,
			Quantita:            formatAmount(&line.Quantity),
			PrezzoUnitario:      formatAmount(&line.Item.Price),
			PrezzoTotale:        formatAmount(&line.Sum),
			AliquotaIVA:         vatRate,
			Natura:              codeNatura,
			ScontoMaggiorazione: extractLinePriceAdjustments(line),
		})
	}

	return dl, nil
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
			AliquotaIVA:       formatPercentage(rate.Percent),
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

func findCodeNatura(line *bill.Line) (string, error) {
	taxRate := extractZeroVatTaxRate(line)

	if taxRate == "" {
		return "", nil
	}

	taxCategoryVat := regime.Category(common.TaxCategoryVAT)

	rate := findRate(taxCategoryVat.Rates, taxRate)
	if rate == nil {
		return "", fmt.Errorf("natura code not found for VAT rate of zero (line number: %d)", line.Index)
	}

	code := rate.Codes[it.KeyFatturaPANatura]
	if code == "" {
		return "", fmt.Errorf("natura code not found for VAT rate of zero (line number: %d)", line.Index)
	}

	return code.String(), nil
}

func findRate(rates []*tax.Rate, taxRateKey cbc.Key) *tax.Rate {
	for _, rate := range rates {
		if rate.Key == taxRateKey {
			return rate
		}
	}
	return nil
}

func extractZeroVatTaxRate(line *bill.Line) cbc.Key {
	var combo *tax.Combo

	for _, tax := range line.Taxes {
		if tax.Category == common.TaxCategoryVAT {
			combo = tax
		}
	}

	if combo == nil {
		return ""
	}

	if combo.Percent == nil || combo.Percent.IsZero() {
		return combo.Rate
	}

	return ""
}
