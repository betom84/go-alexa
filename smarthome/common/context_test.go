package common_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/betom84/go-alexa/smarthome/common"

	"github.com/stretchr/testify/assert"
)

type EndpointHealth struct {
	Connected bool
}

func (h EndpointHealth) IsConnected() bool {
	return h.Connected
}

func TestContextProperty(t *testing.T) {
	timeOfSample, _ := time.Parse(time.RFC3339, "2018-02-25T19:56:05+00:00")

	tt := []struct {
		name string
		add  func(*common.Context)
		json string
	}{
		{
			name: "add temperature property",
			add:  func(c *common.Context) { c.AddTemperatureProperty(20.28, timeOfSample) },
			json: `{"properties":[{"namespace":"Alexa.TemperatureSensor","name":"temperature","value":{"value":20.2,"scale":"CELSIUS"},"timeOfSample":"2018-02-25T19:56:05Z","uncertaintyInMilliseconds":100}]}`,
		},
		{
			name: "add temperature property with many decimal digits",
			add:  func(c *common.Context) { c.AddTemperatureProperty(22.43224575421, timeOfSample) },
			json: `{"properties":[{"namespace":"Alexa.TemperatureSensor","name":"temperature","value":{"value":22.4,"scale":"CELSIUS"},"timeOfSample":"2018-02-25T19:56:05Z","uncertaintyInMilliseconds":100}]}`,
		},
		{
			name: "add negative temperature property",
			add:  func(c *common.Context) { c.AddTemperatureProperty(-5.32, timeOfSample) },
			json: `{"properties":[{"namespace":"Alexa.TemperatureSensor","name":"temperature","value":{"value":-5.3,"scale":"CELSIUS"},"timeOfSample":"2018-02-25T19:56:05Z","uncertaintyInMilliseconds":100}]}`,
		},
		{
			name: "add turned on power state property",
			add:  func(c *common.Context) { c.AddPowerStateProperty(true, timeOfSample) },
			json: `{"properties":[{"namespace":"Alexa.PowerController","name":"powerState","value":"ON","timeOfSample":"2018-02-25T19:56:05Z","uncertaintyInMilliseconds":100}]}`,
		},
		{
			name: "add turned off power state property",
			add:  func(c *common.Context) { c.AddPowerStateProperty(false, timeOfSample) },
			json: `{"properties":[{"namespace":"Alexa.PowerController","name":"powerState","value":"OFF","timeOfSample":"2018-02-25T19:56:05Z","uncertaintyInMilliseconds":100}]}`,
		},
		{
			name: "add endpoint health property of connected endpoint",
			add:  func(c *common.Context) { c.AddEndpointHealthProperty(EndpointHealth{true}, timeOfSample) },
			json: `{"properties":[{"namespace":"Alexa.EndpointHealth","name":"connectivity","value":{"value":"OK"},"timeOfSample":"2018-02-25T19:56:05Z","uncertaintyInMilliseconds":100}]}`,
		},
		{
			name: "add endpoint health property of disconnected endpoint",
			add:  func(c *common.Context) { c.AddEndpointHealthProperty(EndpointHealth{false}, timeOfSample) },
			json: `{"properties":[{"namespace":"Alexa.EndpointHealth","name":"connectivity","value":{"value":"UNREACHABLE"},"timeOfSample":"2018-02-25T19:56:05Z","uncertaintyInMilliseconds":100}]}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := common.NewContext()

			tc.add(c)

			marshaledContext, err := json.Marshal(c)
			assert.NoError(t, err)

			t.Log(string(marshaledContext))

			assert.JSONEq(t, tc.json, string(marshaledContext))
		})
	}
}
