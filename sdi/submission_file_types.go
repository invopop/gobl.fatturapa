package sdi

// SubmissionFile represents the "fileSdI" element
// This type is defined in both receipt types and transmission types.
type SubmissionFile struct {
	SdiIdentifier string `xml:"IdentificativoSdI,omitempty"`
	FileName      string `xml:"NomeFile"`
	FileData      []byte `xml:"File"`
}
