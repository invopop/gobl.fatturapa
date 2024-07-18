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

// Address from IndirizzoType
type Address struct {
	Street   string `xml:"Indirizzo"`              // Street
	Number   string `xml:"NumeroCivico,omitempty"` // Number
	Code     string `xml:"CAP"`                    // Post Code
	Locality string `xml:"Comune"`                 // Locality
	Region   string `xml:"Provincia,omitempty"`    // Region
	Country  string `xml:"Nazione"`                // Country Code
}

func newAddress(addr *org.Address) *Address {
	ad := &Address{
		Street:   addressStreet(addr),
		Number:   addr.Number,
		Locality: addr.Locality,
		Region:   addressRegion(addr),
		Country:  addr.Country.String(),
	}
	if addr.Country == l10n.IT {
		ad.Code = addr.Code
	} else {
		ad.Code = foreignCAP
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
