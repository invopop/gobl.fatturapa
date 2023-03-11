package fatturapa

import (
	"github.com/invopop/gobl/bill"
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

func findCodeNatura(line *bill.Line) string {
	return ""
}
