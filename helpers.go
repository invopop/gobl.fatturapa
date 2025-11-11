package fatturapa

import (
	"strings"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
)

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

// formatAmount2 will format the number always with 2 decimal places.
func formatAmount2(a *num.Amount) string {
	if a == nil {
		return ""
	}
	return a.Rescale(2).String()
}

// formatAmount8 will ensure the number has at least 2 decimal places,
// with a maximum of 8.
func formatAmount8(a *num.Amount) string {
	if a == nil {
		return ""
	}
	return a.RescaleRange(2, 8).String()
}

// parseDate will parse a date string and return a cal.Date
func parseDate(date string) (cal.Date, error) {
	if date == "" {
		return cal.Date{}, nil
	}
	var d cal.Date
	if err := d.UnmarshalJSON([]byte(`"` + date + `"`)); err != nil {
		return cal.Date{}, err
	}
	return d, nil
}

// parseAmount will trim whitespace and parse a string into a num.Amount
func parseAmount(s string) (num.Amount, error) {
	if s == "" {
		return num.Amount{}, nil
	}
	return num.AmountFromString(strings.TrimSpace(s))
}
