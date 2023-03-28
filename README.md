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

## Usage

### CLI

Copy `.env.example` to `.env` and update the values to configure the application.

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

## Current Conversion Limitations

The FatturaPA XML schema is quite large and complex. This library is not complete and only supports a subset of the schema. The current implementation is focused on the most common use cases.

- FatturaPA allows multiple invoices within the document, but this library only supports a single invoice per transmission.
- `DatiBeniServizi.DatiRiepilogo.EsigibilitaIVA` code defaults to "I" (immediata) for all invoices.

Some of the optional elements currently not supported include:
- `Allegati` (attachments)
- `DatiOrdineAcquisto` (data related to purchase orders)
- `DatiContratto` (data related to contracts)
- `DatiConvenzione` (data related to conventions)
- `DatiRicezione` (data related to receipts)
- `DatiFattureCollegate` (data related to linked invoices)
- `DatiBollo` (data related to duty stamps)
