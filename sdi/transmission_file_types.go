package sdi

// Based on "TrasmissioneFileTypes.xsd"

// TransmissionFile struct represents the XML structure for file element.
type TransmissionFile struct {
	Name string `xml:"NomeFile"`
	Type string `xml:"TipoFile"`
	Data []byte `xml:"File"`
}

// ResponseFile struct represents the XML structure for rispostaFile element.
type ResponseFile struct {
	ID       int     `xml:"IDFile"`
	DateTime string  `xml:"DataOraRicezione"`
	Error    *string `xml:"Errore,omitempty"`
}

// Outcome struct represents the XML structure for esito element.
type Outcome struct {
	ID int `xml:"IDFile"`
}

// ResponseOutcome struct represents the XML structure for rispostaEsito element.
type ResponseOutcome struct {
	Status        string           `xml:"Esito"`
	Notification  *Notification    `xml:"Notifica,omitempty"`
	ArchiveDetail []*ArchiveDetail `xml:"DettaglioArchivio,omitempty"`
	Error         *string          `xml:"Errore,omitempty"`
}

// Notification struct represents the XML structure for notifica element.
type Notification struct {
	Name string `xml:"NomeFile"`
	File []byte `xml:"File"`
}

// ArchiveDetail struct represents the XML structure for dettaglioArchivio element.
type ArchiveDetail struct {
	Name string `xml:"NomeFile"`
	ID   int    `xml:"IDFile"`
}
