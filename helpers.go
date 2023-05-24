package fatturapa

import "github.com/invopop/gobl/num"

func formatPercentage(p *num.Percentage) string {
	if p == nil {
		return ""
	}

	return p.Amount.Multiply(*num.NewAmount(100, 0)).Rescale(2).String()
}
