// Package discovery contains the directive processor to handle directives with namepace "Alexa.Discovery"
package discovery

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/betom84/go-alexa/smarthome/common"
)

// Discovery uses the Endpoints file to process discovery directive
type Discovery struct {
	Endpoints io.Reader
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

	ep, err := ioutil.ReadAll(d.Endpoints)
	if err != nil {
		return nil, fmt.Errorf("could not read endpoints; %s", err)
	}

	var endpoints []interface{}
	err = json.Unmarshal(ep, &endpoints)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal endpoints; %s", err)
	}

	resp := new(common.Response)
	resp.Event.Header = common.NewHeader("Discover.Response", "Alexa.Discovery")

	resp.Event.Payload = struct {
		Endpoints []interface{} `json:"endpoints"`
	}{
		Endpoints: endpoints}

	return resp, nil
}
