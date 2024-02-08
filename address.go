package fatturapa

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

const (
	foreignCAP = "00000"
)

// address from IndirizzoType
type address struct {
	Indirizzo    string // Street
	NumeroCivico string `xml:",omitempty"` // Number
	CAP          string // Post Code
	Comune       string // Locality
	Provincia    string `xml:",omitempty"` // Region
	Nazione      string // Country Code
}

func newAddress(addr *org.Address) *address {
	ad := &address{
		Indirizzo:    addressStreet(addr),
		NumeroCivico: addr.Number,
		Comune:       addr.Locality,
		Provincia:    addr.Region,
		Nazione:      addr.Country.String(),
	}
	if addr.Country == l10n.IT {
		ad.CAP = addr.Code
	} else {
		ad.CAP = foreignCAP
	}
	return ad
}

func addressStreet(address *org.Address) string {
	if address.PostOfficeBox != "" {
		return address.PostOfficeBox
	}
	return address.Street
}
