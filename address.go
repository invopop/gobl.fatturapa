package fatturapa

import (
	"errors"

	"github.com/invopop/gobl/org"
)

// Address contains the address of the party
type Address struct {
	Indirizzo    string
	NumeroCivico string `xml:",omitempty"`
	CAP          string
	Comune       string
	Provincia    string `xml:",omitempty"`
	Nazione      string
}

func newAddress(p *org.Party) (*Address, error) {
	if len(p.Addresses) == 0 {
		return nil, errors.New("party missing address")
	}

	address := p.Addresses[0]

	return &Address{
		Indirizzo: addressLine(address),
		CAP:       address.Code,
		Comune:    address.Locality,
		Provincia: address.Region,
		Nazione:   address.Country.String(),
	}, nil
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
