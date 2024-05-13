package sdi

// Based on "RicezioneTypes.xsd"

// ReceiveInvoicesResponse represents the "rispostaRiceviFatture" element
type ReceiveInvoicesResponse struct {
	Outcome string `xml:"Esito"`
}

// FileWithMetadata represents the "fileSdIConMetadati" element
type FileWithMetadata struct {
	SdiIdentifier        string `xml:"IdentificativoSdI"`
	FileName             string `xml:"NomeFile"`
	FileData             []byte `xml:"File"`
	MetadataFileName     string `xml:"NomeFileMetadati"`
	MetadataFileContents []byte `xml:"Metadati"`
}

// NotifyOutcomeResponse represents the "rispostaSdINotificaEsito" element
type NotifyOutcomeResponse struct {
	Outcome      string          `xml:"Esito"`
	RejectedFile *SubmissionFile `xml:"ScartoEsito,omitempty"`
}
