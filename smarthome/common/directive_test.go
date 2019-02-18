package common_test

import (
	"testing"

	"github.com/betom84/go-alexa/smarthome/common"

	"github.com/stretchr/testify/assert"
)

func TestDirective(t *testing.T) {
	tt := []struct {
		name     string
		payload  []byte
		err      string
		toString string
	}{
		{
			name:     "it creates an directive from payload",
			payload:  []byte(`{"header":{"namespace":"Namespace","name":"Name"}}`),
			toString: "Namespace.Name",
		},
		{
			name:     "it creates an directive with endpoint from payload",
			payload:  []byte(`{"header":{"namespace":"Namespace","name":"Name"},"endpoint":{"cookie":{"name":"Awesome Endpoint","type":"Type","id":"ID"}}}`),
			toString: "Namespace.Name (Awesome Endpoint/Type/ID)",
		},
		{
			name:    "it returns an error on empty payload",
			payload: []byte{},
			err:     "unexpected end of JSON input",
		},
		{
			name:    "it returns an error on invalid payload",
			payload: []byte("invalid payload"),
			err:     "invalid character 'i' looking for beginning of value",
		},
		{
			name:    "it return an error on empty json payload",
			payload: []byte(`{}`),
			err:     "directive does not contain a header",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir, err := common.NewDirective(tc.payload)
			if err != nil {
				assert.EqualError(t, err, tc.err)
			}

			if err == nil && len(tc.err) != 0 {
				t.Fatal("expected error not occurred")
			}

			assert.NotNil(t, dir)
			assert.Equal(t, tc.toString, dir.String())
		})
	}
}
