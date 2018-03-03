package directives_test

import (
	"testing"

	"github.com/betom84/go-alexa/smarthome/directives"
	"github.com/betom84/go-alexa/smarthome/testdata/mocks"

	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	assert.NotNil(t, directives.CreateAuthorizeDirectiveProcessor(&mocks.MockAuthority{}))
	assert.NotNil(t, directives.CreateDiscoveryDirectiveProcessor(nil))
	assert.NotNil(t, directives.CreatePowerControllerDirectiveProcessor())
	assert.NotNil(t, directives.CreateReportAlexaDirectiveProcessor())
}
