package fatturapa

// TaxID is the VAT identification number consisting of a country code and the
// actual VAT number.
type TaxID struct {
	// ISO 3166-1 alpha-2 country code
	IdPaese string
	// Actual VAT number
	IdCodice string
}
