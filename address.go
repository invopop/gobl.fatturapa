package fatturapa

import (
	"errors"

	"github.com/invopop/gobl/org"
)

type address struct {
	Indirizzo    string
	NumeroCivico string `xml:",omitempty"`
	CAP          string
	Comune       string
	Provincia    string `xml:",omitempty"`
	Nazione      string
}

func newAddress(p *org.Party) (*address, error) {
	if len(p.Addresses) == 0 {
		return nil, errors.New("party missing address")
	}

	addr := p.Addresses[0]

	return &address{
		Indirizzo: addressLine(addr),
		CAP:       addr.Code,
		Comune:    addr.Locality,
		Provincia: addr.Region,
		Nazione:   addr.Country.String(),
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
