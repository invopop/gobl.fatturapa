package fatturapa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

func findCodeRegimeFiscale(inv *bill.Invoice) string {
	ss := inv.ScenarioSummary()

	return ss.Meta[it.KeyFatturaPARegimeFiscale]
}

func findCodeTipoDocumento(inv *bill.Invoice) string {
	ss := inv.ScenarioSummary()

	return ss.Meta[it.KeyFatturaPATipoDocumento]
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

func findCodeTipoRitenuta(tc cbc.Code) string {
	taxCategory := regime.Category(tc)

	return taxCategory.Meta[it.KeyFatturaPATipoRitenuta]
}

func findCodeCausalePagamento(line *bill.Line, tc cbc.Code) string {
	taxCategory := regime.Category(tc)
	var lineTaxes []tax.Combo

	for _, lt := range line.Taxes {
		if lt.Category == tc {
			lineTaxes = append(lineTaxes, *lt)
		}
	}

	for _, lt := range lineTaxes {
		if len(lt.Tags) == 0 {
			continue
		}

		for _, tag := range taxCategory.Tags {
			for _, t := range lt.Tags {
				if tag.Key == t {
					return tag.Meta[it.KeyFatturaPACausalePagamento]
				}
			}
		}
	}

	return ""
}
