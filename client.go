package fatturapa

// Client contains information related to the entity using this library
// to submit invoices to SDI.
type Client struct {
	CountryCode string
	TaxID       string
}

func NewClient(countryCode, taxID string) Client {
	return Client{
		CountryCode: countryCode,
		TaxID:       taxID,
	}
}
