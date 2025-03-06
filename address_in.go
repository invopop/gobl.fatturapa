package fatturapa

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

func goblOrgAddressFromAddress(address *Address) *org.Address {
	addr := &org.Address{
		Locality: address.Locality,
		Number:   address.Number,
		Street:   address.Street,
	}

	// Convert country string to ISO country code
	addr.Country = l10n.ISOCountryCode(address.Country)

	// Handle region based on country
	if address.Country == l10n.IT.String() {
		addr.Region = address.Region
		addr.Code = cbc.Code(address.Code)
	} else {
		// Only use region if it matches the province pattern for Italy
		if provinceRegexp.MatchString(address.Region) {
			addr.Region = address.Region
		}
		// For non-Italian addresses, we don't set the code
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
	// Simple implementation - could be expanded with regex patterns
	return len(street) >= 5 && (street[:5] == "P.O. " || street[:7] == "PO Box ")
}
