package fatturapa

import (
	"fmt"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

type datiRitenuta struct {
	TipoRitenuta     string
	ImportoRitenuta  string
	AliquotaRitenuta string
	CausalePagamento string
}

func extractRetainedTaxes(inv *bill.Invoice) ([]*datiRitenuta, error) {
	catTotals := findRetainedCategories(inv.Totals)
	var dr []*datiRitenuta

	for _, catTotal := range catTotals {
		for _, rateTotal := range catTotal.Rates {
			drElem, err := newDatiRitenuta(catTotal.Code, rateTotal)
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

func newDatiRitenuta(cat cbc.Code, rateTotal *tax.RateTotal) (*datiRitenuta, error) {
	rate := formatPercentage(rateTotal.Percent)
	amount := formatAmount2(&rateTotal.Amount)

	codeTR, err := findCodeTipoRitenuta(cat)
	if err != nil {
		return nil, err
	}

	return &datiRitenuta{
		TipoRitenuta:     codeTR,
		ImportoRitenuta:  amount,
		AliquotaRitenuta: rate,
		CausalePagamento: retainedExtensionCode(rateTotal.Ext),
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

func findCodeTipoRitenuta(cat cbc.Code) (string, error) {
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
	default:
		return "", fmt.Errorf("could not find TipoRitenuta code for tax category %s", cat)
	}
}
