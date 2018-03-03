package alexa_test

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/betom84/go-alexa/smarthome/common"
	"github.com/betom84/go-alexa/smarthome/directives/alexa"
	"github.com/betom84/go-alexa/smarthome/testdata/helpers"
	"github.com/betom84/go-alexa/smarthome/testdata/mocks"
)

var update = flag.Bool("update", false, "Run test and update golden file")

func TestReport(t *testing.T) {

	alexa.Now = func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2018-02-24T16:42:05+00:00")
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
			name:        "it returns an error on incompatible directive namespace",
			directive:   helpers.CreateDirective(t, `{"header":{"namespace":"NotAlexa"}}`),
			expectError: "incompatible directive",
		},
		{
			name:       "it can handle unspecific devices without reporting anything",
			directive:  helpers.LoadRequest(t, "testdata/request.json"),
			device:     mocks.MockDevice{},
			goldenFile: "testdata/empty_response.json",
		},
		{
			name:       "it can report power device state",
			directive:  helpers.LoadRequest(t, "testdata/request.json"),
			device:     createMockPowerDevice(true, nil),
			goldenFile: "testdata/powerstate_response.json",
		},
		{
			name:        "it returns error from power device",
			directive:   helpers.LoadRequest(t, "testdata/request.json"),
			device:      createMockPowerDevice(false, fmt.Errorf("the end is near")),
			expectError: "the end is near",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := alexa.Report{}.Process(tc.directive, tc.device)
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

func createMockPowerDevice(returnState bool, returnError error) *mocks.MockPowerDevice {
	d := mocks.MockPowerDevice{}
	d.On("State").Return(returnState, returnError)
	return &d
}
