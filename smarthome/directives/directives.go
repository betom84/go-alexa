// Package directives contains the factory to create specific directive processors
package directives

import (
	"io"

	"github.com/betom84/go-alexa/smarthome/common"
	"github.com/betom84/go-alexa/smarthome/directives/alexa"
	"github.com/betom84/go-alexa/smarthome/directives/authorization"
	"github.com/betom84/go-alexa/smarthome/directives/discovery"
	"github.com/betom84/go-alexa/smarthome/directives/power"
)

// DirectiveProcessor describes something which can process an alexa directive
type DirectiveProcessor interface {
	// Process the directive for an optionally device
	Process(*common.Directive, interface{}) (*common.Response, error)

	// IsCapable checks if an common.Directive can be processed
	IsCapable(*common.Directive) bool
}

// CreateAuthorizeDirectiveProcessor returns a DirectiveProcessor to process authorization directives
func CreateAuthorizeDirectiveProcessor(authority authorization.Authority) DirectiveProcessor {
	return authorization.Authorization{Authority: authority}
}

// CreateDiscoveryDirectiveProcessor returns a DirectiveProcessor to process discovery directives
func CreateDiscoveryDirectiveProcessor(endpoints io.ReadCloser) DirectiveProcessor {
	return &discovery.Discovery{Endpoints: endpoints}
}

// CreatePowerControllerDirectiveProcessor returns a DirectiveProcessor to process power controller directives
func CreatePowerControllerDirectiveProcessor() DirectiveProcessor {
	return power.Controller{}
}

// CreateReportAlexaDirectiveProcessor returns a DirectiveProcessor to process report directives
func CreateReportAlexaDirectiveProcessor() DirectiveProcessor {
	return alexa.Report{}
}
