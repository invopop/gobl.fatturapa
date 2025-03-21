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
		env, err := test.ConvertToGOBL(data, test.NewConverter())
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
		env, err := test.ConvertToGOBL(data, test.NewConverter())
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
}
