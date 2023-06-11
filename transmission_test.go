package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransmissionData(t *testing.T) {
	t.Run("should contain transmitting subject info", func(t *testing.T) {
		converter := test.NewConverter()
		env := test.LoadTestFile("invoice-simple.json")
		doc, err := test.ConvertFromGOBL(env, converter)
		require.NoError(t, err)

		dt := doc.FatturaElettronicaHeader.DatiTrasmissione

		assert.Equal(t, converter.Config.Transmitter.CountryCode, dt.IdTrasmittente.IdPaese)
		assert.Equal(t, converter.Config.Transmitter.TaxID, dt.IdTrasmittente.IdCodice)
		assert.Equal(t, "679a2f25", dt.ProgressivoInvio)
		assert.Equal(t, "FPR12", dt.FormatoTrasmissione)
		assert.Equal(t, "ABCDEF1", dt.CodiceDestinatario)
	})

	t.Run("should skip transmitter info and only include codice destinatario if transmitter is not present", func(t *testing.T) {
		converter := test.NewConverter()
		converter.Config.Transmitter = nil

		env := test.LoadTestFile("invoice-simple.json")
		doc, err := test.ConvertFromGOBL(env, converter)
		require.NoError(t, err)

		dt := doc.FatturaElettronicaHeader.DatiTrasmissione

		assert.Equal(t, "ABCDEF1", dt.CodiceDestinatario)
		assert.Equal(t, nil, dt.IdTrasmittente)
		assert.Equal(t, "", dt.ProgressivoInvio)
		assert.Equal(t, "", dt.FormatoTrasmissione)
	})
}
