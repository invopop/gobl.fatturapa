package sdi

// Based on "TransmissioneTypes.xsd"

// ReceivedFile represents the "fileSdI" element from the XSD.
type ReceivedFile struct {
	SdiIdentifier string `xml:"IdentificativoSdI"`
	FileName      string `xml:"NomeFile"`
	FileData      []byte `xml:"File"`
}

// ReceiveFileResponse represents the "rispostaSdIRiceviFile" element from the XSD.
type ReceiveFileResponse struct {
	SdiIdentifier string        `xml:"IdentificativoSdI"`
	ReceiptTime   string        `xml:"DataOraRicezione"`
	Error         *ErrorReceipt `xml:"Errore,omitempty"`
}

// ErrorReceipt represents the error type from the XSD.
type ErrorReceipt struct {
	ErrorCode string `xml:",innerxml"`
}
