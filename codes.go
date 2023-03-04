package fatturapa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

const (
	regimeFiscaleCodeDefault = "RF01"
	tipoDocumentoCodeDefault = "TD01"
)

var regime = tax.RegimeFor(l10n.IT, l10n.CodeEmpty)

func findCodeRegimeFiscale(inv *bill.Invoice) string {
	t := inv.Type.Key()

	// based on the invoice type and the tax schemes used in the invoice,
	// it iterates through the available schemes defined in gobl/it and fetches
	// the corresponding regime fiscale code
	for _, scheme := range regime.Schemes {
		for _, invType := range scheme.InvoiceTypes {
			for _, s := range inv.Tax.Schemes {
				if t == invType && s == scheme.Key {
					return scheme.Meta[it.KeyFatturaPARegimeFiscale]
				}
			}
		}
	}

	return regimeFiscaleCodeDefault
}

func findCodeTipoDocumento(inv *bill.Invoice) string {
	t := inv.Type.Key()

	for _, scheme := range regime.Schemes {
		for _, invType := range scheme.InvoiceTypes {
			for _, s := range inv.Tax.Schemes {
				if invType == t && scheme.Key == s {
					return scheme.Meta[it.KeyFatturaPATipoDocumento]
				}
			}
		}
	}

	return tipoDocumentoCodeDefault
}

func findCodeNatura(line *bill.Line) string {
	var taxCategoryVAT *tax.Category

	// get the italian VAT tax category
	for _, cat := range regime.Categories {
		if cat.Code == common.TaxCategoryVAT {
			taxCategoryVAT = cat
			break
		}
	}

	var vatZeroTags []*tax.Tag

	// get the list of available tags for the zero VAT rate
	for _, rate := range taxCategoryVAT.Rates {
		if rate.Key == common.TaxRateZero {
			vatZeroTags = rate.Tags
			break
		}
	}

	if vatZeroTags == nil {
		return ""
	}

	// check if the line has a tag for the zero VAT rate and return the
	// corresponding code from the tag metadata
	for _, tag := range vatZeroTags {
		for _, tax := range line.Taxes {
			if tax.Tag == tag.Key {
				return tag.Meta[it.KeyFatturaPANatura]
			}
		}
	}

	return ""
}
