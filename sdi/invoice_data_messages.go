package sdi

// Based on "DatiFatturaMessaggi.xsd"

// OutcomeFile struct represents the XML structure for EsitoFile element.
type OutcomeFile struct {
	Type         string          `xml:"TipoFile"`
	ID           int             `xml:"IDFile"`
	Name         string          `xml:"NomeFile"`
	DateTime     string          `xml:"DataOraRicezione"`
	ArchiveRef   *ArchiveRefType `xml:"RifArchivio,omitempty"`
	Outcome      string          `xml:"Esito"`
	ErrorList    *ErrorListType  `xml:"ListaErrori,omitempty"`
	MessageID    string          `xml:"MessageID"`
	PECMessageID *string         `xml:"PECMessageID,omitempty"`
	Note         string          `xml:"Note,omitempty"`
	Signature    string          `xml:"http://www.w3.org/2000/09/xmldsig# Signature"`
	Version      string          `xml:"versione,attr"`
}

// ArchiveRefType represents the RifArchivio element.
type ArchiveRefType struct {
	ID   int    `xml:"IDArchivio"`
	Name string `xml:"NomeArchivio"`
}

// ErrorListType represents the ListaErrori element.
type ErrorListType struct {
	Error []ErrorType `xml:"Errore"`
}

// ErrorType represents the Errore element.
type ErrorType struct {
	Code        string `xml:"Codice"`
	Description string `xml:"Descrizione"`
}
