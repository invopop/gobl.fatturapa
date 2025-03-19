package fatturapa

import (
	"strconv"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// GoodsServices contains all data related to the goods and services sold.
type GoodsServices struct {
	LineDetails []*LineDetail `xml:"DettaglioLinee"`
	TaxSummary  []*TaxSummary `xml:"DatiRiepilogo"`
}

// LineDetail contains line data such as description, quantity, price, etc.
type LineDetail struct {
	LineNumber       string             `xml:"NumeroLinea"`
	Description      string             `xml:"Descrizione"`
	Quantity         string             `xml:"Quantita"`
	Unit             string             `xml:"UnitaMisura,omitempty"`
	UnitPrice        string             `xml:"PrezzoUnitario"`
	PriceAdjustments []*PriceAdjustment `xml:"ScontoMaggiorazione,omitempty"`
	TotalPrice       string             `xml:"PrezzoTotale"`
	TaxRate          string             `xml:"AliquotaIVA"`
	Retained         string             `xml:"Ritenuta,omitempty"`
	TaxNature        string             `xml:"Natura,omitempty"`
}

// TaxSummary contains tax summary data such as tax rate, tax amount, etc.
type TaxSummary struct {
	TaxRate        string `xml:"AliquotaIVA"`
	TaxNature      string `xml:"Natura,omitempty"`
	TaxableAmount  string `xml:"ImponibileImporto"`
	TaxAmount      string `xml:"Imposta"`
	TaxLiability   string `xml:"EsigibilitaIVA,omitempty"`
	LegalReference string `xml:"RiferimentoNormativo,omitempty"`
}

func newGoodsServices(inv *bill.Invoice) *GoodsServices {
	return &GoodsServices{
		LineDetails: generateLineDetails(inv),
		TaxSummary:  generateTaxSummary(inv),
	}
}

func generateLineDetails(inv *bill.Invoice) []*LineDetail {
	var dl []*LineDetail

	for _, line := range inv.Lines {
		d := &LineDetail{
			LineNumber:       strconv.Itoa(line.Index),
			Description:      line.Item.Name,
			Quantity:         formatAmount8(&line.Quantity),
			Unit:             string(line.Item.Unit),
			UnitPrice:        formatAmount8(line.Item.Price),
			TotalPrice:       formatAmount8(line.Total),
			PriceAdjustments: extractLinePriceAdjustments(line),
		}

		// Process taxes
		if len(line.Taxes) > 0 {
			// Process all taxes in a single loop
			for _, t := range line.Taxes {
				// Handle VAT tax
				if t.Category == tax.CategoryVAT {
					d.TaxRate = formatPercentageWithZero(t.Percent)
					d.TaxNature = exemptExtensionCode(t.Ext)
				} else if t.Ext != nil && t.Ext.Has(sdi.ExtKeyRetained) {
					// Check for retained taxes
					d.Retained = "SI"
				}
			}
		}

		dl = append(dl, d)
	}

	return dl
}

func exemptExtensionCode(ext tax.Extensions) string {
	if ext.Has(sdi.ExtKeyExempt) {
		return ext[sdi.ExtKeyExempt].String()
	}
	if ext.Has("it-sdi-nature") { // old key
		return ext["it-sdi-nature"].String()
	}
	return ""
}

func generateTaxSummary(inv *bill.Invoice) []*TaxSummary {
	var dr []*TaxSummary
	var vatRateTotals []*tax.RateTotal

	for _, cat := range inv.Totals.Taxes.Categories {
		if cat.Code == tax.CategoryVAT {
			vatRateTotals = cat.Rates
			break
		}
	}

	for _, rateTotal := range vatRateTotals {
		// Get tax liability from extensions if present
		var taxLiability string
		if rateTotal.Ext != nil && rateTotal.Ext.Has(sdi.ExtKeyVATLiability) {
			taxLiability = rateTotal.Ext[sdi.ExtKeyVATLiability].String()
		}

		dr = append(dr, &TaxSummary{
			TaxRate:        formatPercentageWithZero(rateTotal.Percent),
			TaxNature:      exemptExtensionCode(rateTotal.Ext),
			TaxableAmount:  formatAmount2(&rateTotal.Base),
			TaxAmount:      formatAmount2(&rateTotal.Amount),
			TaxLiability:   taxLiability,
			LegalReference: findRiferimentoNormativo(rateTotal),
		})
	}

	return dr
}

func extractLinePriceAdjustments(line *bill.Line) []*PriceAdjustment {
	list := make([]*PriceAdjustment, 0)

	for _, discount := range line.Discounts {
		// Unlike most formats, FatturaPA applies the discount to the unit price
		// instead of the line sum.
		// Quick ref: https://fex-app.com/FatturaElettronica/FatturaElettronicaBody/DatiBeniServizi/DettaglioLinee/PrezzoTotale
		a := discount.Amount
		if line.Quantity.Value() != 1 {
			a = a.RescaleUp(4).Divide(line.Quantity)
		}
		list = append(list, &PriceAdjustment{
			Type:    scontoMaggiorazioneTypeDiscount,
			Percent: formatPercentage(discount.Percent),
			Amount:  formatAmount8(&a),
		})
	}

	for _, charge := range line.Charges {
		a := charge.Amount
		if line.Quantity.Value() != 1 {
			a = a.RescaleUp(4).Divide(line.Quantity)
		}
		list = append(list, &PriceAdjustment{
			Type:    scontoMaggiorazioneTypeCharge,
			Percent: formatPercentage(charge.Percent),
			Amount:  formatAmount8(&a),
		})
	}

	return list
}

func findRiferimentoNormativo(rateTotal *tax.RateTotal) string {
	def := tax.ExtensionForKey(sdi.ExtKeyExempt)

	nature := exemptExtensionCode(rateTotal.Ext)
	for _, c := range def.Values {
		if c.Code.String() == nature {
			return c.Name[i18n.IT]
		}
	}

	return ""
}
