package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPartiesSupplier(t *testing.T) {
	t.Run("should contain the supplier party info", func(t *testing.T) {
		doc, err := test.LoadGOBL("invoice-vat.json", test.Client)
		require.NoError(t, err)

		s := doc.FatturaElettronicaHeader.CedentePrestatore

		assert.Equal(t, "IT", s.DatiAnagrafici.IdFiscaleIVA.IdPaese)
		assert.Equal(t, "12345678903", s.DatiAnagrafici.IdFiscaleIVA.IdCodice)
		assert.Equal(t, "MªF. Services", s.DatiAnagrafici.Anagrafica.Denominazione)
		assert.Equal(t, "RF01", s.DatiAnagrafici.RegimeFiscale)
		assert.Equal(t, "VIALE DELLA LIBERTÀ, 1", s.Address.Indirizzo)
		assert.Equal(t, "00100", s.Address.CAP)
		assert.Equal(t, "ROMA", s.Address.Comune)
		assert.Equal(t, "RM", s.Address.Provincia)
		assert.Equal(t, "IT", s.Address.Nazione)
	})
}

func TestPartiesCustomer(t *testing.T) {
	t.Run("should contain the customer party info", func(t *testing.T) {
		doc, err := test.LoadGOBL("invoice-vat.json", test.Client)
		require.NoError(t, err)

		c := doc.FatturaElettronicaHeader.CessionarioCommittente

		assert.Equal(t, "IT", c.DatiAnagrafici.IdFiscaleIVA.IdPaese)
		assert.Equal(t, "12345678903", c.DatiAnagrafici.IdFiscaleIVA.IdCodice)
		assert.Equal(t, "MARIO", c.DatiAnagrafici.Anagrafica.Nome)
		assert.Equal(t, "LEONI", c.DatiAnagrafici.Anagrafica.Cognome)
		assert.Equal(t, "Dott.", c.DatiAnagrafici.Anagrafica.Titolo)
		assert.Equal(t, "VIALE DELI LAVORATORI, 32", c.Address.Indirizzo)
		assert.Equal(t, "50100", c.Address.CAP)
		assert.Equal(t, "FIRENZE", c.Address.Comune)
		assert.Equal(t, "FI", c.Address.Provincia)
		assert.Equal(t, "IT", c.Address.Nazione)
	})
}
