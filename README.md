# GOBL to FatturaPA Tools

Convert GOBL into the Italy's FatturaPA format.

Copyright [Invopop Ltd.](https://invopop.com) 2023. Released publicly under the [Apache License Version 2.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com). In order to accept contributions to this library we will require transferring copyrights to Invopop Ltd.

[![Lint](https://github.com/invopop/gobl.factturapa/actions/workflows/lint.yaml/badge.svg)](https://github.com/invopop/gobl.fatturapa/actions/workflows/lint.yaml)
[![Test Go](https://github.com/invopop/gobl.fatturapa/actions/workflows/test.yaml/badge.svg)](https://github.com/invopop/gobl.fatturapa/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/invopop/gobl.fatturapa)](https://goreportcard.com/report/github.com/invopop/gobl.fatturapa)
[![GoDoc](https://godoc.org/github.com/invopop/gobl.fatturapa?status.svg)](https://godoc.org/github.com/invopop/gobl.fatturapa)
![Latest Tag](https://img.shields.io/github/v/tag/invopop/gobl.fatturapa)

## Introduction

FatturaPA defines two versions of invoices:

- Ordinary invoices, `FatturaElettronica` types `FPA12` and `FPR12` defined in the v1.2 schema, usable for all sales.
- Simplified invoices, `FatturaElettronicaSemplificata` type `FSM10` defined in the v1.0 schema, with a reduced set of requirements but can only be used for sales of less then €400, as of writing. **Currently not supported!**

Unlike other tax regimes, Italy requires simplified invoices to include the customer's tax ID. For "cash register" style receipts locally called "Scontrinos", another format and API is used for this from approved hardware.

## Sources

You can find copies of the Italian FatturaPA schema in the [schemas folder](./schema).

Key websites:
- [FatturaPA & SDI service documentation page on  on Italy's tax authority's website ] (https://www.agenziaentrate.gov.it/portale/web/guest/fatturazione-elettronica-e-dati-fatture-transfrontaliere-new)
- [FatturaPA documentation page on FatturaPA's dedicated website](https://www.fatturapa.gov.it/en/norme-e-regole/documentazione-fattura-elettronica/formato-fatturapa/)

Useful files:
- [Ordinary Schema V1.2.1 Spec Table View (EN)](https://www.fatturapa.gov.it/export/documenti/fatturapa/v1.2.1/Table-view-B2B-Ordinary-invoice.pdf) - by far the most comprehensible spec doc. Since the difference between 1.2.2 and 1.2.1 is minimal, this is perfectly usable.
- [Ordinary Schema V1.2.2 PDF (IT)](https://www.fatturapa.gov.it/export/documenti/Specifiche_tecniche_del_formato_FatturaPA_v1.3.1.pdf) - most up-to-date but difficult
- [XSD V1.2.2](https://www.fatturapa.gov.it/export/documenti/fatturapa/v1.2.2/Schema_del_file_xml_FatturaPA_v1.2.2.xsd)
- [XSD V1 (FSM10) - simplified invoices](https://www.agenziaentrate.gov.it/portale/documents/20143/288192/ST+Fatturazione+elettronica+-+Schema+VFSM10_Schema_VFSM10.xsd/010f1b41-6683-1b31-ba36-c8bced659c06)

## Limitations

The FatturaPA XML schema is quite large and complex. This library is not complete and only supports a subset of the schema. The current implementation is focused on the most common use cases.

- Simplified invoices are not currently supported (please get in touch if you need this).
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

The command line interface can be useful for situations when you're using a language other than Golang in your application. Install with:

```bash
go install github.com/invopop/gobl.fatturapa
```

Simply provide the input GOBL JSON file and output to a file or another application:

```bash
gobl.fatturapa convert input.json output.xml
```

If you have a digital certificate, run with:

```bash
gobl.fatturapa convert -c cert.p12 -p password input.json output.xml
```

To include the transmitter information, add the `-T` flag and provide the _country code_ and the _tax ID_:

```bash
gobl.fatturapa convert -T ES12345678 input.json output.xml
```

The command also supports pipes:

```bash
cat input.json > ./gobl.fatturapa output.xml
```

## Notes

- In all cases Go structures have been written using the same naming from the XML style document. This means names are not repeated in tags and generally makes it a bit easier to map the XML output to the internal structures.

## Integration Tests

There are some integration and XML generation tests available in the `/test` path. to generate the FatturaPA XML documents from the GOBL sources, use the digital certificates that are available in the `/test/certificates` path:

```
mage -v TestConversion
```

Sample data sources are contained in the `/test/data` directory. JSON (for tests) documents are stored in the Git repository, but the XML must be generated using the above commands.
