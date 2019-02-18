package discovery

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/betom84/go-alexa/smarthome/common"
	"github.com/betom84/go-alexa/smarthome/testdata/helpers"
)

var update = flag.Bool("update", false, "Run test and update golden file")

func TestDiscovery(t *testing.T) {
	tt := []struct {
		name        string
		processor   Discovery
		directive   *common.Directive
		goldenFile  string
		expectError string
	}{
		{
			name:       "it responds with endpoints from file",
			processor:  createDiscovery(t, "testdata/endpoints.json"),
			directive:  helpers.LoadRequest(t, "testdata/request.json"),
			goldenFile: "testdata/response.json",
		},
		{
			name:        "it returns an error on incompatible directive namespace",
			processor:   Discovery{},
			directive:   helpers.CreateDirective(t, `{"header":{"namespace":"Not.Alexa.Discovery"}}`),
			expectError: "incompatible directive",
		},
		{
			name:        "it returns an error when endpoints are undefined",
			processor:   Discovery{},
			directive:   helpers.LoadRequest(t, "testdata/request.json"),
			expectError: "endpoints not specified",
		},
		{
			name:        "it returns an error when endpoints could not be unmarshaled",
			processor:   createDiscovery(t, "testdata/endpoints.invalid"),
			directive:   helpers.LoadRequest(t, "testdata/request.json"),
			expectError: "could not unmarshal endpoints; invalid character 'o' in literal null (expecting 'u')",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := tc.processor.Process(tc.directive, nil)
			if err != nil && err.Error() != tc.expectError {
				t.Fatalf("unexpected error; %v", err)
			}

			if err == nil && len(tc.expectError) > 0 {
				t.Fatal("expected an error; but got none")
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

func createDiscovery(t *testing.T, endpoints string) Discovery {
	ep, err := os.Open(endpoints)
	if err != nil {
		t.Fatal(err)
	}
	return Discovery{Endpoints: ep}
}

func TestDiscoveryUnmarshalEndpointsOnDemand(t *testing.T) {
	ep, err := os.Open("testdata/endpoints.json")
	if err != nil {
		t.Fatal(err)
	}

	discovery := Discovery{Endpoints: ep}
	assert.Len(t, discovery.unmarshaledEndpoints, 0)

	err = discovery.unmarshalEndpointsOnDemand()
	assert.NoError(t, err)
	assert.Len(t, discovery.unmarshaledEndpoints, 1)

	err = ep.Close()
	assert.Error(t, err, "file already closed")

	err = discovery.unmarshalEndpointsOnDemand()
	assert.NoError(t, err)
	assert.Len(t, discovery.unmarshaledEndpoints, 1)
}
