package validator_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/betom84/go-alexa/smarthome/validator"
)

type Handler struct {
	t *testing.T
}

func (h Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	content, err := ioutil.ReadFile("testdata/schema.json")
	if err != nil {
		h.t.Fatalf("could not read json schema; %v", err)
	}
	writer.Write(content)
}

func TestValidator(t *testing.T) {
	h := Handler{t}
	srv := httptest.NewServer(h)

	tt := []struct {
		name      string
		validator validator.Validator
	}{
		{
			name:      "validate with a file schema reference",
			validator: validator.Validator{SchemaReference: "testdata/schema.json"},
		},
		{
			name:      "validate with a url schema reference",
			validator: validator.Validator{SchemaReference: srv.URL},
		},
		{
			name:      "validate with default url schema reference",
			validator: validator.Validator{SchemaReference: ""},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			testValidator(tc.validator, t)
		})
	}
}

func testValidator(v validator.Validator, t *testing.T) {

	tt := []struct {
		name  string
		value []byte
		err   string
	}{
		{
			name:  "validate a valid value",
			value: []byte(`{"event":{"header":{"namespace":"Alexa","name":"Response","messageId":"1","payloadVersion":"3"},"payload":{}}}`),
		},
		{
			name:  "validate a invalid value",
			value: []byte(`{"event":{"header":{"namespace":"Alexa","name":"Response","messageId":"1","payloadVersion":""},"payload":{}}}`),
			err:   "event.header.payloadVersion: event.header.payloadVersion must be one of the following: \"3\"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := v.Validate(tc.value)

			if err != nil && tc.err == "" {
				t.Fatalf("validate returned an error without expecting one; %v", err)
			}

			if err == nil && tc.err != "" {
				t.Fatal("validate returned no error while expecting one")
			}

			if err != nil && !strings.Contains(err.Error(), tc.err) {
				t.Fatalf("validate returned unexpected error; %v", err)
			}
		})
	}
}
