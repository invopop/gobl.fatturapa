package fatturapa

import (
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

// poBoxPrefixes contains common PO Box prefix patterns
var poBoxPrefixes = []string{"P.O. ", "PO Box ", "P.O.Box", "P.O Box", "PO BOX"}

func goblOrgAddressFromAddress(address *Address) *org.Address {
	addr := &org.Address{
		Locality: address.Locality,
		Number:   address.Number,
		Code:     cbc.Code(address.Code),
		Street:   address.Street,
	}

	// Convert country string to ISO country code
	addr.Country = l10n.ISOCountryCode(address.Country)

	// Handle region based on country
	if address.Country == l10n.IT.String() {
		addr.Region = address.Region
	}

	// Handle street vs post office box
	if address.Street != "" && isPostOfficeBox(address.Street) {
		addr.PostOfficeBox = address.Street
		addr.Street = ""
	}

	return addr
}

// isPostOfficeBox is a helper function to determine if a street address is actually a PO Box
func isPostOfficeBox(street string) bool {
	for _, prefix := range poBoxPrefixes {
		if strings.HasPrefix(street, prefix) {
			return true
		}
	}
	return false
}
