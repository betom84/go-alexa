package common

import (
	"github.com/betom84/go-alexa/smarthome/common/capabilities"
	"time"
)

// Context holds properties of an endpoint
type Context struct {
	Properties []property `json:"properties"`
}

// Property represents some state of a device.
type property struct {
	Namespace                 string      `json:"namespace"`
	Name                      string      `json:"name"`
	Value                     interface{} `json:"value"`
	TimeOfSample              time.Time   `json:"timeOfSample"`
	UncertaintyInMilliseconds int         `json:"uncertaintyInMilliseconds"`
}

func (c *Context) addProperty(p property) {
	c.Properties = append(c.Properties, p)
}

// AddPowerStateProperty adds a powerState property to context
func (c *Context) AddPowerStateProperty(state bool, timeOfSample time.Time) {

	value := "OFF"
	if state == true {
		value = "ON"
	}

	c.addProperty(property{
		Namespace:                 "Alexa.PowerController",
		Name:                      "powerState",
		Value:                     value,
		TimeOfSample:              timeOfSample,
		UncertaintyInMilliseconds: 100})
}

// AddTemperatureProperty adds a temperature property to context
func (c *Context) AddTemperatureProperty(temperature float32, timeOfSample time.Time) {

	iv := int(temperature * 10)
	v := float32(iv) / 10

	c.addProperty(property{
		Namespace: "Alexa.TemperatureSensor",
		Name:      "temperature",
		Value: struct {
			Value float32 `json:"value"`
			Scale string  `json:"scale"`
		}{
			Value: v,
			Scale: "CELSIUS",
		},
		TimeOfSample:              timeOfSample,
		UncertaintyInMilliseconds: 100})
}

// AddEndpointHealthProperty adds a connectivity property to context
func (c *Context) AddEndpointHealthProperty(health capabilities.HealthConscious, timeOfSample time.Time) {

	var value = "OK"
	if !health.IsConnected() {
		value = "UNREACHABLE"
	}

	c.addProperty(property{
		Namespace: "Alexa.EndpointHealth",
		Name:      "connectivity",
		Value: struct {
			Value string `json:"value"`
		}{
			Value: value,
		},
		TimeOfSample:              timeOfSample,
		UncertaintyInMilliseconds: 100})
}

// NewContext creates a new empty Context
func NewContext() *Context {
	return new(Context)
}
