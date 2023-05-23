package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPartiesSupplier(t *testing.T) {
	t.Run("should contain the supplier info", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		s := doc.FatturaElettronicaHeader.CedentePrestatore

		assert.Equal(t, "IT", s.DatiAnagrafici.IdFiscaleIVA.IdPaese)
		assert.Equal(t, "12345678903", s.DatiAnagrafici.IdFiscaleIVA.IdCodice)
		assert.Equal(t, "MªF. Services", s.DatiAnagrafici.Anagrafica.Denominazione)
		assert.Equal(t, "RF01", s.DatiAnagrafici.RegimeFiscale)
		assert.Equal(t, "VIALE DELLA LIBERTÀ, 1", s.Sede.Indirizzo)
		assert.Equal(t, "00100", s.Sede.CAP)
		assert.Equal(t, "ROMA", s.Sede.Comune)
		assert.Equal(t, "RM", s.Sede.Provincia)
		assert.Equal(t, "IT", s.Sede.Nazione)
		assert.Equal(t, "RM", s.IscrizioneREA.Ufficio)
		assert.Equal(t, "123456", s.IscrizioneREA.NumeroREA)
		assert.Equal(t, "50000.00", s.IscrizioneREA.CapitaleSociale)
		assert.Equal(t, "LN", s.IscrizioneREA.StatoLiquidazione)
	})
}

func TestPartiesCustomer(t *testing.T) {
	t.Run("should contain the customer info", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		c := doc.FatturaElettronicaHeader.CessionarioCommittente

		assert.Equal(t, "IT", c.DatiAnagrafici.IdFiscaleIVA.IdPaese)
		assert.Equal(t, "09876543210", c.DatiAnagrafici.IdFiscaleIVA.IdCodice)
		assert.Equal(t, "", c.DatiAnagrafici.CodiceFiscale)
		assert.Equal(t, "MARIO", c.DatiAnagrafici.Anagrafica.Nome)
		assert.Equal(t, "LEONI", c.DatiAnagrafici.Anagrafica.Cognome)
		assert.Equal(t, "Dott.", c.DatiAnagrafici.Anagrafica.Titolo)
		assert.Equal(t, "VIALE DELI LAVORATORI, 32", c.Sede.Indirizzo)
		assert.Equal(t, "50100", c.Sede.CAP)
		assert.Equal(t, "FIRENZE", c.Sede.Comune)
		assert.Equal(t, "FI", c.Sede.Provincia)
		assert.Equal(t, "IT", c.Sede.Nazione)
	})

	t.Run("should contain customer info with codice fiscale", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Customer.TaxID.Code = "RSSGNC73A02F205X"
		})

		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		c := doc.FatturaElettronicaHeader.CessionarioCommittente

		assert.Nil(t, c.DatiAnagrafici.IdFiscaleIVA)
		assert.Equal(t, "RSSGNC73A02F205X", c.DatiAnagrafici.CodiceFiscale)
	})

	t.Run("should contain customer info for EU citizen with Tax ID given", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Customer.TaxID.Code = "81237984062783472"
			inv.Customer.TaxID.Country = l10n.AT
		})

		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		c := doc.FatturaElettronicaHeader.CessionarioCommittente

		assert.Equal(t, "AT", c.DatiAnagrafici.IdFiscaleIVA.IdPaese)
		assert.Equal(t, "81237984062783472", c.DatiAnagrafici.IdFiscaleIVA.IdCodice)
	})

	t.Run("should contain customer info for EU citizen with no Tax ID given", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Customer.TaxID.Code = ""
			inv.Customer.TaxID.Country = l10n.SE
		})

		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		c := doc.FatturaElettronicaHeader.CessionarioCommittente

		assert.Equal(t, "SE", c.DatiAnagrafici.IdFiscaleIVA.IdPaese)
		assert.Equal(t, "0000000", c.DatiAnagrafici.IdFiscaleIVA.IdCodice)
	})
}

// [x] test codice fiscale: RSSGNC73A02F205X
// [ ] test italian missing tax id
// [ ] test EU citizen with tax id
// [ ] test EU citizen without tax id
// [ ] test non-eu citizen with tax id
// [ ] test non-eu citizen without tax id
