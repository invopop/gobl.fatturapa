# SdI (Sistema di Interscambio)

## SDICoop Web Service

### Test

- [Test URI](https://testservizi.fatturapa.it/)

### Production

- [Production URI](https://fatturapa.it/)

## Manage the channel

1. Go to the
  [Manage the channel](https://sdi.fatturapa.gov.it/SdI2FatturaPAWebSpa/GestireCanaleAction.do)
  page.
2. Choose file to load the initial registration request.
  The required file is `RichiestaAccreditamento.zip.p7m`.
  Note the extension "zip.p7m" of the file.
3. Fill the "Security Code" field. This is a captcha.

After logging in, we have access to the following options:

- Interoperability Test
- View service agreement
- Change Endpoint

## Send an invoice

To send an invoice you need:
- invoice file in XML format
- file with CA certificates in PEM format
- PEM certificate (see the Certificato_Client folder)
- PEM RSA private key (generated during initial registration process)

Example of use:

```console
$ go run ./cmd/gobl.fatturapa transmit --ca-cert ./ca-all.pem --cert ./SDI-B85905495.pem --key ./key_client.key --verbose --env test invoice.xml
2024/05/28 17:00:00.000000 DEBUG RESTY
==============================================================================
~~~ REQUEST ~~~
POST  /ricevi_file  HTTP/1.1
HOST   : testservizi.fatturapa.it
HEADERS:
	Content-Type: text/plain; charset=utf-8
	User-Agent: go-resty/2.13.1 (https://github.com/go-resty/resty)
BODY   :
<?xml version='1.0' encoding='UTF-8'?><soapenv:Envelope xmlns:soapenv='http://schemas.xmlsoap.org/soap/envelope/' xmlns:typ='http://www.fatturapa.gov.it/sdi/ws/trasmissione/v1.0/types'><soapenv:Header/><soapenv:Body><typ:fileSdIAccoglienza><NomeFile>invoice.xml</NomeFile><File>FILE CONTENT</File></typ:fileSdIAccoglienza></soapenv:Body></soapenv:Envelope>
------------------------------------------------------------------------------
~~~ RESPONSE ~~~
STATUS       : 200 OK
PROTO        : HTTP/1.1
RECEIVED AT  : 2024-05-28T17:00:00.500000000+02:00
TIME DURATION: 500.000000ms
HEADERS      :
	Cache-Control: no-cache="set-cookie, set-cookie2"
	Content-Language: en-US
	Content-Length: 716
	Content-Type: multipart/related; boundary="MIMEBoundary_"; type="application/xop+xml"; start="<0.182152212d14b303688146ac0db42b507ac086a479534824@apache.org>"; start-info="text/xml"
	Date: Tue, 28 May 2024 15:00:00 GMT
	Expires: Thu, 01 Dec 1994 16:00:00 GMT
	Set-Cookie: ==; HTTPOnly; Path=/; Domain=.fatturapa.it; HttpOnly
	X-Powered-By: Servlet/3.0
BODY         :
--MIMEBoundary_
Content-Type: application/xop+xml; charset=utf-8; type="text/xml"
Content-Transfer-Encoding: binary
Content-ID: <0.182152212d14b303688146ac0db42b507ac086a479534824@apache.org>

<?xml version="1.0" encoding="utf-8"?><soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"><soapenv:Body><ns2:rispostaSdIRiceviFile xmlns:ns2="http://www.fatturapa.gov.it/sdi/ws/trasmissione/v1.0/types"><IdentificativoSdI>29208092</IdentificativoSdI><DataOraRicezione>2024-05-28T17:00:00.200+02:00</DataOraRicezione></ns2:rispostaSdIRiceviFile></soapenv:Body></soapenv:Envelope>
--MIMEBoundary_--
==============================================================================
```

## Links

- [Manage the channel](https://sdi.fatturapa.gov.it/SdI2FatturaPAWebSpa/GestireCanaleAction.do)

## Pfx server certificate

```
$ openssl pkcs12 -export -certpbe PBE-SHA1-3DES -keypbe PBE-SHA1-3DES -nomac -out SDI-PIVA-SERVER.pfx -inkey key_server.key -in SDI-IT.INVOPOP.COM.pem -certfile ca-all.pem
Enter Export Password:
Verifying - Enter Export Password:
```
