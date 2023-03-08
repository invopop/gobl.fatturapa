package fatturapa

import (
	"github.com/invopop/gobl/bill"
)

const (
	regimeFiscaleCodeDefault = "RF01"
	tipoDocumentoCodeDefault = "TD01"
)

func findCodeRegimeFiscale(inv *bill.Invoice) string {
	return regimeFiscaleCodeDefault
}

func findCodeTipoDocumento(inv *bill.Invoice) string {
	return tipoDocumentoCodeDefault
}

func findCodeNatura(line *bill.Line) string {
	return ""
}
