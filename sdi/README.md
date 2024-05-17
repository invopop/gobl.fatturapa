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

## Links

- [Manage the channel](https://sdi.fatturapa.gov.it/SdI2FatturaPAWebSpa/GestireCanaleAction.do)

## Pfx server certificate

```
$ openssl pkcs12 -export -certpbe PBE-SHA1-3DES -keypbe PBE-SHA1-3DES -nomac -out SDI-PIVA-SERVER.pfx -inkey key_server.key -in SDI-IT.INVOPOP.COM.pem -certfile ca-all.pem
Enter Export Password:
Verifying - Enter Export Password:
```
