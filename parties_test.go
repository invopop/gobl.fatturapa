package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPartiesSupplier(t *testing.T) {
	t.Run("should contain the supplier party info", func(t *testing.T) {
		doc, err := test.LoadGOBL("invoice-vat.json")
		require.NoError(t, err)

		supplier := doc.FatturaElettronicaHeader.CedentePrestatore

		assert.Equal(t, "IT", supplier.DatiAnagrafici.IdFiscaleIVA.IdPaese)
		assert.Equal(t, "12345678903", supplier.DatiAnagrafici.IdFiscaleIVA.IdCodice)
		assert.Equal(t, "MªF. Services", supplier.DatiAnagrafici.Anagrafica.Denominazione)
		assert.Equal(t, "RF01", supplier.DatiAnagrafici.RegimeFiscale)
		assert.Equal(t, "VIALE DELLA LIBERTÀ, 1", supplier.Address.Indirizzo)
		assert.Equal(t, "00100", supplier.Address.CAP)
		assert.Equal(t, "ROMA", supplier.Address.Comune)
		assert.Equal(t, "RM", supplier.Address.Provincia)
		assert.Equal(t, "IT", supplier.Address.Nazione)
	})
}

func TestPartiesCustomer(t *testing.T) {
	t.Run("should contain the customer party info", func(t *testing.T) {
		doc, err := test.LoadGOBL("invoice-vat.json")
		require.NoError(t, err)

		customer := doc.FatturaElettronicaHeader.CessionarioCommittente

		assert.Equal(t, "RSSGNC73A02F205X", customer.DatiAnagrafici.CodiceFiscale)
		assert.Equal(t, "MARIO", customer.DatiAnagrafici.Anagrafica.Nome)
		assert.Equal(t, "LEONI", customer.DatiAnagrafici.Anagrafica.Cognome)
		assert.Equal(t, "Dott.", customer.DatiAnagrafici.Anagrafica.Titolo)
		assert.Equal(t, "VIALE DELI LAVORATORI, 32", customer.Address.Indirizzo)
		assert.Equal(t, "50100", customer.Address.CAP)
		assert.Equal(t, "FIRENZE", customer.Address.Comune)
		assert.Equal(t, "FI", customer.Address.Provincia)
		assert.Equal(t, "IT", customer.Address.Nazione)
	})
}
