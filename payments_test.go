package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentsSimple(t *testing.T) {
	t.Run("should contain the supplier party info", func(t *testing.T) {
		doc, err := test.LoadGOBL("invoice-simple.json", test.TestConverter())
		require.NoError(t, err)

		dp := doc.FatturaElettronicaBody[0].DatiPagamento

		require.NotNil(t, dp)
		assert.Equal(t, "TP02", dp.CondizioniPagamento)
		assert.Len(t, dp.DettaglioPagamento, 1)
		assert.Equal(t, "MP08", dp.DettaglioPagamento[0].ModalitaPagamento)
		assert.Equal(t, "1388.40", dp.DettaglioPagamento[0].ImportoPagamento)
	})
}

func TestPaymentsMultipleDueDates(t *testing.T) {
	t.Run("should contain the customer party info", func(t *testing.T) {
		doc, err := test.LoadGOBL("invoice-irpef.json", test.TestConverter())
		require.NoError(t, err)

		dp := doc.FatturaElettronicaBody[0].DatiPagamento

		require.NotNil(t, dp)
		assert.Equal(t, "TP01", dp.CondizioniPagamento)
		assert.Len(t, dp.DettaglioPagamento, 2)
		assert.Equal(t, "MP05", dp.DettaglioPagamento[0].ModalitaPagamento)
		assert.Equal(t, "2023-03-02", dp.DettaglioPagamento[0].DataScadenzaPagamento)
		assert.Equal(t, "500.00", dp.DettaglioPagamento[0].ImportoPagamento)
		assert.Equal(t, "MP05", dp.DettaglioPagamento[1].ModalitaPagamento)
		assert.Equal(t, "2023-04-02", dp.DettaglioPagamento[1].DataScadenzaPagamento)
		assert.Equal(t, "544.40", dp.DettaglioPagamento[1].ImportoPagamento)
	})
}
