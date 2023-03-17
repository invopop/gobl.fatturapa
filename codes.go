package fatturapa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/regimes/it"
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

func findCodeCausalePagamento(inv *bill.Invoice, tc cbc.Code) string {
	taxCategory := regime.Category(tc)

	for _, line := range inv.Lines {
		for _, tax := range line.Taxes {
			if tax.Category == tc {
				if len(tax.Tags) == 0 {
					continue
				}

				for _, tag := range taxCategory.Tags {
					for _, t := range tax.Tags {
						if tag.Key == t {
							return tag.Meta[it.KeyFatturaPACausalePagamento]
						}
					}
				}
			}
		}
	}

	return ""
}
