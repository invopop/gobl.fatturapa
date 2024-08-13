package fatturapa

import "github.com/invopop/gobl/num"

func formatPercentage(p *num.Percentage) string {
	if p == nil {
		return ""
	}
	return p.Rescale(4).StringWithoutSymbol()
}

func formatPercentageWithZero(p *num.Percentage) string {
	if p == nil {
		p = num.NewPercentage(0, 4)
	}
	return formatPercentage(p)
}

func formatAmount(a *num.Amount) string {
	if a == nil {
		return ""
	}
	return a.RescaleUp(2).String()
}
