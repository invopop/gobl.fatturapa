# Certificates

[Source link](https://www.fatturapa.gov.it/it/norme-e-regole/DocumentazioneSDI/)

## Test Certificates

Contents of the `test` folder:

- `testservizi.fatturapa.it.cer`: SERVER certificate exposed by
   the test services of SdI (Sistema di Interscambio)
- `SistemaInterscambioFatturaPATest.cer`: Public part of the CLIENT certificate
   used by SdI (Sistema di Interscambio) to invoke the test services exposed by you

## Production Certificates

Contents of the `production` folder:

- `servizi.fatturapa.it.cer`: SERVER certificate exposed by
  the services of SdI (Sistema di Interscambio)
- `Sistema_Interscambio_Fattura_PA.cer`: Public part of the CLIENT certificate
  used by SdI (Sistema di Interscambio) to invoke the services exposed by you

## CA Certificates

Contents of the `ca` folder:

- `caentrate.cer`: CA certificate for the production environment
- `CAEntratetest.cer`: CA certificate to validate the test SdI certificate
- `Sectigo RSA.cer`: CA certificate for servizi.fatturapa.it.cer
  for the production environment
- `UserTrustCA.cer`: CA certificate for servizi.fatturapa.it.cer
  for the production environment
