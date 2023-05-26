package fatturapa

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

// datiPagamento contains all data related to the payment of the document.
type datiPagamento struct {
	CondizioniPagamento string
	DettaglioPagamento  []*dettaglioPagamento
}

// dettaglioPagamento contains data related to a single payment.
type dettaglioPagamento struct {
	ModalitaPagamento     string
	DataScadenzaPagamento string `xml:",omitempty"`
	ImportoPagamento      string
}

func newDatiPagamento(inv *bill.Invoice) (*datiPagamento, error) {
	if inv.Payment == nil {
		return nil, nil
	}

	dp, err := newDettalgioPagamento(inv)
	if err != nil {
		return nil, err
	}

	return &datiPagamento{
		CondizioniPagamento: determinePaymentConditions(inv.Payment),
		DettaglioPagamento:  dp,
	}, nil
}

func newDettalgioPagamento(inv *bill.Invoice) ([]*dettaglioPagamento, error) {
	var dp []*dettaglioPagamento
	payment := inv.Payment

	codeModalitaPagamento, err := findCodeModalitaPagamento(payment.Instructions.Key)
	if err != nil {
		return nil, err
	}

	// First check if there are multiple due dates, and if so, create a
	// DettaglioPagamento for each one.
	if terms := payment.Terms; terms != nil {
		for _, dueDate := range payment.Terms.DueDates {
			dp = append(dp, &dettaglioPagamento{
				ModalitaPagamento:     codeModalitaPagamento,
				DataScadenzaPagamento: dueDate.Date.String(), // ISO 8601 YYYY-MM-DD format
				ImportoPagamento:      formatAmount(&dueDate.Amount),
			})
		}
	}

	// If there are no due dates, then a single DettaglioPagamento is created
	// with the total payable amount.
	if len(dp) == 0 {
		dp = append(dp, &dettaglioPagamento{
			ModalitaPagamento: codeModalitaPagamento,
			ImportoPagamento:  formatAmount(&inv.Totals.Payable),
		})
	}

	return dp, nil
}

func findCodeModalitaPagamento(key cbc.Key) (string, error) {
	keyDef := findPaymentKeyDefinition(key)

	if keyDef == nil {
		return "", fmt.Errorf("ModalitaPagamento Code not found for payment method key '%s'", key)
	}

	code := keyDef.Codes[it.KeyFatturaPAModalitaPagamento]
	if code == "" {
		return "", fmt.Errorf("ModalitaPagamento Code not found for payment method key '%s'", key)
	}

	return code.String(), nil
}

func findPaymentKeyDefinition(key cbc.Key) *tax.KeyDefinition {
	for _, keyDef := range regime.PaymentMeansKeys {
		if key == keyDef.Key {
			return keyDef
		}
	}
	return nil
}

func determinePaymentConditions(payment *bill.Payment) string {
	switch {
	case payment.Terms == nil:
		return condizioniPagamentoFull
	case len(payment.Terms.DueDates) > 1:
		return condizioniPagamentoInstallments
	case payment.Terms.Key == pay.TermKeyAdvanced:
		return condizioniPagamentoAdvance
	default:
		return condizioniPagamentoFull
	}
}
