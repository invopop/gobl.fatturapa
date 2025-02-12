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

		dl := doc.FatturaElettronicaBody[0].DatiBeniServizi.DettaglioLinee[0]

		assert.Equal(t, "1", dl.NumeroLinea)
		assert.Equal(t, "Development services", dl.Descrizione)
		assert.Equal(t, "20.00", dl.Quantita)
		assert.Equal(t, "90.00", dl.PrezzoUnitario)
		assert.Equal(t, "1620.00", dl.PrezzoTotale)
		assert.Equal(t, "22.00", dl.AliquotaIVA)
		assert.Equal(t, "", dl.Natura)

		sm := dl.ScontoMaggiorazione[0]

		assert.Equal(t, "SC", sm.Tipo)
		assert.Equal(t, "10.00", sm.Percentuale)
		assert.Equal(t, "9.0000", sm.Importo)

		dl = doc.FatturaElettronicaBody[0].DatiBeniServizi.DettaglioLinee[1]

		assert.Equal(t, "2", dl.NumeroLinea)
		assert.Equal(t, "N2.2", dl.Natura)
	})
}

func TestDatiRiepilogo(t *testing.T) {
	t.Run("should contain the tax summary info", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		dr := doc.FatturaElettronicaBody[0].DatiBeniServizi.DatiRiepilogo[0]

		assert.Equal(t, "22.00", dr.AliquotaIVA)
		assert.Equal(t, "1620.00", dr.ImponibileImporto)
		assert.Equal(t, "356.40", dr.Imposta)
		assert.Equal(t, "", dr.Natura)
		assert.Equal(t, "", dr.RiferimentoNormativo)

		dr = doc.FatturaElettronicaBody[0].DatiBeniServizi.DatiRiepilogo[1]

		assert.Equal(t, "N2.2", dr.Natura)
		assert.Equal(t, "S", dr.EsigibilitaIVA)
		assert.Equal(t, "Non soggette - altri casi", dr.RiferimentoNormativo)
	})
}
