package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDettaglioLinee(t *testing.T) {
	t.Run("should contain the line info", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json", test.PathGOBLFatturaPA)
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dl := doc.Body[0].GoodsServices.LineDetails[0]

		assert.Equal(t, "1", dl.LineNumber)
		assert.Equal(t, "Development services", dl.Description)
		assert.Equal(t, "20.00", dl.Quantity)
		assert.Equal(t, "90.00", dl.UnitPrice)
		assert.Equal(t, "1620.00", dl.TotalPrice)
		assert.Equal(t, "22.00", dl.TaxRate)
		assert.Equal(t, "", dl.TaxNature)

		sm := dl.PriceAdjustments[0]

		assert.Equal(t, "SC", sm.Type)
		assert.Equal(t, "10.00", sm.Percent)
		assert.Equal(t, "9.0000", sm.Amount)

		dl = doc.Body[0].GoodsServices.LineDetails[1]

		assert.Equal(t, "2", dl.LineNumber)
		assert.Equal(t, "N2.2", dl.TaxNature)
	})
}

func TestDettaglioLineePeriod(t *testing.T) {
	t.Run("should map period dates from GOBL to FatturaPA", func(t *testing.T) {
		env := test.LoadTestFile("invoice-services-period.json", test.PathGOBLFatturaPA)
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dl := doc.Body[0].GoodsServices.LineDetails[0]
		assert.Equal(t, "2024-01-01", dl.PeriodStart)
		assert.Equal(t, "2024-01-31", dl.PeriodEnd)

		// Line without period should have empty dates
		dl2 := doc.Body[0].GoodsServices.LineDetails[1]
		assert.Empty(t, dl2.PeriodStart)
		assert.Empty(t, dl2.PeriodEnd)
	})

	t.Run("should omit missing period date", func(t *testing.T) {
		env := test.LoadTestFile("invoice-services-period.json", test.PathGOBLFatturaPA)
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Lines[0].Period.End = cal.Date{}
		})
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dl := doc.Body[0].GoodsServices.LineDetails[0]
		assert.Equal(t, "2024-01-01", dl.PeriodStart)
		assert.Empty(t, dl.PeriodEnd)
	})
}

func TestAltriDatiGestionaliINVCONT(t *testing.T) {
	t.Run("should add INVCONT for N2.1 with reverse-charge tag", func(t *testing.T) {
		env := test.LoadTestFile("invoice-reverse-charge.json", test.PathGOBLFatturaPA)
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		for _, dl := range doc.Body[0].GoodsServices.LineDetails {
			require.Len(t, dl.OtherData, 1)
			assert.Equal(t, "INVCONT", dl.OtherData[0].DataType)
			assert.Equal(t, "Inversione contabile - art. 21 c.6 bis lett. a) DPR 633/72", dl.OtherData[0].TextReference)
		}
	})

	t.Run("should not add INVCONT for N2.1 without reverse-charge tag", func(t *testing.T) {
		env := test.LoadTestFile("invoice-hotel.json", test.PathGOBLFatturaPA)
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		// First line has N2.1 but no reverse-charge tag on the invoice
		dl := doc.Body[0].GoodsServices.LineDetails[0]
		assert.Equal(t, "N2.1", dl.TaxNature)
		assert.Empty(t, dl.OtherData)
	})
}

func TestDatiRiepilogo(t *testing.T) {
	t.Run("should contain the tax summary info", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json", test.PathGOBLFatturaPA)
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dr := doc.Body[0].GoodsServices.TaxSummary[0]

		assert.Equal(t, "22.00", dr.TaxRate)
		assert.Equal(t, "1620.00", dr.TaxableAmount)
		assert.Equal(t, "356.40", dr.TaxAmount)
		assert.Equal(t, "", dr.TaxNature)
		assert.Equal(t, "", dr.LegalReference)

		dr = doc.Body[0].GoodsServices.TaxSummary[1]

		assert.Equal(t, "N2.2", dr.TaxNature)
		assert.Equal(t, "S", dr.TaxLiability)
		assert.Equal(t, "Non soggette - altri casi", dr.LegalReference)
	})
}

func TestNegativeQuantityConversion(t *testing.T) {
	t.Run("should convert negative quantities to negative prices", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json", test.PathGOBLFatturaPA)

		// Modify the invoice to have a negative quantity
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Lines[0].Quantity = num.MakeAmount(-2000, 2) // -20.00
			require.NoError(t, inv.Calculate())
		})

		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dl := doc.Body[0].GoodsServices.LineDetails[0]

		// Quantity should be positive
		assert.Equal(t, "20.00", dl.Quantity)

		// Unit price should be negative
		assert.Equal(t, "-90.00", dl.UnitPrice)

		// Total price should still be negative
		assert.Equal(t, "-1620.00", dl.TotalPrice)
	})

	t.Run("should handle discounts correctly with negative quantities", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json", test.PathGOBLFatturaPA)

		// Modify the invoice to have a negative quantity
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Lines[0].Quantity = num.MakeAmount(-2000, 2) // -20.00
			require.NoError(t, inv.Calculate())
		})

		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dl := doc.Body[0].GoodsServices.LineDetails[0]

		// Price adjustments (discounts/charges) should have positive amounts
		// regardless of quantity sign
		require.NotEmpty(t, dl.PriceAdjustments)

		// The discount amount should be positive (9.0000), not negative
		// Original line has 10% discount = 180 total / 20 units = 9 per unit
		assert.Equal(t, "SC", dl.PriceAdjustments[0].Type)
		assert.Equal(t, "10.00", dl.PriceAdjustments[0].Percent)
		assert.Equal(t, "9.0000", dl.PriceAdjustments[0].Amount)
	})
}
