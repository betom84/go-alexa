package alexa

import (
	"fmt"
	"time"

	"github.com/betom84/go-alexa/smarthome/common"
	"github.com/betom84/go-alexa/smarthome/common/capabilities"
)

// Now is used to change the current time for tets, defaults to time.Now()
var Now = time.Now

// Report processes the alexa directive to report state of an device
type Report struct{}

// IsCapable checks if an common.Directive can be processed
func (r Report) IsCapable(dir *common.Directive) bool {
	return dir.Header.Namespace == "Alexa"
}

// Process reports the current state of a device in response of an alexa directive
func (r Report) Process(dir *common.Directive, ed interface{}) (*common.Response, error) {
	if !r.IsCapable(dir) {
		return nil, fmt.Errorf("incompatible directive")
	}

	var resp = r.createResponse(dir)
	resp.Context = common.NewContext()

	switch rd := ed.(type) {
	case capabilities.PowerDevice:
		state, err := rd.State()
		if err != nil {
			return nil, err
		}
		resp.Context.AddPowerStateProperty(state, Now())
	case capabilities.TemperatureSensor:
		resp.Context.AddTemperatureProperty(rd.Temperature(), Now())
	}

	if h, ok := ed.(capabilities.HealthConscious); ok == true {
		resp.Context.AddEndpointHealthProperty(h, Now())
	}

	return resp, nil
}

func (r Report) createResponse(dir *common.Directive) (resp *common.Response) {
	resp = new(common.Response)

	resp.Event.Header = common.NewHeader("StateReport", "Alexa")
	resp.Event.Header.CorrelationToken = dir.Header.CorrelationToken
	resp.Event.Endpoint = dir.Endpoint
	resp.Event.Payload = struct{}{}

	return
}
