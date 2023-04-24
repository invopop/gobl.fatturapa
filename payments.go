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

type DatiPagamento struct {
	CondizioniPagamento string
	DettaglioPagamento  []*DettaglioPagamento
}

type DettaglioPagamento struct {
	ModalitaPagamento     string
	DataScadenzaPagamento string `xml:",omitempty"`
	ImportoPagamento      string
}

func newDatiPagamento(inv *bill.Invoice) *DatiPagamento {
	if inv.Payment == nil {
		return nil
	}

	return &DatiPagamento{
		CondizioniPagamento: determinePaymentConditions(inv.Payment),
		DettaglioPagamento:  newDettalgioPagamento(inv),
	}
}

func newDettalgioPagamento(inv *bill.Invoice) []*DettaglioPagamento {
	var dp []*DettaglioPagamento
	payment := inv.Payment

	paymentMethod := paymentMethods[payment.Instructions.Key]

	// First check if there are multiple due dates, and if so, create a
	// DettaglioPagamento for each one.
	if terms := payment.Terms; terms != nil {
		for _, dueDate := range payment.Terms.DueDates {
			dp = append(dp, &DettaglioPagamento{
				ModalitaPagamento:     paymentMethod,
				DataScadenzaPagamento: dueDate.Date.String(), // ISO 8601 YYYY-MM-DD format
				ImportoPagamento:      dueDate.Amount.String(),
			})
		}
	}

	// If there are no due dates, then a single DettaglioPagamento is created
	// with the total payable amount.
	if len(dp) == 0 {
		dp = append(dp, &DettaglioPagamento{
			ModalitaPagamento: paymentMethod,
			ImportoPagamento:  inv.Totals.Payable.String(),
		})
	}

	return dp
}

func determinePaymentConditions(payment *bill.Payment) string {
	switch {
	case payment.Terms == nil:
		return CondizioniPagamentoFull
	case len(payment.Terms.DueDates) > 1:
		return CondizioniPagamentoInstallments
	case payment.Terms.Key == pay.TermKeyAdvance:
		return CondizioniPagamentoAdvance
	default:
		return CondizioniPagamentoFull
	}
}
