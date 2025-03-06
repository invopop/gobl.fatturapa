package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDettaglioLinee(t *testing.T) {
	t.Run("should contain the line info", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
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

func TestDatiRiepilogo(t *testing.T) {
	t.Run("should contain the tax summary info", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
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
