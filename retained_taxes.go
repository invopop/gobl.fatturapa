package fatturapa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

type DatiRitenuta struct {
	TipoRitenuta     string
	ImportoRitenuta  string
	AliquotaRitenuta string
	CausalePagamento string
}

func extractRetainedTaxes(inv *bill.Invoice) []*DatiRitenuta {
	var dr []*DatiRitenuta
	var retCats []cbc.Code

	// First we need to find all the retained tax categoriesfrom Totals
	for _, tax := range inv.Totals.Taxes.Categories {
		if tax.Retained {
			retCats = append(retCats, tax.Code)
		}
	}

	// Then we iterate through the invoice lines and build DatiRitenuta taking
	// into account the attached tags
	for _, line := range inv.Lines {
		for _, tax := range line.Taxes {
			if !includesCode(retCats, tax.Category) {
				continue
			}

			codeTR := findCodeTipoRitenuta(tax.Category)
			amount := tax.Percent.Multiply(line.Total).Rescale(2).String()
			rate := tax.Percent.String()
			codeCP := findCodeCausalePagamento(line, tax.Category)

			dr = append(dr, &DatiRitenuta{
				TipoRitenuta:     codeTR,
				ImportoRitenuta:  amount,
				AliquotaRitenuta: rate,
				CausalePagamento: codeCP,
			})
		}
	}

	return dr
}

func includesCode(codes []cbc.Code, code cbc.Code) bool {
	for _, c := range codes {
		if c == code {
			return true
		}
	}

	return false
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
