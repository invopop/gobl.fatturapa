package fatturapa

import "github.com/invopop/gobl/num"

func formatPercentage(p *num.Percentage) string {
	if p == nil {
		return num.MakePercentage(0, 4).StringWithoutSymbol()
	}

	return p.Rescale(4).StringWithoutSymbol()
}

func formatAmount(a *num.Amount) string {
	if a == nil {
		return ""
	}
	return a.RescaleUp(2).String()
}
