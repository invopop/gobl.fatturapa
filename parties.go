package fatturapa

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

const (
	statoLiquidazioneDefault   = "LN"
	euCitizenTaxCodeDefault    = "0000000"
	nonEUCitizenTaxCodeDefault = "99999999999"
)

type supplier struct {
	DatiAnagrafici *datiAnagrafici
	Sede           *address
	IscrizioneREA  *iscrizioneREA `xml:",omitempty"`
	Contatti       *contatti      `xml:",omitempty"`
}

type customer struct {
	DatiAnagrafici *datiAnagrafici
	Sede           *address
}

// datiAnagrafici contains information related to an individual or company
type datiAnagrafici struct {
	IdFiscaleIVA *taxID `xml:",omitempty"` // nolint:revive
	// CodiceFiscale is the Italian fiscal code, distinct from TaxID
	CodiceFiscale string `xml:",omitempty"`
	Anagrafica    *anagrafica
	// RegimeFiscale identifies the tax system to be applied
	// Has the form RFXX where XX is numeric; required only for the supplier
	RegimeFiscale string `xml:",omitempty"`
}

// anagrafica contains further party information
type anagrafica struct {
	// Name of the organization
	Denominazione string `xml:",omitempty"`
	// Name of the person
	Nome string `xml:",omitempty"`
	// Surname of the person
	Cognome string `xml:",omitempty"`
	// Title of the person
	Titolo string `xml:",omitempty"`
	// EORI (Economic Operator Registration and Identification) code
	CodEORI string `xml:",omitempty"`
}

// iscrizioneREA contains information related to the company registration details (REA)
type iscrizioneREA struct {
	// Initials of the province where the company's Registry Office is located
	Ufficio string
	// Company's REA registration number
	NumeroREA string
	// Company's share capital
	CapitaleSociale string `xml:",omitempty"`
	// Indication of whether the Company is in liquidation or not.
	// Possible values: LS (in liquidation), LN (not in liquidation)
	StatoLiquidazione string
}

type contatti struct {
	Telefono string `xml:",omitempty"`
	Email    string `xml:",omitempty"`
}

func newCedentePrestatore(s *org.Party) *supplier {
	ns := &supplier{
		DatiAnagrafici: &datiAnagrafici{
			IdFiscaleIVA: &taxID{
				IdPaese:  s.TaxID.Country.String(),
				IdCodice: s.TaxID.Code.String(),
			},
			Anagrafica: newAnagrafica(s),
		},
		IscrizioneREA: newIscrizioneREA(s),
		Contatti:      newContatti(s),
	}

	if v, ok := s.Ext[it.ExtKeySDIFiscalRegime]; ok {
		ns.DatiAnagrafici.RegimeFiscale = v.String()
	} else {
		ns.DatiAnagrafici.RegimeFiscale = "RF01"
	}

	if len(s.Addresses) > 0 {
		ns.Sede = newAddress(s.Addresses[0])
	}

	return ns
}

func newCessionarioCommittente(c *org.Party) *customer {
	if c == nil {
		return nil
	}

	nc := new(customer)
	if len(c.Addresses) > 0 {
		nc.Sede = newAddress(c.Addresses[0])
	}

	da := &datiAnagrafici{
		Anagrafica: newAnagrafica(c),
	}

	if c.TaxID != nil {
		if isCodiceFiscale(c.TaxID) {
			da.CodiceFiscale = c.TaxID.Code.String()
		} else if isEUCountry(c.TaxID.Country) {
			da.IdFiscaleIVA = customerFiscaleIVA(c.TaxID, euCitizenTaxCodeDefault)
		} else {
			da.IdFiscaleIVA = customerFiscaleIVA(c.TaxID, nonEUCitizenTaxCodeDefault)
		}
	}

	nc.DatiAnagrafici = da

	return nc
}

func newAnagrafica(party *org.Party) *anagrafica {
	if len(party.People) > 0 && party.TaxID.Type == it.TaxIdentityTypeIndividual {
		name := party.People[0].Name

		return &anagrafica{
			Nome:    name.Given,
			Cognome: name.Surname,
			Titolo:  name.Prefix,
		}
	}

	return &anagrafica{
		Denominazione: party.Name,
	}
}

func newContatti(party *org.Party) *contatti {
	c := &contatti{}

	if len(party.Emails) > 0 {
		c.Email = party.Emails[0].Address
	}

	if len(party.Telephones) > 0 {
		c.Telefono = party.Telephones[0].Number
	}

	return c
}

func customerFiscaleIVA(id *tax.Identity, fallBack string) *taxID {
	idCodice := id.Code.String()

	if idCodice == "" {
		idCodice = fallBack
	}

	return &taxID{
		IdPaese:  id.Country.String(),
		IdCodice: idCodice,
	}
}

func newIscrizioneREA(supplier *org.Party) *iscrizioneREA {
	if supplier.Registration == nil {
		return nil
	}

	capital := supplier.Registration.Capital
	var capitalFormatted string

	if capital == nil {
		capitalFormatted = ""
	} else {
		capitalFormatted = capital.Rescale(2).String()
	}

	return &iscrizioneREA{
		Ufficio:           supplier.Registration.Office,
		NumeroREA:         supplier.Registration.Entry,
		CapitaleSociale:   capitalFormatted,
		StatoLiquidazione: statoLiquidazioneDefault,
	}
}

func isCodiceFiscale(taxID *tax.Identity) bool {
	if taxID.Country != l10n.IT {
		return false
	}

	return len(taxID.Code.String()) == 16
}
