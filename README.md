# gobl.fatturapa
GOBL Conversion into FatturaPA in Italy

## Elements Not covered in GOBL

### DatiTrasmissione.ProgressivoInvio

The ProgressivoInvio field is a unique identifier assigned by the sender to each invoice that is being sent to the recipient, whereas the Numero field under DatiGeneraliDocumento is the invoice number assigned by the issuer of the invoice.

### DatiTrasmissione.FormatoTrasmissione

The FormatoTrasmissione field is used to specify the format of the invoice. The value of this field is either FPA12 (invoice to public administrations) or FPR12 (invoice to private parties)

### DatiTrasmissione.CodiceDestinatario

Number of the invoice office this invoice is being sent to. 6-digit code if FormatoTrasmissione is FPA12 and 7-digit code if FormatoTrasmissione is FPR12.

### DatiPagamento.CondizioniPagamento

allowed values:

- TP01: Payment by instalments
- TP02: full payment
- TP03: advance payment

### DatiBeniServizi.DatiRiepilogo.EsigibilitaIVA

allowed values:

- I: VAT payable immediately
- D: unrealized VAT
- S: split payments

