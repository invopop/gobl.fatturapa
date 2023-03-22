package fatturapa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

const (
	keyRegimeFiscale    cbc.Key = it.KeyFatturaPARegimeFiscale
	keyTipoDocumento    cbc.Key = it.KeyFatturaPATipoDocumento
	keyNatura           cbc.Key = it.KeyFatturaPANatura
	keyTipoRitenuta     cbc.Key = it.KeyFatturaPATipoRitenuta
	keyCausalePagamento cbc.Key = it.KeyFatturaPACausalePagamento
)

func findCodeRegimeFiscale(inv *bill.Invoice) string {
	ss := inv.ScenarioSummary()

	return ss.Meta[keyRegimeFiscale]
}

func findCodeTipoDocumento(inv *bill.Invoice) string {
	ss := inv.ScenarioSummary()

	return ss.Meta[keyTipoDocumento]
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
			return tag.Meta[keyNatura]
		}
	}

	return ""
}

func findCodeTipoRitenuta(tc cbc.Code) string {
	taxCategory := regime.Category(tc)

	return taxCategory.Meta[keyTipoRitenuta]
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
					return tag.Meta[keyCausalePagamento]
				}
			}
		}
	}

	return ""
}
