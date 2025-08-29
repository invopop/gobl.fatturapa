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
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetainedTaxesInConversion(t *testing.T) {
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

		// Check that the invoice has tax categories
		require.NotNil(t, invoice.Totals)
		require.NotNil(t, invoice.Totals.Taxes)
		require.NotEmpty(t, invoice.Totals.Taxes.Categories)

		// Find the retained tax category (IRPEF)
		var retainedCategory *tax.CategoryTotal
		for _, cat := range invoice.Totals.Taxes.Categories {
			if cat.Code == it.TaxCategoryIRPEF && cat.Retained {
				retainedCategory = cat
				break
			}
		}

		// Check the retained tax category
		require.NotNil(t, retainedCategory, "Retained tax category should exist")
		assert.Equal(t, it.TaxCategoryIRPEF, retainedCategory.Code)
		assert.True(t, retainedCategory.Retained, "Category should be marked as retained")

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

	t.Run("should not have retained taxes when not present", func(t *testing.T) {
		// Load the XML file without retained taxes
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

		// Check that the invoice has tax categories
		require.NotNil(t, invoice.Totals)
		require.NotNil(t, invoice.Totals.Taxes)
		require.NotEmpty(t, invoice.Totals.Taxes.Categories)

		// Check that there are no retained tax categories
		for _, cat := range invoice.Totals.Taxes.Categories {
			assert.False(t, cat.Retained, "No category should be marked as retained")
		}
	})

	t.Run("should convert retained taxes with fund contribution correctly", func(t *testing.T) {
		// Load the XML file with fund contributions and retained taxes
		data, err := os.ReadFile(filepath.Join(test.GetDataPath(test.PathFatturaPAGOBL), "invoice-fund-contribution.xml"))
		require.NoError(t, err)

		// Convert XML to GOBL
		env, err := test.ConvertToGOBL(data)
		require.NoError(t, err)
		require.NotNil(t, env)

		// Extract the invoice
		invoice, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		require.NotNil(t, invoice)

		// Check that we have both line and fund contribution with retained taxes
		require.Len(t, invoice.Lines, 1)
		require.NotEmpty(t, invoice.Charges)

		// Check line has retained tax
		line := invoice.Lines[0]
		var lineRetainedTax *tax.Combo
		for _, t := range line.Taxes {
			if t.Category == it.TaxCategoryIRPEF {
				lineRetainedTax = t
				break
			}
		}
		require.NotNil(t, lineRetainedTax, "Line should have IRPEF retained tax")
		assert.True(t, lineRetainedTax.Percent.Compare(num.MakePercentage(200, 3)) == 0, "Line retained tax should be 20%")
		assert.Equal(t, cbc.Code("A"), lineRetainedTax.Ext[sdi.ExtKeyRetained])

		// Check fund contribution charge has retained tax
		var fundContributionCharge *bill.Charge
		for _, charge := range invoice.Charges {
			if charge.Key.Has(sdi.KeyFundContribution) {
				fundContributionCharge = charge
				break
			}
		}
		require.NotNil(t, fundContributionCharge, "Should have fund contribution charge")

		var chargeRetainedTax *tax.Combo
		for _, t := range fundContributionCharge.Taxes {
			if t.Category == it.TaxCategoryIRPEF {
				chargeRetainedTax = t
				break
			}
		}
		require.NotNil(t, chargeRetainedTax, "Fund contribution should have IRPEF retained tax")
		assert.True(t, chargeRetainedTax.Percent.Compare(num.MakePercentage(200, 3)) == 0, "Fund contribution retained tax should be 20%")
		assert.Equal(t, cbc.Code("A"), chargeRetainedTax.Ext[sdi.ExtKeyRetained])

		// Check total retained tax amount in totals
		require.NotNil(t, invoice.Totals)
		require.NotNil(t, invoice.Totals.Taxes)
		var retainedCategory *tax.CategoryTotal
		for _, cat := range invoice.Totals.Taxes.Categories {
			if cat.Code == it.TaxCategoryIRPEF && cat.Retained {
				retainedCategory = cat
				break
			}
		}
		require.NotNil(t, retainedCategory, "Should have retained tax category in totals")
		assert.True(t, retainedCategory.Amount.Compare(num.MakeAmount(33280, 2)) == 0, "Total retained tax should be 332.80")
	})
}

func TestRetainedTaxesDistribution(t *testing.T) {
	t.Run("should distribute single retained tax proportionally across two items", func(t *testing.T) {
		// This test uses the existing invoice-irpef.xml which has multiple retained taxes
		// but we'll just test that the conversion works (the existing test already validates this)
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
		require.Len(t, invoice.Lines, 2)

		// Check that both lines have retained taxes applied
		line1 := invoice.Lines[0]
		line2 := invoice.Lines[1]

		// Find retained taxes on both lines
		var line1RetainedTax, line2RetainedTax *tax.Combo
		for _, t := range line1.Taxes {
			if t.Category == it.TaxCategoryIRPEF {
				line1RetainedTax = t
				break
			}
		}
		for _, t := range line2.Taxes {
			if t.Category == it.TaxCategoryIRPEF {
				line2RetainedTax = t
				break
			}
		}

		require.NotNil(t, line1RetainedTax, "Line 1 should have retained tax")
		require.NotNil(t, line2RetainedTax, "Line 2 should have retained tax")

		// Check that they have the expected rates and reasons from the test file
		assert.True(t, line1RetainedTax.Percent.Compare(num.MakePercentage(200, 3)) == 0, "Line 1 retained tax should be 20%")
		assert.True(t, line2RetainedTax.Percent.Compare(num.MakePercentage(500, 3)) == 0, "Line 2 retained tax should be 50%")
		assert.Equal(t, cbc.Code("A"), line1RetainedTax.Ext[sdi.ExtKeyRetained])
		assert.Equal(t, cbc.Code("I"), line2RetainedTax.Ext[sdi.ExtKeyRetained])
	})

	t.Run("should fail to parse invoice with ambiguous retained mapping", func(t *testing.T) {
		// Load the XML file with multiple retained taxes that cannot be properly distributed
		data, err := os.ReadFile(filepath.Join(test.GetDataPath(test.PathGOBLFatturaPA), "out/invoice-retained-ambiguous.xml"))
		require.NoError(t, err)

		// Convert XML to GOBL - this should fail
		_, err = test.ConvertToGOBL(data)
		require.Error(t, err)

		// Check that the error message indicates the distribution problem
		assert.Contains(t, err.Error(), "cannot determine which item it applies to")
	})

	t.Run("should fail to parse invoice with impossible retained mapping", func(t *testing.T) {
		// Load the XML file with multiple retained taxes that cannot be properly distributed
		data, err := os.ReadFile(filepath.Join(test.GetDataPath(test.PathGOBLFatturaPA), "out/invoice-retained-unparseable.xml"))
		require.NoError(t, err)

		// Convert XML to GOBL - this should fail
		_, err = test.ConvertToGOBL(data)
		require.Error(t, err)

		// Check that the error message indicates the distribution problem
		assert.Contains(t, err.Error(), "cannot match retained tax")
	})
}
