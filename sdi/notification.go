package sdi

import (
	"encoding/xml"
	"fmt"
)

// ReceiptRejection represents the structure for "RicevutaScarto" message.
type ReceiptRejection struct {
	XMLName    xml.Name `xml:"RicevutaScarto"`
	HashNumber string   `xml:"Hash"`
	RejectionMessage
}

// ParseReceiptRejection converts XML to a Receipt Rejection structure
func ParseReceiptRejection(receipt []byte) (ReceiptRejection, error) {
	var rr ReceiptRejection
	err := xml.Unmarshal(receipt, &rr)
	if err != nil {
		return rr, fmt.Errorf("xml parsing error: %v", err)
	}

	errors := rr.ErrorList.Error
	if len(errors) > 0 {
		errCodes := make([]string, len(errors))
		for i, err := range errors {
			errCodes[i] = err.Code
		}
		return rr, fmt.Errorf("sdi error code list: %v", errCodes)
	}

	return rr, nil
}
