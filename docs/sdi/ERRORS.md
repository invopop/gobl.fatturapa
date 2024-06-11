# Error List

## Empty SOAP request

Response to an empty string sent instead of a SOAP request with invoice.

```xml
<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
  <soapenv:Body>
    <soapenv:Fault>
      <soapenv:Code>
        <soapenv:Value>soapenv:Receiver</soapenv:Value>
      </soapenv:Code>
      <soapenv:Reason>
        <soapenv:Text xml:lang="en-US">javax.xml.stream.XMLStreamException: The root element is required in a well-formed document.</soapenv:Text>
      </soapenv:Reason>
      <soapenv:Detail></soapenv:Detail>
    </soapenv:Fault>
  </soapenv:Body>
</soapenv:Envelope>
```

Status Code: 500 Internal Server Error

## SOAP request without Body

Response to sending SOAP request without `<soapenv:Body>` part.

```xml
<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
  <soapenv:Body>
    <soapenv:Fault xmlns:axis2ns1="http://schemas.xmlsoap.org/soap/envelope/">
      <faultcode>axis2ns1:Server</faultcode>
      <faultstring>Internal Error</faultstring>
      <detail></detail>
    </soapenv:Fault>
  </soapenv:Body>
</soapenv:Envelope>
```

Status Code: 500 Internal Server Error

## Empty invoice

Response to sending an empty invoice file.

```xml
<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
  <soapenv:Body>
    <ns2:rispostaSdIRiceviFile
      xmlns:ns2="http://www.fatturapa.gov.it/sdi/ws/trasmissione/v1.0/types">
      <IdentificativoSdI>0</IdentificativoSdI>
      <DataOraRicezione>2024-06-11T12:00:00.000+02:00</DataOraRicezione>
      <Errore>EI01</Errore>
    </ns2:rispostaSdIRiceviFile>
  </soapenv:Body>
</soapenv:Envelope>
```

Status Code: 200 OK
