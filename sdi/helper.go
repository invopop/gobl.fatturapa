package sdi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"mime"
	"mime/multipart"
	"strings"

	resty "github.com/go-resty/resty/v2"
)

// parseMultipartResponse parses a multipart HTTP response and deserializes the content into the provided structure
func parseMultipartResponse(resp *resty.Response, response interface{}) error {
	mediaType, params, err := mime.ParseMediaType(resp.Header().Get("Content-Type"))
	if err != nil {
		return err
	}

	if !strings.HasPrefix(mediaType, "multipart/related") {
		return fmt.Errorf("unexpected content type: %s", mediaType)
	}

	reader := strings.NewReader(string(resp.Body()))
	mr := multipart.NewReader(reader, params["boundary"])

	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}
		if part.Header.Get("Content-Type") != "application/xop+xml; charset=utf-8; type=\"text/xml\"" {
			continue
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(part)
		if err != nil {
			return fmt.Errorf("multipart reading error: %s", err)
		}
		xmlData := buf.String()

		err = xml.Unmarshal([]byte(xmlData), &response)
		if err != nil {
			return fmt.Errorf("parsing xml error: %s", err)
		}
	}
	return nil
}
