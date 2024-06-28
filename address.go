package fatturapa

import (
	"regexp"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

const (
	foreignCAP = "00000"
)

var (
	provinceRegexp = regexp.MustCompile(`^[A-Z]{2}$`)
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
		Provincia:    addressRegion(addr),
		Nazione:      addr.Country.String(),
	}
	if addr.Country == l10n.IT {
		ad.CAP = addr.Code
	} else {
		ad.CAP = foreignCAP
	}
	return ad
}

// addressRegion will simply check if the region is using the
// standard two digital capital letter code for the Italian province,
// or return an empty string to avoid FatturaPA validation issues.
// The province is optional, so it's not a problem if it's not set.
func addressRegion(address *org.Address) string {
	if address.Country == l10n.IT {
		if provinceRegexp.MatchString(address.Region) {
			return address.Region
		}
	}
	return ""
}

func addressStreet(address *org.Address) string {
	if address.PostOfficeBox != "" {
		return address.PostOfficeBox
	}
	return address.Street
}
