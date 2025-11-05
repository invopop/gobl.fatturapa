package fatturapa

import (
	"fmt"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

// RetainedTax represents a retained tax.
type RetainedTax struct {
	Type   string `xml:"TipoRitenuta"`
	Amount string `xml:"ImportoRitenuta"`
	Rate   string `xml:"AliquotaRitenuta"`
	Reason string `xml:"CausalePagamento"`
}

func extractRetainedTaxes(inv *bill.Invoice) ([]*RetainedTax, error) {
	catTotals := findRetainedCategories(inv.Totals)
	var dr []*RetainedTax

	for _, catTotal := range catTotals {
		for _, rateTotal := range catTotal.Rates {
			drElem, err := newRetainedTax(catTotal.Code, rateTotal)
			if err != nil {
				return nil, err
			}
			dr = append(dr, drElem)
		}
	}

	return dr, nil
}

func findRetainedCategories(totals *bill.Totals) []*tax.CategoryTotal {
	var catTotals []*tax.CategoryTotal

	for _, catTotal := range totals.Taxes.Categories {
		if catTotal.Retained {
			catTotals = append(catTotals, catTotal)
		}
	}

	return catTotals
}

func newRetainedTax(cat cbc.Code, rateTotal *tax.RateTotal) (*RetainedTax, error) {
	rate := formatPercentage(rateTotal.Percent)
	amount := formatAmount2(&rateTotal.Amount)

	codeTR, err := findCodeTaxType(cat)
	if err != nil {
		return nil, err
	}

	return &RetainedTax{
		Type:   codeTR,
		Amount: amount,
		Rate:   rate,
		Reason: retainedExtensionCode(rateTotal.Ext),
	}, nil
}

func retainedExtensionCode(ext tax.Extensions) string {
	if ext.Has(sdi.ExtKeyRetained) {
		return ext[sdi.ExtKeyRetained].String()
	}
	if ext.Has("it-sdi-retained-tax") { // old key
		return ext["it-sdi-retained-tax"].String()
	}
	return ""
}

func findCodeTaxType(cat cbc.Code) (string, error) {
	switch cat {
	case it.TaxCategoryIRPEF:
		return "RT01", nil
	case it.TaxCategoryIRES:
		return "RT02", nil
	case it.TaxCategoryINPS:
		return "RT03", nil
	case it.TaxCategoryENASARCO:
		return "RT04", nil
	case it.TaxCategoryENPAM:
		return "RT05", nil
	case it.TaxCategoryCP:
		return "RT06", nil
	default:
		return "", fmt.Errorf("could not find TipoRitenuta code for tax category %s", cat)
	}
}
