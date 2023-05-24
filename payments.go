package fatturapa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
)

var paymentMethods = map[cbc.Key]string{
	pay.MethodKeyCash:           "MP01",
	pay.MethodKeyCheque:         "MP02",
	pay.MethodKeyBankDraft:      "MP03",
	pay.MethodKeyCreditTransfer: "MP05",
	pay.MethodKeyCard:           "MP08",
	pay.MethodKeyDirectDebit:    "MP09",
	pay.MethodKeyDebitTransfer:  "MP09",
}

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

func newDatiPagamento(inv *bill.Invoice) *datiPagamento {
	if inv.Payment == nil {
		return nil
	}

	return &datiPagamento{
		CondizioniPagamento: determinePaymentConditions(inv.Payment),
		DettaglioPagamento:  newDettalgioPagamento(inv),
	}
}

func newDettalgioPagamento(inv *bill.Invoice) []*dettaglioPagamento {
	var dp []*dettaglioPagamento
	payment := inv.Payment

	paymentMethod := paymentMethods[payment.Instructions.Key]

	// First check if there are multiple due dates, and if so, create a
	// DettaglioPagamento for each one.
	if terms := payment.Terms; terms != nil {
		for _, dueDate := range payment.Terms.DueDates {
			dp = append(dp, &dettaglioPagamento{
				ModalitaPagamento:     paymentMethod,
				DataScadenzaPagamento: dueDate.Date.String(), // ISO 8601 YYYY-MM-DD format
				ImportoPagamento:      formatAmount(&dueDate.Amount),
			})
		}
	}

	// If there are no due dates, then a single DettaglioPagamento is created
	// with the total payable amount.
	if len(dp) == 0 {
		dp = append(dp, &dettaglioPagamento{
			ModalitaPagamento: paymentMethod,
			ImportoPagamento:  formatAmount(&inv.Totals.Payable),
		})
	}

	return dp
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
