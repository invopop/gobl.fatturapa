package fatturapa

import (
	"fmt"

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
	amount := formatAmount(&rateTotal.Amount)

	codeTR, err := findCodeTipoRitenuta(cat)
	if err != nil {
		return nil, err
	}

	return &datiRitenuta{
		TipoRitenuta:     codeTR,
		ImportoRitenuta:  amount,
		AliquotaRitenuta: rate,
		CausalePagamento: rateTotal.Ext[it.ExtKeySDIRetainedTax].String(),
	}, nil
}

func findCodeTipoRitenuta(cat cbc.Code) (string, error) {
	taxCategory := regime.Category(cat)

	code := taxCategory.Map[it.KeyFatturaPATipoRitenuta]

	if code == "" {
		return "", fmt.Errorf("could not find TipoRitenuta code for tax category %s", cat)
	}

	return code.String(), nil
}
