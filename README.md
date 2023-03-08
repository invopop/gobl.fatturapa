# gobl.fatturapa
GOBL Conversion into FatturaPA in Italy

# GOBL to FatturaPA Toolkit

Convert GOBL documents into the Italy's FatturaPA format.

TODO: copyright, license, build statuses

## Usage

### Go

There are a couple of entry points to build a new Fatturapa document. If you already have a GOBL Envelope available in Go, you could convert and output to a data file like this:

```golang
doc, err := fatturapa.NewInvoice(env)
if err != nil {
    panic(err)
}

data, err := doc.Bytes()
if err != nil {
    panic(err)
}

if err = os.WriteFile("./test.xml", data, 0644); err != nil {
    panic(err)
}
```

If you're loading from a file, you can use the `LoadGOBL` convenience method:

```golang
doc, err := fatturapa.LoadGOBL(file)
if err != nil {
    panic(err)
}
// do something with doc
```

Outputting to a FatturaPA XML is most useful when the document is signed. Use a certificate to sign the document as follows:

```golang
// import from github.com/invopop/xmldsig
cert, err := xmldsig.LoadCertificate(filename, password)
if err != nil {
    panic(err)
}

doc, err := fatturapa.NewInvoice(env, fatturapa.WithCertificate(cert))
if err != nil {
    panic(err)
}
```

### CLI

The command line interface can be useful for situations when you're using a language other than Golang in your application.

```bash
# install example
```

Simply provide the input GOBL JSON file and output to a file or another application:

```bash
./gobl.fatturapa convert input.json output.xml
```

If you have a digital certificate, run with:

```bash
./gobl.fatturapa convert -c cert.p12 -p password input.json > output.xml
```

The command also supports pipes:

```bash
cat input.json > ./gobl.fatturapa > output.xml
```

## Notes

- In all cases Go structures have been written using the same naming from the XML style document. This means names are not repeated in tags and generally makes it a bit easier map the XML output to the internal structures.

## Current Conversion Limitations

TODO

## Integration Tests

There are some integration and XML generation tests available in the `/test` path. To execute them, there are two [Magefile](https://magefile.org/) commands.

The first will convert YAML source data into GOBL JSON documents:

```
mage -v convertFromYAML
```

The second will generate the FatturaPA XML documents from the GOBL sources, using the digital certificates that are available in the `/test/certificates` path:

```
mage -v convertToXML
```

Sample data sources are contained in the `/test/data` directory. YAML and JSON (for tests) documents are stored in the Git repository, but the XML must be generated using the above commands.

## Mapping to GOBL

### DatiTrasmissione.ProgressivoInvio

The ProgressivoInvio field is a unique identifier assigned by the sender to each invoice that is being sent to the recipient, whereas the Numero field under DatiGeneraliDocumento is the invoice number assigned by the issuer of the invoice.

### DatiTrasmissione.FormatoTrasmissione

The FormatoTrasmissione field is used to specify the format of the invoice. The value of this field is either FPA12 (invoice to public administrations) or FPR12 (invoice to private parties)

### DatiTrasmissione.CodiceDestinatario

Number of the invoice office this invoice is being sent to. 6-digit code if FormatoTrasmissione is FPA12 and 7-digit code if FormatoTrasmissione is FPR12.

`GOBL`: Using `Inbox.Code` from `bill.Invoice.Supplier.Inboxes` where key is
`codice-destinario`.

### DatiPagamento.CondizioniPagamento

allowed values:

- TP01: Payment by instalments
- TP02: full payment
- TP03: advance payment

### DatiBeniServizi.DatiRiepilogo.EsigibilitaIVA

allowed values:

- I: VAT payable immediately
- D: deferred VAT payments
- S: split payments

