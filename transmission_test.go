package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransmissionData(t *testing.T) {
	t.Run("should contain transmitting subject info", func(t *testing.T) {
		doc, err := test.LoadGOBL("invoice-vat.json", test.Client)
		require.NoError(t, err)

		dt := doc.FatturaElettronicaHeader.DatiTrasmissione

		assert.Equal(t, test.Client.CountryCode, dt.IdTrasmittente.IdPaese)
		assert.Equal(t, test.Client.TaxID, dt.IdTrasmittente.IdCodice)
		assert.Equal(t, "679a2f25", dt.ProgressivoInvio)
		assert.Equal(t, "FPR12", dt.FormatoTrasmissione)
		assert.Equal(t, "ABCDEF1", dt.CodiceDestinatario)
	})
}
