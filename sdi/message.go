package sdi

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

// MessageHandler processes SOAP requests from SdI (Sistema di Interscambio)
func MessageHandler(handler HandleSOAPRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Printf("Failed to dump incoming request: %s\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("Incoming request:\n%s", requestDump)

		err = ParseMessage(req.Body, handler)
		if err != nil {
			log.Printf("Failed to parse body: %s\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		responseBody := []byte(soapEmptyResponse())
		response := &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusOK,
			Proto:         "HTTP/2.0",
			ProtoMajor:    2,
			ProtoMinor:    0,
			Body:          io.NopCloser(bytes.NewReader(responseBody)),
			ContentLength: int64(len(responseBody)),
			Header:        make(http.Header),
		}
		// contentType := http.DetectContentType(responseBody)
		response.Header.Set("Content-Type", "application/soap+xml")

		responseDump, err := httputil.DumpResponse(response, true)
		if err != nil {
			log.Printf("Failed to dump outgoing response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("Outgoing response:\n%s", responseDump)

		err = responseToWriter(w, response)
		if err != nil {
			log.Printf("Failed to send response: %v", err)
			return
		}
	}
}

func responseToWriter(w http.ResponseWriter, response *http.Response) error {
	w.WriteHeader(response.StatusCode)

	for k, v := range response.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	_, ok := io.Copy(w, response.Body)
	return ok
}

func soapEmptyResponse() string {
	return `<?xml version='1.0' encoding='UTF-8'?>` +
		`<soapenv:Envelope xmlns:soapenv='http://schemas.xmlsoap.org/soap/envelope/' xmlns:typ='http://www.fatturapa.gov.it/sdi/ws/trasmissione/v1.0/types'>` +
		`<soapenv:Header/>` +
		`<soapenv:Body/>` +
		`</soapenv:Envelope>`
}
