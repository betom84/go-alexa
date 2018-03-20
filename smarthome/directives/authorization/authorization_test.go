package authorization_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/betom84/go-alexa/smarthome/common"
	"github.com/betom84/go-alexa/smarthome/directives/authorization"
	"github.com/betom84/go-alexa/smarthome/testdata/helpers"
	"github.com/betom84/go-alexa/smarthome/testdata/mocks"

	"github.com/stretchr/testify/mock"
)

var update = flag.Bool("update", false, "Run test and update golden file")

func TestAuthorization(t *testing.T) {
	// never make real world requests in tests
	authorization.RequestTokenURL = ""
	authorization.RequestUserProfileURL = ""

	tt := []struct {
		name            string
		directive       *common.Directive
		expectError     string
		goldenFile      string
		mockHTTPHandler *mocks.MockHTTPHandler
		mockAuthority   *mocks.MockAuthority
	}{
		{
			name:        "it returns an error on incompatible directive namespace",
			directive:   helpers.CreateDirective(t, `{"header":{"namespace":"Alexa.Authorization.Not"}}`),
			expectError: "incompatible directive",
		},
		{
			name:        "it returns an error when authority is not set",
			directive:   helpers.CreateDirective(t, `{"header":{"namespace":"Alexa.Authorization"}}`),
			expectError: "authority is missing",
		},
		{
			name:            "it responds successful when authority returns no error",
			directive:       helpers.LoadRequest(t, "testdata/request.json"),
			goldenFile:      "testdata/response.json",
			mockHTTPHandler: createMockHTTPHandler(t, "testdata/profile.json", "testdata/tokens.json"),
			mockAuthority:   createMockAuthority(t, "mhashimoto-04@plaxo.com"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockHTTPHandler != nil {
				defer tc.mockHTTPHandler.AssertExpectations(t)

				srv := httptest.NewServer(tc.mockHTTPHandler)
				defer srv.Close()

				authorization.RequestTokenURL = fmt.Sprintf("%s/auth/o2/token", srv.URL)
				authorization.RequestUserProfileURL = fmt.Sprintf("%s/user/profile", srv.URL)
			}

			auth := authorization.Authorization{}
			if tc.mockAuthority != nil {
				auth.Authority = tc.mockAuthority
				defer tc.mockAuthority.AssertExpectations(t)
			}

			resp, err := auth.Process(tc.directive, nil)
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

func createMockHTTPHandler(t *testing.T, profileFile string, tokensFile string) *mocks.MockHTTPHandler {
	t.Helper()

	handler := mocks.MockHTTPHandler{T: t}

	if len(profileFile) > 0 {
		bytes, err := ioutil.ReadFile(profileFile)
		if err != nil {
			t.Fatalf("could not read profile; %v", err)
		}

		handler.On("RequestURI", "GET", "/user/profile", mock.Anything).Return(bytes, http.StatusOK)
	}

	if len(tokensFile) > 0 {
		bytes, err := ioutil.ReadFile(tokensFile)
		if err != nil {
			t.Fatalf("could not read tokens; %v", err)
		}

		handler.On("RequestURI", "POST", "/auth/o2/token", mock.Anything).Return(bytes, http.StatusOK)
	}

	return &handler
}

func createMockAuthority(t *testing.T, expectedEmail string) *mocks.MockAuthority {
	t.Helper()

	authority := mocks.MockAuthority{}
	authority.On("GetClientID").Return("clientID")
	authority.On("GetClientSecret").Return("clientSecret")
	authority.On("AcceptGrant", expectedEmail, mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(nil)

	return &authority
}
