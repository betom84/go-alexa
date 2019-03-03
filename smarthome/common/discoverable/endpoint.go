package discoverable

import "github.com/betom84/go-alexa/smarthome/common"

// DisplayCategory is used to categorize an endpoint
type DisplayCategory string

//
const (
	Light            DisplayCategory = "LIGHT"
	Switch           DisplayCategory = "SWITCH"
	Other            DisplayCategory = "OTHER"
	TemperaturSensor DisplayCategory = "TEMPERATURE_SENSOR"
)

// Endpoint describes a discoverable device with multiple capabilities gets controlled with Alexa
type Endpoint struct {
	EndpointID        string            `json:"endpointId"`
	FriendlyName      string            `json:"friendlyName"`
	Description       string            `json:"description"`
	ManufacturerName  string            `json:"manufacturerName"`
	DisplayCategories []DisplayCategory `json:"displayCategories"`
	Cookie            common.Cookie     `json:"cookie"`
	Capabilities      []Capability      `json:"capabilities"`
}
