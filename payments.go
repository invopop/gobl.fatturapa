package fatturapa

import (
	"fmt"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/pay"
)

// datiPagamento contains all data related to the payment of the document.
type datiPagamento struct {
	CondizioniPagamento string
	DettaglioPagamento  []*dettaglioPagamento
}

// dettaglioPagamento contains data related to a single payment.
type dettaglioPagamento struct {
	ModalitaPagamento string
	Date              string `xml:"DataRiferimentoTerminiPagamento,omitempty"`
	DueDate           string `xml:"DataScadenzaPagamento,omitempty"`
	ImportoPagamento  string
}

func newDatiPagamento(inv *bill.Invoice) (*datiPagamento, error) {
	if inv.Payment == nil {
		return nil, nil
	}

	dp, err := preparePaymentDetails(inv)
	if err != nil {
		return nil, err
	}

	return &datiPagamento{
		CondizioniPagamento: determinePaymentConditions(inv),
		DettaglioPagamento:  dp,
	}, nil
}

func preparePaymentDetails(inv *bill.Invoice) ([]*dettaglioPagamento, error) {
	var dp []*dettaglioPagamento
	payment := inv.Payment

	if len(payment.Advances) == 0 && payment.Instructions == nil {
		return nil, fmt.Errorf("missing payment advances or instructions")
	}

	// First deal with payment advances
	for _, advance := range payment.Advances {
		row := &dettaglioPagamento{
			ModalitaPagamento: advance.Ext[sdi.ExtKeyPaymentMeans].String(),
			ImportoPagamento:  formatAmount(&advance.Amount),
		}
		if advance.Date != nil {
			row.Date = advance.Date.String()
		}
		dp = append(dp, row)
	}

	if payment.Instructions == nil {
		// No instructions, ignore anything else
		return dp, nil
	}

	codeModalitaPagamento := payment.Instructions.Ext[sdi.ExtKeyPaymentMeans].String()

	// First check if there are multiple due dates, and if so, create a
	// DettaglioPagamento for each one.
	if terms := payment.Terms; terms != nil {
		for _, dueDate := range payment.Terms.DueDates {
			dp = append(dp, &dettaglioPagamento{
				ModalitaPagamento: codeModalitaPagamento,
				DueDate:           dueDate.Date.String(), // ISO 8601 YYYY-MM-DD format
				ImportoPagamento:  formatAmount(&dueDate.Amount),
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

func determinePaymentConditions(inv *bill.Invoice) string {
	p := inv.Payment
	switch {
	case inv.Totals.Paid() || (p.Terms != nil && p.Terms.Key == pay.TermKeyAdvanced):
		return condizioniPagamentoAdvance
	case len(p.Advances) > 0 || (p.Terms != nil && len(p.Terms.DueDates) > 1):
		return condizioniPagamentoInstallments
	default:
		return condizioniPagamentoFull
	}
}
