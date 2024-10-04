package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentsSimple(t *testing.T) {
	t.Run("should contain the supplier party info", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dp := doc.FatturaElettronicaBody[0].DatiPagamento

		require.NotNil(t, dp)
		assert.Equal(t, "TP02", dp.CondizioniPagamento)
		assert.Len(t, dp.DettaglioPagamento, 1)
		assert.Equal(t, "MP08", dp.DettaglioPagamento[0].ModalitaPagamento)
		assert.Equal(t, "1388.40", dp.DettaglioPagamento[0].ImportoPagamento)
	})
}

func TestPayments(t *testing.T) {
	t.Run("multiple due dates", func(t *testing.T) {
		env := test.LoadTestFile("invoice-irpef.json")
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dp := doc.FatturaElettronicaBody[0].DatiPagamento

		require.NotNil(t, dp)
		assert.Equal(t, "TP01", dp.CondizioniPagamento)
		assert.Len(t, dp.DettaglioPagamento, 2)
		assert.Equal(t, "MP05", dp.DettaglioPagamento[0].ModalitaPagamento)
		assert.Equal(t, "2023-03-02", dp.DettaglioPagamento[0].DueDate)
		assert.Equal(t, "500.00", dp.DettaglioPagamento[0].ImportoPagamento)
		assert.Equal(t, "MP05", dp.DettaglioPagamento[1].ModalitaPagamento)
		assert.Equal(t, "2023-04-02", dp.DettaglioPagamento[1].DueDate)
		assert.Equal(t, "544.40", dp.DettaglioPagamento[1].ImportoPagamento)
	})

	t.Run("advance payment", func(t *testing.T) {
		env := test.LoadTestFile("invoice-hotel-private.json")
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dp := doc.FatturaElettronicaBody[0].DatiPagamento

		require.NotNil(t, dp)
		assert.Equal(t, "TP01", dp.CondizioniPagamento)
		assert.Len(t, dp.DettaglioPagamento, 1)
		assert.Equal(t, "MP08", dp.DettaglioPagamento[0].ModalitaPagamento)
		assert.Equal(t, "29.00", dp.DettaglioPagamento[0].ImportoPagamento)
	})

	t.Run("prepaid", func(t *testing.T) {
		env := test.LoadTestFile("invoice-hotel.json")
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dp := doc.FatturaElettronicaBody[0].DatiPagamento

		require.NotNil(t, dp)
		assert.Equal(t, "TP03", dp.CondizioniPagamento)
		assert.Len(t, dp.DettaglioPagamento, 1)
		assert.Equal(t, "MP08", dp.DettaglioPagamento[0].ModalitaPagamento)
		assert.Equal(t, "241.00", dp.DettaglioPagamento[0].ImportoPagamento)
	})
}
