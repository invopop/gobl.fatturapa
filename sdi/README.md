# SdI (Sistema di Interscambio)

## SDICoop Web Service

### Test endpoint

- [Test URI](https://testservizi.fatturapa.it/)

### Production endpoint

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

## Receive an invoice

To receive invoices and notifications, you need to set up a server communicating with SdI.

### Development

To communicate with the server the client need:

- file with CA certificates in PEM format
- PEM certificate (see the Certificator_Server folder)
- PEM RSA private key (generated during initial registration process)

The server must use HTTPS.
Since certificates are associated with a domain,
we need to convince our local DNS that this is the right domain.

In this example the domain used will be `sdi-it.invopop.com`.

Set this domain in the `/etc/hosts` file to point to `0.0.0.0`.

```console
$ cat /etc/hosts | grep sdi
0.0.0.0		sdi-it.invopop.com
```

Run the server in one console:

```console
$ go run ./cmd/gobl.fatturapa server --ca-cert ./ca-all.pem --cert ./SDI-IT.INVOPOP.COM.pem --key ./key_server.key --verbose sdi-it.invopop.com 8080
Server start: sdi-it.invopop.com:8080
Client auth: RequireAndVerifyClientCert
Incoming request:
GET / HTTP/2.0
Host: sdi-it.invopop.com:8080
Accept: */*
User-Agent: curl/7.88.1

Outgoing response:
HTTP/2.0 200 OK
Content-Length: 2
Content-Type: text/plain

OK
```

In another console, send requests:

```console
$ curl --cacert ./ca-all.pem --cert ./SDI-IT.INVOPOP.COM.pem --key ./key_server.key https://sdi-it.invopop.com:8080
OK
```

### Production

#### Pre-requirements

Read [instructions for creating the CSR using OpenSSL commands](https://www.fatturapa.gov.it/it/norme-e-regole/DocumentazioneSDI/) / SDICoop Service / Example of CSR generation using OpenSSL commands (for expert users).

To ensure proper functionality, the Common Name (CN) of the SSL certificate
should match the DNS name of the endpoint it is used with.
Mismatch lead to authentication issues.

SdI uses two-way SSL authentication, where both the client (SdI)
and the server (your endpoint) exchange certificates for authentication.
The server must have a certificate issued by
[Agenzia delle Entrate](https://www.agenziaentrate.gov.it/portale/)
configured to the IP address, not the domain,
due to SdI not supporting SNI (Server Name Indication).

Failure to meet the above requirements will result in communication
being blocked at the TLS level.
Any attempts from SdI will only be visible in the simulation section
of the "Manage the channel" interface.

---

## TODO

- [X] completed tasks
- [ ] uncompleted tasks

Below is a checklist.
It is intended to help us determine the current status of the project.

### Sending invoices

The functionality related to sending invoices
and simple invoices works correctly (tested in test env for SdI).
The status of sent invoices is visible
in the "Interoperability Test" section of "Manage the channel".

- [X] Sending invoices to SdI
  - [X] Building the client along with configuration
  - [X] Preparation of SOAP request for SdI
  - [X] Creation of "transmit" command for CLI
  - [X] Adding SSL support to the HTTP client for sending invoices
- [X] Receiving responses after sending invoices:
  - [X] Parses a multipart HTTP response (XOP+XML)
  - [X] Deserialization of response into appropriate structure
- [X] Handling errors:
  - [X] Empty invoice (File vuoto)
  - [X] Service unavailable (Servizio non disponibile)
  - [X] Unauthorized user (Utente non abilitato)
  - [X] File type not correct (Tipo file non corretto)
- [X] Preparation of structures for Message Types (based on MessaggiTypes_v1.1.xsd)

### Receiving invoices

The functionality of receiving invoices and notifications is not working
because we couldn't configure SSL communication correctly.
When we send an invoice to the recipient specified in the "Manage the channel"
simulation section, the our server should receive it.
SdI makes three attempts to deliver the invoice or notification.
However, for this to happen, SdI tries to authenticate the server
it sends data to at the SSL level.

If any requests were to arrive, we would see them in the server logs.
Currently, the status of sent requests is visible only in "Manage the channel",
with the following message:

> javax.net.ssl.SSLHandshakeException [1]: General SSLEngine problem

Despite SSL issues, some tasks were successfully completed (based on tests).

- [ ] Receiving invoices from SdI
  - [X] Simple HTTP server listening on the selected port
  - [X] Adding SSL support to the HTTPS server
  - [X] Building a message handler
  - [ ] Handling transmission service endpoint (out)
  - [ ] Handling reception service endpoint (in)
  - [X] Creation of "server" command for CLI
  - [ ] Parsing request from SdI with the invoice
  - [ ] Assign parsed invoice to GOBL format
- [ ] Receiving notifications from SdI
  - [X] Handling different SdI requests
  - [ ] Receiving actual requests from SDI for testing purposes
  - [X] Mocking real communication with SDI for testing purposes
- [X] Parsing SDI messages:
  - [X] Rejection Receipt (RicevutaScarto)
  - [ ] Invoice Transmission Confirmation (AttestazioneTrasmissioneFattura)
  - [X] File Submission Metadata (MetadatiInvioFile)
  - [ ] Delivery Failure Notification (NotificaMancataConsegna)
- [ ] Server configuration
  - [X] Analysis of requirements needed for server configuration
  - [ ] Setting up the server for "Interoperability Tests"
  - [ ] Setting up the server for production
- [X] Preparing the structure based on XSD files:
  - [X] Invoice Data Message (DatiFatturaMessaggi_v2.0.xsd)
  - [X] Receipt Types (RicezioneTypes_v1.0.xsd)
  - [X] Submission File Types
  - [X] Transmission File Types (TrasmissioneFileTypes_v2.0.xsd)
  - [X] Transmission Types (TransmissioneTypes_v1.1.xsd)

The above tasks likely do not cover all the preparations required.
However, without the ability to receive requests from SdI,
the development process is hindered.

### Interoperability Test

To pass the "Interoperability Test", a correctly configured environment
with an active server is required. SdI does not allow the use of two different
certificates, so the staging environment should be on the same domain
as the production environment. The difference may lie in a different path,
which can be set in the "Change Endpoint" section of "Manage the Channel".

Necessary tests to pass:

- [ ] Invoice Reception (Ricezione Fattura)
- [ ] Delivery Receipt (Ricevuta consegna)
- [ ] Non-Delivery Notification (B2G)/Undeliverable Notification (B2B, B2C) (Notifica mancata consegna (B2G)/Notifica impossibilità di recapito (B2B, B2C))
- [ ] Rejection Notification (B2G)/Rejection Receipt (B2B, B2C) (Notifica scarto (B2G)/Ricevuta di scarto (B2B, B2C))

Further tests for FatturaPA:

- [ ] Outcome Notification from PA (Notifica di esito da PA)
- [ ] Rejection Notification of Outcome to PA (Notifica di Scarto esito a PA)
- [ ] Deadline Notification to PA (Notifica Decorrenza Termini a PA)
- [ ] Outcome Notification to Economic Operator (Notifica esito a Operatore Economico)
- [ ] Deadline Notification to Economic Operator (Notifica Decorrenza Termini a Operatore Economico)
- [ ] Transmission Confirmation (Attestazione avvenuta trasmissione)

After successfully passing the interoperability tests, SdI will be unlocked
to communicate with the production environment of the SDICoop Web Service.
