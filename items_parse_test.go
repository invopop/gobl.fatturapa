package fatturapa_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestItemsInConversion tests the conversion of line items and tax summaries from FatturaPA XML to GOBL
func TestItemsInConversion(t *testing.T) {
	t.Run("should convert line items correctly", func(t *testing.T) {
		// Load the XML file with line items
		data, err := os.ReadFile(filepath.Join(test.GetDataPath(test.PathFatturaPAGOBL), "invoice-simple.xml"))
		require.NoError(t, err)

		// Convert XML to GOBL
		env, err := test.ConvertToGOBL(data)
		require.NoError(t, err)
		require.NotNil(t, env)

		// Extract the invoice
		invoice, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		require.NotNil(t, invoice)

		// Check the line items
		require.NotEmpty(t, invoice.Lines)
		require.GreaterOrEqual(t, len(invoice.Lines), 2, "Invoice should have at least 2 lines")

		// Check the first line
		line1 := invoice.Lines[0]
		assert.Equal(t, "Development services", line1.Item.Name)
		assert.True(t, line1.Quantity.Compare(num.MakeAmount(2000, 2)) == 0, "Quantity should be 20.00")
		assert.True(t, line1.Item.Price.Compare(num.MakeAmount(9000, 2)) == 0, "Unit price should be 90.00")
		assert.True(t, line1.Total.Compare(num.MakeAmount(162000, 2)) == 0, "Total should be 1620.00")

		// Check the tax rate on the first line
		require.NotEmpty(t, line1.Taxes)
		taxRate := num.MakePercentage(220, 3) // 22.0% with correct precision
		require.NotNil(t, line1.Taxes[0].Percent)
		assert.True(t, line1.Taxes[0].Percent.Compare(taxRate) == 0, "Tax rate should be 22.0%")

		// Check for price adjustments (discounts)
		require.NotEmpty(t, line1.Discounts)
		discountPercent := num.MakePercentage(100, 3) // 10.0% with correct precision
		require.NotNil(t, line1.Discounts[0].Percent)
		assert.True(t, line1.Discounts[0].Percent.Compare(discountPercent) == 0, "Discount percent should be 10.0%")
		assert.True(t, line1.Discounts[0].Amount.Compare(num.MakeAmount(18000, 2)) == 0, "Discount amount should be 180.00")

		// Check the second line
		line2 := invoice.Lines[1]
		assert.NotEmpty(t, line2.Item.Name)

		// Check for tax exemption on the second line
		require.NotEmpty(t, line2.Taxes)
		assert.Equal(t, tax.CategoryVAT, line2.Taxes[0].Category)
		assert.Equal(t, cbc.Code("N1"), line2.Taxes[0].Ext[sdi.ExtKeyExempt])
	})

	t.Run("should convert tax summaries correctly", func(t *testing.T) {
		// Load the XML file with tax summaries
		data, err := os.ReadFile(filepath.Join(test.GetDataPath(test.PathFatturaPAGOBL), "invoice-simple.xml"))
		require.NoError(t, err)

		// Convert XML to GOBL
		env, err := test.ConvertToGOBL(data)
		require.NoError(t, err)
		require.NotNil(t, env)

		// Extract the invoice
		invoice, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		require.NotNil(t, invoice)

		// Check the tax totals
		require.NotNil(t, invoice.Totals)
		require.NotNil(t, invoice.Totals.Taxes)

		// Check the VAT tax category
		var vatCategory *tax.CategoryTotal
		for _, cat := range invoice.Totals.Taxes.Categories {
			if cat.Code == "VAT" {
				vatCategory = cat
				break
			}
		}

		require.NotNil(t, vatCategory, "VAT tax category should exist")

		// Check the VAT tax amount
		assert.True(t, vatCategory.Amount.Compare(num.MakeAmount(35640, 2)) == 0, "VAT amount should be 356.40")

		// Check the VAT tax rates
		require.NotEmpty(t, vatCategory.Rates)

		// Find the 22% VAT rate
		var vat22Rate *tax.RateTotal
		for _, rt := range vatCategory.Rates {
			if rt.Percent != nil && rt.Percent.Compare(num.MakePercentage(220, 3)) == 0 {
				vat22Rate = rt
				break
			}
		}

		require.NotNil(t, vat22Rate, "22% VAT rate should exist")
		assert.True(t, vat22Rate.Base.Compare(num.MakeAmount(162000, 2)) == 0, "VAT base should be 1620.00")
		assert.True(t, vat22Rate.Amount.Compare(num.MakeAmount(35640, 2)) == 0, "VAT amount should be 356.40")

		// Find the exempt rate
		var exemptRate *tax.RateTotal
		for _, rt := range vatCategory.Rates {
			if rt.Percent == nil || rt.Percent.IsZero() {
				exemptRate = rt
				break
			}
		}

		require.NotNil(t, exemptRate, "Exempt rate should exist")
		assert.NotEmpty(t, exemptRate.Base)
		assert.True(t, exemptRate.Amount.IsZero(), "Exempt rate amount should be zero")
	})

	t.Run("should convert retained taxes correctly", func(t *testing.T) {
		// Load the XML file with retained taxes
		data, err := os.ReadFile(filepath.Join(test.GetDataPath(test.PathFatturaPAGOBL), "invoice-irpef.xml"))
		require.NoError(t, err)

		// Convert XML to GOBL
		env, err := test.ConvertToGOBL(data)
		require.NoError(t, err)
		require.NotNil(t, env)

		// Extract the invoice
		invoice, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		require.NotNil(t, invoice)

		// Check the retained taxes
		require.NotNil(t, invoice.Totals)
		require.NotNil(t, invoice.Totals.Taxes)

		// Find the retained tax category
		var retainedCategory *tax.CategoryTotal
		for _, cat := range invoice.Totals.Taxes.Categories {
			if cat.Retained {
				retainedCategory = cat
				break
			}
		}

		require.NotNil(t, retainedCategory, "Retained tax category should exist")

		// Check the retained tax amount
		assert.NotEmpty(t, retainedCategory.Amount)
		assert.True(t, retainedCategory.Amount.Compare(num.MakeAmount(37400, 2)) == 0, "Retained tax amount should be 374.00")

		// Check the retained tax rates
		require.NotEmpty(t, retainedCategory.Rates)
		require.Len(t, retainedCategory.Rates, 2)

		// Check first rate (20% IRPEF)
		rate1 := retainedCategory.Rates[0]
		require.NotNil(t, rate1.Percent)
		assert.True(t, rate1.Percent.Compare(num.MakePercentage(200, 3)) == 0, "Rate should be 20.0%")
		assert.True(t, rate1.Amount.Compare(num.MakeAmount(32400, 2)) == 0, "Amount should be 324.00")
		assert.Equal(t, cbc.Code("A"), rate1.Ext[sdi.ExtKeyRetained])

		// Check second rate (50% IRPEF)
		rate2 := retainedCategory.Rates[1]
		require.NotNil(t, rate2.Percent)
		assert.True(t, rate2.Percent.Compare(num.MakePercentage(500, 3)) == 0, "Rate should be 50.0%")
		assert.True(t, rate2.Amount.Compare(num.MakeAmount(5000, 2)) == 0, "Amount should be 50.00")
		assert.Equal(t, cbc.Code("I"), rate2.Ext[sdi.ExtKeyRetained])
	})
}
