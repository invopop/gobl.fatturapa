package fatturapa

import (
	"fmt"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
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

func newPaymentData(inv *bill.Invoice) ([]*PaymentData, error) {
	fmt.Println("Starting newPaymentData function")
	if inv.Payment == nil {
		fmt.Println("Invoice payment is nil, returning nil")
		return nil, nil
	}

	paymentData := []*PaymentData{}

	if inv.Payment.Advances != nil {
		fmt.Printf("Found %d advances, adding to payment data\n", len(inv.Payment.Advances))
		paymentData = append(paymentData, &PaymentData{
			Conditions: condizioniPagamentoAdvance,
			Payments:   prepareAdvancePaymentDetails(inv),
		})
	}

	if inv.Payment.Instructions == nil {
		fmt.Println("Payment instructions are nil, returning current payment data")
		return paymentData, nil
	}

	// Determine payment condition based on number of due dates
	condition := condizioniPagamentoFull
	if inv.Payment.Terms != nil && len(inv.Payment.Terms.DueDates) > 1 {
		fmt.Printf("Multiple due dates found (%d), using installments condition\n", len(inv.Payment.Terms.DueDates))
		condition = condizioniPagamentoInstallments
	} else {
		fmt.Println("Using full payment condition")
	}

	paymentDetails := preparePaymentDetails(inv)
	if len(paymentDetails) == 0 {
		fmt.Println("No payment details prepared, returning current payment data")
		return paymentData, nil
	}

	fmt.Printf("Adding payment data with condition %s and %d payment details\n", condition, len(paymentDetails))
	paymentData = append(paymentData, &PaymentData{
		Conditions: condition,
		Payments:   paymentDetails,
	})

	return paymentData, nil
}

func prepareAdvancePaymentDetails(inv *bill.Invoice) []*PaymentDetailRow {
	fmt.Println("Starting prepareAdvancePaymentDetails function")
	var dp []*PaymentDetailRow
	payment := inv.Payment

	for i, advance := range payment.Advances {
		fmt.Printf("Processing advance payment %d\n", i+1)
		row := &PaymentDetailRow{
			Method: advance.Ext[sdi.ExtKeyPaymentMeans].String(),
			Amount: formatAmount2(&advance.Amount),
		}
		fmt.Printf("Advance payment method: %s, amount: %s\n", row.Method, row.Amount)

		if advance.Date != nil {
			row.Date = advance.Date.String()
			fmt.Printf("Advance payment date: %s\n", row.Date)
		}
		if advance.CreditTransfer != nil {
			row.IBAN = advance.CreditTransfer.IBAN
			row.BIC = advance.CreditTransfer.BIC
			fmt.Printf("Advance payment credit transfer - IBAN: %s, BIC: %s\n", row.IBAN, row.BIC)
		}
		dp = append(dp, row)
	}

	fmt.Printf("Prepared %d advance payment details\n", len(dp))
	return dp
}

func preparePaymentDetails(inv *bill.Invoice) []*PaymentDetailRow {
	fmt.Println("Starting preparePaymentDetails function")
	var dp []*PaymentDetailRow
	payment := inv.Payment

	fmt.Printf("Payment method from instructions: %s\n", payment.Instructions.Ext[sdi.ExtKeyPaymentMeans].String())
	br := PaymentDetailRow{
		Method: payment.Instructions.Ext[sdi.ExtKeyPaymentMeans].String(),
	}
	if len(payment.Instructions.CreditTransfer) > 0 {
		ct1 := payment.Instructions.CreditTransfer[0]
		br.IBAN = ct1.IBAN
		br.BIC = ct1.BIC
		fmt.Printf("Credit transfer details - IBAN: %s, BIC: %s\n", br.IBAN, br.BIC)
	}

	// First check if there are multiple due dates, and if so, create a
	// DettaglioPagamento for each one.
	if terms := payment.Terms; terms != nil {
		fmt.Printf("Processing %d due dates\n", len(terms.DueDates))
		for i, dueDate := range payment.Terms.DueDates {
			fmt.Printf("Processing due date %d: %s\n", i+1, dueDate.Date.String())
			r := br // copy
			r.DueDate = dueDate.Date.String()
			r.Amount = formatAmount2(&dueDate.Amount)
			fmt.Printf("Due date %d details - date: %s, amount: %s\n", i+1, r.DueDate, r.Amount)
			dp = append(dp, &r)
		}
	}

	// If there are no due dates, then a single DettaglioPagamento is created
	// with the total payable amount.
	if len(dp) == 0 {
		fmt.Println("No due dates found, creating single payment detail with total payable amount")
		br.Amount = formatAmount2(&inv.Totals.Payable)
		fmt.Printf("Single payment detail amount: %s\n", br.Amount)
		dp = append(dp, &br)
	}

	fmt.Printf("Prepared %d payment details\n", len(dp))
	return dp
}
