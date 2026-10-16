package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.it.sdi/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Under split payment (scissione dei pagamenti) the customer pays the VAT
// straight to the tax authority, so the supplier is only ever owed the net.
// ImportoPagamento must reflect that while ImportoTotaleDocumento stays gross.
func TestPaymentsSplitPayment(t *testing.T) {
	t.Run("percent due date is net of split payment VAT", func(t *testing.T) {
		env := test.LoadTestFile("invoice-split-payment.json", test.PathGOBLFatturaPA)
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		assert.Equal(t, "1388.40", doc.Body[0].GeneralData.Document.TotalAmount)

		dp := doc.Body[0].PaymentsData[0]
		assert.Equal(t, "TP02", dp.Conditions)
		require.Len(t, dp.Payments, 1)
		assert.Equal(t, "1032.00", dp.Payments[0].Amount)
		assert.Equal(t, "2023-04-30", dp.Payments[0].DueDate)
	})

	t.Run("explicit net due amount is a full payment, not an installment", func(t *testing.T) {
		env := test.LoadTestFile("invoice-split-payment.json", test.PathGOBLFatturaPA)
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			net := num.MakeAmount(103200, 2)
			inv.Payment.Terms.DueDates[0].Percent = nil
			inv.Payment.Terms.DueDates[0].Amount = &net
		})
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dp := doc.Body[0].PaymentsData[0]
		assert.Equal(t, "TP02", dp.Conditions)
		require.Len(t, dp.Payments, 1)
		assert.Equal(t, "1032.00", dp.Payments[0].Amount)
	})

	t.Run("percent installments are each their own rounded share of the net", func(t *testing.T) {
		env := test.LoadTestFile("invoice-split-payment.json", test.PathGOBLFatturaPA)
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			// 100.03 + 22% = 122.04 with 22.01 of VAT, so each half of the
			// net 100.03 rounds on its own, as GOBL rounds percentage dues.
			inv.Lines = inv.Lines[:1]
			inv.Lines[0].Quantity = num.MakeAmount(1, 0)
			inv.Lines[0].Item.Price = num.NewAmount(10003, 2)
			inv.Lines[0].Discounts = nil
			inv.Discounts = nil
			inv.Charges = nil
			inv.Payment.Terms.DueDates = []*pay.DueDate{
				{Date: cal.NewDate(2023, 4, 30), Percent: num.NewPercentage(50, 2)},
				{Date: cal.NewDate(2023, 5, 31), Percent: num.NewPercentage(50, 2)},
			}
		})
		require.NoError(t, env.Calculate())
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		assert.Equal(t, "122.04", doc.Body[0].GeneralData.Document.TotalAmount)

		dp := doc.Body[0].PaymentsData[0]
		assert.Equal(t, "TP01", dp.Conditions)
		require.Len(t, dp.Payments, 2)
		assert.Equal(t, "50.02", dp.Payments[0].Amount)
		assert.Equal(t, "50.02", dp.Payments[1].Amount)
	})

	t.Run("an advance comes off the net before a percent due date", func(t *testing.T) {
		env := test.LoadTestFile("invoice-split-payment.json", test.PathGOBLFatturaPA)
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Payment.Advances = []*pay.Record{{
				Description: "Deposit",
				Amount:      num.MakeAmount(10000, 2),
			}}
		})
		require.NoError(t, env.Calculate())
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		// 1388.40 gross, 356.40 split VAT, 100.00 already paid.
		dp := doc.Body[0].PaymentsData[1]
		assert.Equal(t, "TP02", dp.Conditions)
		require.Len(t, dp.Payments, 1)
		assert.Equal(t, "932.00", dp.Payments[0].Amount)
	})

	t.Run("an advance comes off the net when there are no due dates", func(t *testing.T) {
		env := test.LoadTestFile("invoice-split-payment.json", test.PathGOBLFatturaPA)
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Payment.Terms = nil
			inv.Payment.Advances = []*pay.Record{{
				Description: "Deposit",
				Amount:      num.MakeAmount(10000, 2),
			}}
		})
		require.NoError(t, env.Calculate())
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dp := doc.Body[0].PaymentsData[1]
		assert.Equal(t, "TP02", dp.Conditions)
		require.Len(t, dp.Payments, 1)
		assert.Equal(t, "932.00", dp.Payments[0].Amount)
	})

	t.Run("without due dates the single payment is net of split payment VAT", func(t *testing.T) {
		env := test.LoadTestFile("invoice-split-payment.json", test.PathGOBLFatturaPA)
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Payment.Terms = nil
		})
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dp := doc.Body[0].PaymentsData[0]
		assert.Equal(t, "TP02", dp.Conditions)
		require.Len(t, dp.Payments, 1)
		assert.Equal(t, "1032.00", dp.Payments[0].Amount)
	})
}
