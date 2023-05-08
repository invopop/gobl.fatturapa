package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatiRitenuta(t *testing.T) {
	t.Run("When retained taxes are NOT present", func(t *testing.T) {
		t.Run("should be empty", func(t *testing.T) {
			doc, err := test.LoadGOBL("invoice-simple.json", test.TestConverter())
			require.NoError(t, err)

			dr := doc.FatturaElettronicaBody[0].DatiGenerali.DatiGeneraliDocumento.DatiRitenuta

			assert.Empty(t, dr)
		})
	})

	t.Run("When retained taxes are present", func(t *testing.T) {
		t.Run("should contain the correct retainted taxes", func(t *testing.T) {
			doc, err := test.LoadGOBL("invoice-irpef.json", test.TestConverter())
			require.NoError(t, err)

			dr := doc.FatturaElettronicaBody[0].DatiGenerali.DatiGeneraliDocumento.DatiRitenuta

			require.Len(t, dr, 2)

			assert.Equal(t, "RT01", dr[0].TipoRitenuta)
			assert.Equal(t, "324.00", dr[0].ImportoRitenuta)
			assert.Equal(t, "20.0%", dr[0].AliquotaRitenuta)
			assert.Equal(t, "A", dr[0].CausalePagamento)

			assert.Equal(t, "RT01", dr[1].TipoRitenuta)
			assert.Equal(t, "20.00", dr[1].ImportoRitenuta)
			assert.Equal(t, "20.0%", dr[1].AliquotaRitenuta)
			assert.Equal(t, "J", dr[1].CausalePagamento)
		})
	})
}
