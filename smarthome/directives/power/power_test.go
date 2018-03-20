package power_test

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/betom84/go-alexa/smarthome/common"
	"github.com/betom84/go-alexa/smarthome/directives/power"
	"github.com/betom84/go-alexa/smarthome/testdata/helpers"
	"github.com/betom84/go-alexa/smarthome/testdata/mocks"
)

var update = flag.Bool("update", false, "Run test and update golden file")

func TestController(t *testing.T) {
	power.Now = func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2018-02-23T22:57:05+00:00")
		return t
	}

	tt := []struct {
		name        string
		directive   *common.Directive
		device      interface{}
		expectError string
		goldenFile  string
	}{
		{
			name:       "it can turn on an power device",
			directive:  helpers.LoadRequest(t, "testdata/turnon_request.json"),
			device:     createMockPowerDevice(true, nil),
			goldenFile: "testdata/turnon_response.json",
		},
		{
			name:       "it can turn off an power device",
			directive:  helpers.LoadRequest(t, "testdata/turnoff_request.json"),
			device:     createMockPowerDevice(false, nil),
			goldenFile: "testdata/turnoff_response.json",
		},
		{
			name:        "it returns an error on incompatible device given",
			directive:   helpers.LoadRequest(t, "testdata/turnon_request.json"),
			device:      &mocks.MockDevice{},
			expectError: "endpoint device does not support change of powerState",
		},
		{
			name:        "it returns an error on incompatible directive name",
			directive:   helpers.CreateDirective(t, `{"header":{"namespace":"Alexa.PowerController", "name":"Explode"}}`),
			expectError: "directive name should be TurnOn or TurnOff",
		},
		{
			name:        "it returns an error on incompatible directive namespace",
			directive:   helpers.CreateDirective(t, `{"header":{"namespace":"Alexa.Whatever", "name":"TurnOn"}}`),
			expectError: "incompatible directive",
		},
		{
			name:        "it returns error when changing state fails",
			directive:   helpers.LoadRequest(t, "testdata/turnon_request.json"),
			device:      createMockPowerDevice(true, fmt.Errorf("something horrible happened")),
			expectError: "something horrible happened",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer helpers.FailOnPanic(t)

			resp, err := power.Controller{}.Process(tc.directive, tc.device)
			if err != nil || len(tc.expectError) > 0 {
				if err != nil && err.Error() != tc.expectError {
					t.Fatalf("unexpected error; %v", err)
				}

				if err == nil {
					t.Fatal("expected an error; but got none")
				}

				return
			}

			if len(tc.goldenFile) > 0 {
				if *update {
					helpers.UpdateGolden(t, tc.goldenFile, resp)
				}

				helpers.AssertEqualsGolden(t, tc.goldenFile, resp)
			}
		})
	}
}

func createMockPowerDevice(expectedState bool, returnError error) *mocks.MockPowerDevice {
	d := mocks.MockPowerDevice{}
	d.On("SetState", expectedState).Return(expectedState, returnError)
	return &d
}
