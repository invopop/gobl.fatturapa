package fatturapa

// taxID is the VAT identification number consisting of a country code and the
// actual VAT number.
type taxID struct {
	// ISO 3166-1 alpha-2 country code
	IdPaese string // nolint:revive
	// Actual VAT number
	IdCodice string // nolint:revive
}
