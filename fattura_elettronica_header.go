package fatturapa

import (
	"errors"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
)

const (
	FormatoTrasmissione  = "FPA12"
	RegimeFiscaleDefault = "RF01"
)

type FatturaElettronicaHeader struct {
	DatiTrasmissione       DatiTrasmissione       `xml:",omitempty"`
	CedentePrestatore      CedentePrestatore      `xml:",omitempty"`
	CessionarioCommittente CessionarioCommittente `xml:",omitempty"`
}

// Data related to the transmitting subject
type DatiTrasmissione struct {
	IdTrasmittente      VATID
	ProgressivoInvio    string
	FormatoTrasmissione string
	CodiceDestinatario  string
}

// Data related to the supplier
type CedentePrestatore struct {
	DatiAnagrafici DatiAnagrafici
	Sede           Sede
}

// Data related to the customer
type CessionarioCommittente struct {
	DatiAnagrafici DatiAnagrafici
	Sede           Sede
}

type DatiAnagrafici struct {
	IdFiscaleIVA VATID `xml:",omitempty"`
	// CodiceFiscale is the Italian fiscal code, distinct from VATID
	CodiceFiscale string `xml:",omitempty"`
	Anagrafica    Anagrafica
	// RegimeFiscale identifies the tax system to be applied
	// Has the form RFXX where XX is numeric; required only for the supplier
	RegimeFiscale string
}

// Anagrafica contains information related to an individual or company
type Anagrafica struct {
	// Name of the organization
	Denominazione string
}

// VATID is the VAT identification number consisting of a country code and the
// actual VAT number.
type VATID struct {
	// ISO 3166-1 alpha-2 country code
	IdPaese string
	// Actual VAT number
	IdCodice string
}

// Sede contains the address of the party
type Sede struct {
	Indirizzo string
	CAP       string
	Comune    string
	Provincia string
	Nazione   string
}

func newFatturaElettronicaHeader(inv bill.Invoice) (*FatturaElettronicaHeader, error) {
	sedeSupplier, err := newSede(inv.Customer)
	if err != nil {
		return nil, err
	}

	sedeCustomer, err := newSede(inv.Customer)
	if err != nil {
		return nil, err
	}

	dataAnagraficiCustomer, err := newCustomerDataAnagrafici(inv.Customer)
	if err != nil {
		return nil, err
	}

	return &FatturaElettronicaHeader{
		DatiTrasmissione: DatiTrasmissione{
			IdTrasmittente: VATID{
				IdPaese:  inv.Supplier.TaxID.Country.String(),
				IdCodice: inv.Supplier.TaxID.Code.String(),
			},
			ProgressivoInvio:    inv.Code,
			FormatoTrasmissione: FormatoTrasmissione,
			CodiceDestinatario:  inv.Meta["fatturapa-codice-destinatario"],
		},
		CedentePrestatore: CedentePrestatore{
			DatiAnagrafici: DatiAnagrafici{
				IdFiscaleIVA: VATID{
					IdPaese:  inv.Supplier.TaxID.Country.String(),
					IdCodice: inv.Supplier.TaxID.Code.String(),
				},
				Anagrafica: Anagrafica{
					Denominazione: inv.Supplier.Name,
				},
				RegimeFiscale: RegimeFiscaleDefault,
			},
			Sede: *sedeSupplier,
		},
		CessionarioCommittente: CessionarioCommittente{
			DatiAnagrafici: *dataAnagraficiCustomer,
			Sede:           *sedeCustomer,
		},
	}, nil
}

func newSede(p *org.Party) (*Sede, error) {
	if len(p.Addresses) == 0 {
		return nil, errors.New("party missing address")
	}

	address := p.Addresses[0]

	return &Sede{
		Indirizzo: addressLine(address),
		CAP:       address.Code,
		Comune:    address.Locality,
		Provincia: address.Region,
		Nazione:   address.Country.String(),
	}, nil
}

func newCustomerDataAnagrafici(c *org.Party) (*DatiAnagrafici, error) {
	da := &DatiAnagrafici{
		Anagrafica: Anagrafica{
			Denominazione: c.Name,
		},
	}

	// Apply VATID or fiscal code. At least one of them is required.
	// FatturaPA only evaluates VATID if both are present
	if c.TaxID != nil {
		da.IdFiscaleIVA = VATID{
			IdPaese:  c.TaxID.Country.String(),
			IdCodice: c.TaxID.Code.String(),
		}
	} else {
		for _, id := range c.Identities {
			if id.Type == "CF" {
				da.CodiceFiscale = id.Code.String()
			}
		}

		if da.CodiceFiscale == "" {
			return nil, errors.New("customer has no VATID or fiscal code")
		}
	}

	return da, nil
}

func addressLine(address *org.Address) string {
	if address.PostOfficeBox != "" {
		return address.PostOfficeBox
	}

	return address.Street +
		", " + address.Number +
		addressMaybe(address.Block) +
		addressMaybe(address.Floor) +
		addressMaybe(address.Door)
}

func addressMaybe(element string) string {
	if element != "" {
		return ", " + element
	}
	return ""
}
