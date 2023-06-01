package fatturapa

import (
	"fmt"
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

var taxCategoryVat = regime.Category(common.TaxCategoryVAT)

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

func newDatiBeniServizi(inv *bill.Invoice) (*datiBeniServizi, error) {
	dl, err := generateLineDetails(inv)
	if err != nil {
		return nil, err
	}

	dr, err := generateTaxSummary(inv)
	if err != nil {
		return nil, err
	}

	return &datiBeniServizi{
		DettaglioLinee: dl,
		DatiRiepilogo:  dr,
	}, nil
}

func generateLineDetails(inv *bill.Invoice) ([]*dettaglioLinee, error) {
	var dl []*dettaglioLinee

	for _, line := range inv.Lines {
		vatRate := ""

		for _, tax := range line.Taxes {
			if tax.Category == common.TaxCategoryVAT {
				vatRate = formatPercentage(tax.Percent)
				break
			}
		}

		codeNatura, err := findCodeNaturaLine(line)
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

func generateTaxSummary(inv *bill.Invoice) ([]*datiRiepilogo, error) {
	var dr []*datiRiepilogo
	var vatRateTotals []*tax.RateTotal

	for _, cat := range inv.Totals.Taxes.Categories {
		if cat.Code == common.TaxCategoryVAT {
			vatRateTotals = cat.Rates
			break
		}
	}

	for _, rateTotal := range vatRateTotals {
		codeNatura, err := findCodeNaturaSummary(rateTotal)
		if err != nil {
			return nil, err
		}

		dr = append(dr, &datiRiepilogo{
			AliquotaIVA:          formatPercentage(rateTotal.Percent),
			Natura:               codeNatura,
			ImponibileImporto:    formatAmount(&rateTotal.Base),
			Imposta:              formatAmount(&rateTotal.Amount),
			RiferimentoNormativo: findRiferimentoNormativo(rateTotal),
		})
	}

	return dr, nil
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

func findCodeNaturaLine(line *bill.Line) (string, error) {
	rateKey := findZeroVatTaxRate(line)
	if rateKey == "" {
		return "", nil
	}

	return findCodeNatura(rateKey)
}

func findCodeNaturaSummary(rateTotal *tax.RateTotal) (string, error) {
	if !isZeroRate(rateTotal.Percent) {
		return "", nil
	}

	return findCodeNatura(rateTotal.Key)
}

func findCodeNatura(rateKey cbc.Key) (string, error) {
	rate := findRate(taxCategoryVat.Rates, rateKey)
	if rate == nil {
		return "", fmt.Errorf("'Natura' code not found for VAT rate of zero with key '%s'", rateKey)
	}

	code := rate.Codes[it.KeyFatturaPANatura]
	if code == "" {
		return "", fmt.Errorf("'Natura' code not found for VAT rate of zero with key '%s'", rateKey)
	}

	return code.String(), nil
}

func findRiferimentoNormativo(rateTotal *tax.RateTotal) string {
	if !isZeroRate(rateTotal.Percent) {
		return ""
	}

	rate := findRate(taxCategoryVat.Rates, rateTotal.Key)
	if rate == nil {
		return ""
	}

	return rate.Name[i18n.IT]
}

func findRate(rates []*tax.Rate, taxRateKey cbc.Key) *tax.Rate {
	for _, rate := range rates {
		if rate.Key == taxRateKey {
			return rate
		}
	}
	return nil
}

func findZeroVatTaxRate(line *bill.Line) cbc.Key {
	var combo *tax.Combo

	for _, tax := range line.Taxes {
		if tax.Category == common.TaxCategoryVAT {
			combo = tax
			break
		}
	}

	if combo == nil || !isZeroRate(combo.Percent) {
		return ""
	}

	return combo.Rate
}

func isZeroRate(percent *num.Percentage) bool {
	return percent == nil || percent.IsZero()
}
