package power

import (
	"fmt"
	"github.com/betom84/go-alexa/smarthome/common"
	"time"
)

// Now is used to change the current time for tets, defaults to time.Now()
var Now = time.Now

// Controller process TurnOn or TurnOff power directives to control devices
type Controller struct {
}

// IsCapable checks if an common.Directive is a powercontroller directive
func (c Controller) IsCapable(dir *common.Directive) bool {
	return dir.Header.Namespace == "Alexa.PowerController"
}

// Process change the current state of the given endpoint
func (c Controller) Process(dir *common.Directive, ed interface{}) (*common.Response, error) {
	if !c.IsCapable(dir) {
		return nil, fmt.Errorf("incompatible directive")
	}

	var value bool
	switch dir.Header.Name {
	case "TurnOn":
		value = true
	case "TurnOff":
		value = false
	default:
		return nil, common.NewInvalidDirectiveError("directive name should be TurnOn or TurnOff")
	}

	pd, ok := ed.(common.PowerDevice)
	if !ok {
		return nil, fmt.Errorf("endpoint device does not support change of powerState")
	}

	state, err := pd.SetState(value)
	if err != nil {
		return nil, err
	}

	resp := new(common.Response)
	resp.Event.Header = common.NewHeader("Response", "Alexa")
	resp.Event.Header.CorrelationToken = dir.Header.CorrelationToken
	resp.Event.Endpoint = dir.Endpoint
	resp.Event.Payload = struct{}{}
	resp.Context = common.NewContext()
	resp.Context.AddPowerStateProperty(state, Now())

	return resp, nil
}
