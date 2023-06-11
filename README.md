# gobl.fatturapa

GOBL Conversion into FatturaPA in Italy

# GOBL to FatturaPA Toolkit

Convert GOBL documents into the Italy's FatturaPA format.

TODO: copyright, license, build statuses

## Usage

### Go

There are a couple of entry points to build a new Fatturapa document. If you already have a GOBL Envelope available in Go, you could convert and output to a data file like this:

```golang
converter := fatturapa.NewConverter()

doc, err := converter.ConvertFromGOBL(env)
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

See the following example for signing the XML with a certificate:

```golang
// import from github.com/invopop/xmldsig
cert, err := xmldsig.LoadCertificate(filename, password)
if err != nil {
    panic(err)
}

converter := fatturapa.NewConverter(
    fatturapa.WithCertificate(cert),
	fatturapa.WithTimestamp(), // if you want to include a timestamp in the digital signature
)

doc, err := converter.ConvertFromGOBL(env)
if err != nil {
    panic(err)
}
```

If you want to include the fiscal data of the entity integrating with the SDI (Italy's e-invoice system) and `ProgressivoInvio` (transmission number) in the XML, you can use the `WithTransmitterData` option. This option must be used if you are integrating diredctly with the SDI, but if you are working with a third party service to send the XML, it would be on their side to include this data.

```golang
transmitter := fatturapa.Transmitter{
    CountryCode: countryCode, // ISO 3166-1 alpha-2
    TaxID:       taxID,       // Valid tax ID of transmitter
}

converter := fatturapa.NewConverter(
    fatturapa.WithTransmitterData(transmitter),
    // other options
)
```

### CLI

The command line interface can be useful for situations when you're using a language other than Golang in your application.

```bash
cd cmd/gobl.fatturapa
go build
```

Simply provide the input GOBL JSON file and output to a file or another application:

```bash
./gobl.fatturapa convert input.json output.xml
```

If you have a digital certificate, run with:

```bash
./gobl.fatturapa convert -c cert.p12 -p password input.json output.xml
```

To include the transmitter information, add the `-T` flag and provide the _country code_ and the _tax ID_:

```bash
./gobl.fatturapa convert -T ES12345678 input.json output.xml
```

The command also supports pipes:

```bash
cat input.json > ./gobl.fatturapa output.xml
```

## Notes

- In all cases Go structures have been written using the same naming from the XML style document. This means names are not repeated in tags and generally makes it a bit easier map the XML output to the internal structures.

## Integration Tests

There are some integration and XML generation tests available in the `/test` path. to generate the FatturaPA XML documents from the GOBL sources, use the digital certificates that are available in the `/test/certificates` path:

```
mage -v TestConversion
```

Sample data sources are contained in the `/test/data` directory. JSON (for tests) documents are stored in the Git repository, but the XML must be generated using the above commands.

## Current Conversion Limitations

The FatturaPA XML schema is quite large and complex. This library is not complete and only supports a subset of the schema. The current implementation is focused on the most common use cases.

- FatturaPA allows multiple invoices within the document, but this library only supports a single invoice per transmission.
- Only a subset of payment methods (ModalitaPagamento) are supported. See `payments.go` for the list of supported codes.

Some of the optional elements currently not supported include:

- `Allegati` (attachments)
- `DatiOrdineAcquisto` (data related to purchase orders)
- `DatiContratto` (data related to contracts)
- `DatiConvenzione` (data related to conventions)
- `DatiRicezione` (data related to receipts)
- `DatiFattureCollegate` (data related to linked invoices)
- `DatiBollo` (data related to duty stamps)
