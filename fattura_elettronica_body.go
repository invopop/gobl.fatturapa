package fatturapa

import (
	"github.com/invopop/gobl/bill"
)

type FatturaElettronicaBody struct {
}

func newFatturaElettronicaBody(inv bill.Invoice) (*FatturaElettronicaBody, error) {
	return &FatturaElettronicaBody{}, nil
}
