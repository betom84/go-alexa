// Package discovery contains the directive processor to handle directives with namepace "Alexa.Discovery"
package discovery

import (
	"fmt"

	"github.com/betom84/go-alexa/smarthome/common"
	"github.com/betom84/go-alexa/smarthome/common/discoverable"
)

// Discovery uses the Endpoints file to process discovery directive
type Discovery struct {
	Endpoints            []discoverable.Endpoint
	unmarshaledEndpoints []interface{}
}

// IsCapable checks if an common.Directive is an discovery directive
func (d Discovery) IsCapable(dir *common.Directive) bool {
	return dir.Header.Namespace == "Alexa.Discovery"
}

// Process the discovery directive, device should be nil
func (d Discovery) Process(dir *common.Directive, device interface{}) (*common.Response, error) {
	if !d.IsCapable(dir) {
		return nil, fmt.Errorf("incompatible directive")
	}

	if d.Endpoints == nil {
		return nil, fmt.Errorf("endpoints not specified")
	}

	resp := new(common.Response)
	resp.Event.Header = common.NewHeader("Discover.Response", "Alexa.Discovery")

	resp.Event.Payload = struct {
		Endpoints []discoverable.Endpoint `json:"endpoints"`
	}{
		Endpoints: d.Endpoints,
	}

	return resp, nil
}
