package fatturapa

import (
	"github.com/invopop/gobl/org"
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
	return &address{
		Indirizzo:    addressStreet(addr),
		NumeroCivico: addr.Number,
		CAP:          addr.Code,
		Comune:       addr.Locality,
		Provincia:    addr.Region,
		Nazione:      addr.Country.String(),
	}
}

func addressStreet(address *org.Address) string {
	if address.PostOfficeBox != "" {
		return address.PostOfficeBox
	}
	return address.Street
}
