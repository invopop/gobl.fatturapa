package sdi

// Based on "MessaggiTypes.xsd"

// ReceiptDeliveryMessage represents the structure for "RicevutaConsegna" message.
type ReceiptDeliveryMessage struct {
	SDIIdentifier            int                   `xml:"IdentificativoSdI"`
	FileName                 string                `xml:"NomeFile"`
	ReceptionDateTime        string                `xml:"DataOraRicezione"`
	DeliveryDateTime         string                `xml:"DataOraConsegna"`
	Recipient                RecipientType         `xml:"Destinatario"`
	ArchiveReference         *ArchiveReferenceType `xml:"RiferimentoArchivio,omitempty"`
	MessageID                string                `xml:"MessageId"`
	PecMessageID             *string               `xml:"PecMessageId,omitempty"`
	Notes                    string                `xml:"Note,omitempty"`
	Signature                string                `xml:"ds:Signature"`
	Version                  string                `xml:"versione,attr"`
	IntermediaryWithDualRole string                `xml:"IntermediarioConDupliceRuolo,attr,omitempty"`
}

// RejectionMessage represents the structure for "NotificaScarto" message.
type RejectionMessage struct {
	SDIIdentifier     int                   `xml:"IdentificativoSdI"`
	FileName          string                `xml:"NomeFile"`
	ReceptionDateTime string                `xml:"DataOraRicezione"`
	ArchiveReference  *ArchiveReferenceType `xml:"RiferimentoArchivio,omitempty"`
	ErrorList         ErrorListType         `xml:"ListaErrori"`
	MessageID         string                `xml:"MessageId"`
	PecMessageID      *string               `xml:"PecMessageId,omitempty"`
	Notes             string                `xml:"Note,omitempty"`
	Signature         string                `xml:"ds:Signature"`
	Version           string                `xml:"versione,attr"`
}

// NonDeliveryNotificationMessage represents the structure for "NotificaMancataConsegna" message.
type NonDeliveryNotificationMessage struct {
	SDIIdentifier     int                   `xml:"IdentificativoSdI"`
	FileName          string                `xml:"NomeFile"`
	ReceptionDateTime string                `xml:"DataOraRicezione"`
	ArchiveReference  *ArchiveReferenceType `xml:"RiferimentoArchivio,omitempty"`
	Description       string                `xml:"Descrizione,omitempty"`
	MessageID         string                `xml:"MessageId"`
	PecMessageID      *string               `xml:"PecMessageId,omitempty"`
	Notes             string                `xml:"Note,omitempty"`
	Signature         string                `xml:"ds:Signature"`
	Version           string                `xml:"versione,attr"`
}

// OutcomeNotificationMessage represents the structure for "NotificaEsito" message.
type OutcomeNotificationMessage struct {
	SDIIdentifier            int     `xml:"IdentificativoSdI"`
	FileName                 string  `xml:"NomeFile"`
	ClientOutcome            string  `xml:"EsitoCommittente"`
	MessageID                string  `xml:"MessageId"`
	PecMessageID             *string `xml:"PecMessageId,omitempty"`
	Notes                    string  `xml:"Note,omitempty"`
	Signature                string  `xml:"ds:Signature"`
	Version                  string  `xml:"versione,attr"`
	IntermediaryWithDualRole string  `xml:"IntermediarioConDupliceRuolo,attr,omitempty"`
}

// InvoiceTransmissionCertificate represents the structure for "AttestazioneTrasmissioneFattura" message.
type InvoiceTransmissionCertificate struct {
	SDIIdentifier     int                   `xml:"IdentificativoSdI"`
	FileName          string                `xml:"NomeFile"`
	ReceptionDateTime string                `xml:"DataOraRicezione"`
	ArchiveReference  *ArchiveReferenceType `xml:"RiferimentoArchivio,omitempty"`
	Recipient         RecipientType         `xml:"Destinatario"`
	MessageID         string                `xml:"MessageId"`
	PecMessageID      *string               `xml:"PecMessageId,omitempty"`
	Notes             string                `xml:"Note,omitempty"`
	OriginalFileHash  string                `xml:"HashFileOriginale"`
	Signature         string                `xml:"ds:Signature"`
	Version           string                `xml:"versione,attr"`
}

// FileSubmissionMetadata represents the structure for "MetadatiInvioFile" message.
type FileSubmissionMetadata struct {
	SDIIdentifier      int    `xml:"IdentificativoSdI"`
	FileName           string `xml:"NomeFile"`
	RecipientCode      string `xml:"CodiceDestinatario"`
	Format             string `xml:"Formato"`
	SubmissionAttempts int    `xml:"TentativiInvio"`
	MessageID          string `xml:"MessageId"`
	Notes              string `xml:"Note,omitempty"`
	Version            string `xml:"versione,attr"`
}

// ClientOutcomeNotification represents the structure for "NotificaEsitoCommittente" message.
type ClientOutcomeNotification struct {
	SDIIdentifier    int                   `xml:"IdentificativoSdI"`
	InvoiceReference *InvoiceReferenceType `xml:"RiferimentoFattura,omitempty"`
	Outcome          string                `xml:"Esito"`
	Description      string                `xml:"Descrizione,omitempty"`
	ClientMessageID  *string               `xml:"MessageIdCommittente,omitempty"`
	Signature        *string               `xml:"ds:Signature,omitempty"`
	Version          string                `xml:"versione,attr"`
}

// DeadlineNotification represents the structure for "NotificaDecorrenzaTermini" message.
type DeadlineNotification struct {
	SDIIdentifier            int                   `xml:"IdentificativoSdI"`
	InvoiceReference         *InvoiceReferenceType `xml:"RiferimentoFattura,omitempty"`
	FileName                 string                `xml:"NomeFile"`
	Description              string                `xml:"Descrizione,omitempty"`
	MessageID                string                `xml:"MessageId"`
	PecMessageID             *string               `xml:"PecMessageId,omitempty"`
	Notes                    string                `xml:"Note,omitempty"`
	Signature                string                `xml:"ds:Signature"`
	Version                  string                `xml:"versione,attr"`
	IntermediaryWithDualRole string                `xml:"IntermediarioConDupliceRuolo,attr,omitempty"`
}

// RecipientType represents the recipient information.
type RecipientType struct {
	Code        string `xml:"Codice"`      // Recipient code (6-7 alphanumeric characters)
	Description string `xml:"Descrizione"` // Recipient description
}

// ArchiveReferenceType represents a reference to an archive.
type ArchiveReferenceType struct {
	SDIIdentifier int    `xml:"IdentificativoSdI"` // SDI identifier (exactly 12 digits)
	FileName      string `xml:"NomeFile"`          // File name (maximum 50 characters)
}

// InvoiceReferenceType struct represents the XML structure for RiferimentoFattura_Type element.
type InvoiceReferenceType struct {
	InvoiceNumber   string `xml:"NumeroFattura"`              // Invoice number (1-20 latin characters)
	InvoiceYear     uint   `xml:"AnnoFattura"`                // Invoice year (non negative integer)
	InvoicePosition *uint  `xml:"PosizioneFattura,omitempty"` // Invoice position (optional, positive integer)
}
