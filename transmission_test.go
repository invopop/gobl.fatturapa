package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransmissionData(t *testing.T) {
	t.Run("should contain transmitting subject info", func(t *testing.T) {
		converter := test.TestConverter()
		doc, err := test.LoadGOBL("invoice-simple.json", converter)
		require.NoError(t, err)

		dt := doc.FatturaElettronicaHeader.DatiTrasmissione

		assert.Equal(t, converter.Config.Transmitter.CountryCode, dt.IdTrasmittente.IdPaese)
		assert.Equal(t, converter.Config.Transmitter.TaxID, dt.IdTrasmittente.IdCodice)
		assert.Equal(t, "679a2f25", dt.ProgressivoInvio)
		assert.Equal(t, "FPR12", dt.FormatoTrasmissione)
		assert.Equal(t, "ABCDEF1", dt.CodiceDestinatario)
	})

	t.Run("should skip transmission info if transmitter is not present", func(t *testing.T) {
		converter := test.TestConverter()

		converter.Config.Transmitter = nil

		doc, err := test.LoadGOBL("invoice-simple.json", converter)
		require.NoError(t, err)

		assert.Nil(t, doc.FatturaElettronicaHeader.DatiTrasmissione)
	})
}
