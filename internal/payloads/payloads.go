package payloads

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"html/template"
)

// Payload is a generic map[string]interface{} so we can unmarshal all sorts
// of hooks into our struct without issue
type Payload map[string]interface{}

// DecodeReceivedPayload decodes a response body to a payload object
func DecodeReceivedPayload(r *http.Request, payload *Payload) error {
	switch r.Header.Get("Content-type") {
	case "application/json":
		return json.NewDecoder(r.Body).Decode(payload)

	case "application/x-www-form-urlencoded":
		bufferBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}

		parsed, err := url.ParseQuery(string(bufferBody))
		if err != nil {
			return err
		}

		for k, v := range parsed {
			(*payload)[k] = v
		}

	default:
		return errors.New("Invalid or missing Content-Type header")
	}

	return nil
}

// ParseTargetPayloadFormat converts the target payload format to a json payload
func ParseTargetPayloadFormat(targetFormat interface{}) ([]byte, error) {
	targetPayload, err := json.Marshal(targetFormat)
	return targetPayload, err
}

// CreateTargetPayload will generate the target paylaod
func CreateTargetPayload(format string, payload interface{}) (bytes.Buffer, error) {
	var output bytes.Buffer

	tmpl, err := template.New("").Funcs(transformTemplates).Parse(format)
	if err != nil {
		return output, err
	}

	if err := tmpl.Execute(&output, payload); err != nil {
		return output, err
	}

	return output, nil
}
