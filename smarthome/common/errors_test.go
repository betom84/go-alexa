package common_test

import (
	"github.com/betom84/go-alexa/smarthome/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {

	tt := []struct {
		name    string
		err     common.AlexaError
		errMsg  string
		errType string
		errNS   string
	}{
		{
			name:    "it creates 'accept grant failed' error",
			err:     common.NewAcceptGrantFailedError("message for test"),
			errMsg:  "message for test",
			errType: "ACCEPT_GRANT_FAILED",
			errNS:   "Alexa.Authorization",
		},
		{
			name:    "it creates 'internal' error",
			err:     common.NewInternalError("message for test"),
			errMsg:  "message for test",
			errType: "INTERNAL_ERROR",
			errNS:   "Alexa",
		},
		{
			name:    "it creates 'invalid directive' error",
			err:     common.NewInvalidDirectiveError("message for test"),
			errMsg:  "message for test",
			errType: "INVALID_DIRECTIVE",
			errNS:   "Alexa",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.err.Error(), tc.errMsg)
			assert.Equal(t, tc.err.Type, tc.errType)
			assert.Equal(t, tc.err.Namespace, tc.errNS)
		})
	}
}
