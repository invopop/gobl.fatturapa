package fatturapa_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/it"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPartiesSupplier(t *testing.T) {
	t.Run("should contain the supplier info", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		s := doc.FatturaElettronicaHeader.Supplier

		assert.Equal(t, "IT", s.Identity.TaxID.Country)
		assert.Equal(t, "12345678903", s.Identity.TaxID.Code)
		assert.Equal(t, "MªF. Services", s.Identity.Profile.Name)
		assert.Equal(t, "RF02", s.Identity.FiscalRegime)
		assert.Equal(t, "VIALE DELLA LIBERTÀ", s.Address.Street)
		assert.Equal(t, "1", s.Address.Number)
		assert.Equal(t, "00100", s.Address.Code)
		assert.Equal(t, "ROMA", s.Address.Locality)
		assert.Equal(t, "RM", s.Address.Region)
		assert.Equal(t, "IT", s.Address.Country)
		assert.Equal(t, "RM", s.Registration.Office)
		assert.Equal(t, "123456", s.Registration.Entry)
		assert.Equal(t, "50000.00", s.Registration.Capital)
		assert.Equal(t, "LN", s.Registration.LiquidationState)
	})

	t.Run("should set the supplier fiscal regime", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		inv := env.Extract().(*bill.Invoice)
		inv.Supplier.Ext = nil
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		s := doc.FatturaElettronicaHeader.Supplier

		assert.Equal(t, "RF01", s.Identity.FiscalRegime)
	})
}

func TestPartiesCustomer(t *testing.T) {
	t.Run("should contain the customer info", func(t *testing.T) {
		env := test.LoadTestFile("invoice-irpef.json")
		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		c := doc.FatturaElettronicaHeader.Customer

		assert.Nil(t, c.Identity.TaxID)
		assert.Equal(t, "MRALNE80E05H501C", c.Identity.FiscalCode)
		assert.Equal(t, "", c.Identity.Profile.Name)
		assert.Equal(t, "MARIO", c.Identity.Profile.Given)
		assert.Equal(t, "LEONI", c.Identity.Profile.Surname)
		assert.Equal(t, "Dott.", c.Identity.Profile.Title)
		assert.Equal(t, "VIALE DELI LAVORATORI", c.Address.Street)
		assert.Equal(t, "32", c.Address.Number)
		assert.Equal(t, "50100", c.Address.Code)
		assert.Equal(t, "FIRENZE", c.Address.Locality)
		assert.Equal(t, "FI", c.Address.Region)
		assert.Equal(t, "IT", c.Address.Country)
	})

	t.Run("should contain customer info with codice fiscale", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Customer.TaxID.Code = ""
			inv.Customer.Identities = org.AddIdentity(inv.Customer.Identities,
				&org.Identity{
					Key:  it.IdentityKeyFiscalCode,
					Code: "RSSGNC73A02F205X",
				},
			)
		})

		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		c := doc.FatturaElettronicaHeader.Customer

		assert.Nil(t, c.Identity.TaxID)
		assert.Equal(t, "RSSGNC73A02F205X", c.Identity.FiscalCode)
	})

	t.Run("should contain customer info for EU citizen with Tax ID given", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Customer.TaxID.Code = "81237984062783472"
			inv.Customer.TaxID.Country = l10n.AT
		})

		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		c := doc.FatturaElettronicaHeader.Customer

		assert.Equal(t, "AT", c.Identity.TaxID.Country)
		assert.Equal(t, "81237984062783472", c.Identity.TaxID.Code)
	})

	t.Run("should contain customer info for EU citizen with no Tax ID given", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Customer.TaxID.Code = ""
			inv.Customer.TaxID.Country = l10n.SE
		})

		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		c := doc.FatturaElettronicaHeader.Customer

		assert.Equal(t, "SE", c.Identity.TaxID.Country)
		assert.Equal(t, "0000000", c.Identity.TaxID.Code)
	})

	t.Run("should replace customer ID info for non-EU citizen with Tax ID given", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Customer.TaxID.Code = "09823876432"
			inv.Customer.TaxID.Country = l10n.GB
		})

		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		c := doc.FatturaElettronicaHeader.Customer

		assert.Equal(t, "GB", c.Identity.TaxID.Country)
		assert.Equal(t, "OO99999999999", c.Identity.TaxID.Code)
	})

	t.Run("should contain customer info for non-EU citizen with no Tax ID given", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Customer.TaxID.Code = ""
			inv.Customer.TaxID.Country = l10n.JP
		})

		doc, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)

		c := doc.FatturaElettronicaHeader.Customer

		assert.Equal(t, "JP", c.Identity.TaxID.Country)
		assert.Equal(t, "0000000", c.Identity.TaxID.Code)
	})

	t.Run("should not fail if missing key data", func(t *testing.T) {
		env := test.LoadTestFile("invoice-simple.json")
		test.ModifyInvoice(env, func(inv *bill.Invoice) {
			inv.Customer.TaxID.Country = ""
		})

		_, err := test.ConvertFromGOBL(env)
		require.NoError(t, err)
	})
}
