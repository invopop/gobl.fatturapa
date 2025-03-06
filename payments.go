package fatturapa

import (
	"fmt"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/pay"
)

// PaymentData contains all data related to the payment of the document.
type PaymentData struct {
	Conditions string              `xml:"CondizioniPagamento"`
	Payments   []*PaymentDetailRow `xml:"DettaglioPagamento,omitempty"`
}

// PaymentDetailRow contains data related to a single payment.
type PaymentDetailRow struct {
	Beneficiary string `xml:"Beneficiario,omitempty"`
	Method      string `xml:"ModalitaPagamento"`
	Date        string `xml:"DataRiferimentoTerminiPagamento,omitempty"`
	Days        int64  `xml:"GiorniTerminiPagamento,omitempty"`
	DueDate     string `xml:"DataScadenzaPagamento,omitempty"`
	Amount      string `xml:"ImportoPagamento"`
	IBAN        string `xml:"IBAN,omitempty"`
	ABI         string `xml:"ABI,omitempty"`
	CAB         string `xml:"CAB,omitempty"`
	BIC         string `xml:"BIC,omitempty"`
	Code        string `xml:"CodicePagamento,omitempty"`
}

func newPaymentData(inv *bill.Invoice) (*PaymentData, error) {
	if inv.Payment == nil {
		return nil, nil
	}

	dp, err := preparePaymentDetails(inv)
	if err != nil {
		return nil, err
	}

	return &PaymentData{
		Conditions: determinePaymentConditions(inv),
		Payments:   dp,
	}, nil
}

func preparePaymentDetails(inv *bill.Invoice) ([]*PaymentDetailRow, error) {
	var dp []*PaymentDetailRow
	payment := inv.Payment

	if len(payment.Advances) == 0 && payment.Instructions == nil {
		return nil, fmt.Errorf("missing payment advances or instructions")
	}

	// First deal with payment advances
	for _, advance := range payment.Advances {
		row := &PaymentDetailRow{
			Method: advance.Ext[sdi.ExtKeyPaymentMeans].String(),
			Amount: formatAmount2(&advance.Amount),
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

	br := PaymentDetailRow{
		Method: payment.Instructions.Ext[sdi.ExtKeyPaymentMeans].String(),
	}
	if len(payment.Instructions.CreditTransfer) > 0 {
		ct1 := payment.Instructions.CreditTransfer[0]
		br.IBAN = ct1.IBAN
		br.BIC = ct1.BIC
	}

	// First check if there are multiple due dates, and if so, create a
	// DettaglioPagamento for each one.
	if terms := payment.Terms; terms != nil {
		for _, dueDate := range payment.Terms.DueDates {
			r := br // copy
			r.DueDate = dueDate.Date.String()
			r.Amount = formatAmount2(&dueDate.Amount)
			dp = append(dp, &r)
		}
	}

	// If there are no due dates, then a single DettaglioPagamento is created
	// with the total payable amount.
	if len(dp) == 0 {
		br.Amount = formatAmount2(&inv.Totals.Payable)
		dp = append(dp, &br)
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
