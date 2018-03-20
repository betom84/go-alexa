package smarthome_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/betom84/go-alexa/smarthome"
	"github.com/betom84/go-alexa/smarthome/common"
	"github.com/betom84/go-alexa/smarthome/testdata/mocks"
	"github.com/betom84/go-alexa/smarthome/validator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthorization(t *testing.T) {
	tt := []struct {
		name               string
		handlerUsername    string
		handlerPassword    string
		requestUsername    string
		requestPassword    string
		expectedStatusCode int
	}{
		{
			name:               "authorization will fail on wrong password",
			handlerUsername:    "tester",
			handlerPassword:    "right",
			requestUsername:    "tester",
			requestPassword:    "wrong",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "authorization will pass",
			handlerUsername:    "tester",
			handlerPassword:    "right",
			requestUsername:    "tester",
			requestPassword:    "right",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "authorization will be skipped when not configured",
			requestUsername:    "tester",
			requestPassword:    "right",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", bytes.NewReader([]byte(`{"directive":{"header":{"namespace":""}}}`)))
			req.SetBasicAuth(tc.requestUsername, tc.requestPassword)

			handler := smarthome.Handler{}
			handler.BasicAuth.Username = tc.handlerUsername
			handler.BasicAuth.Password = tc.handlerPassword

			handler.ServeHTTP(rec, req)
			defer rec.Result().Body.Close()

			responseBody, _ := ioutil.ReadAll(rec.Result().Body)
			assert.Equal(t, tc.expectedStatusCode, rec.Result().StatusCode, "got unexpected status code %s; %s", rec.Result().Status, string(responseBody))
		})
	}
}

func TestNewDefaultHandler(t *testing.T) {
	handler := smarthome.NewDefaultHandler(nil, nil)

	hv := reflect.ValueOf(*handler)
	assert.Equal(t, 4, hv.FieldByName("directiveProcessors").Len())
}

func TestHandler(t *testing.T) {
	common.ConstMessageID = "any-const-message-id-for-test"

	tt := []struct {
		name                   string
		request                []byte
		expectedStatusCode     int
		expectedResponse       []byte
		mockDeviceFactory      *mocks.MockDeviceFactory
		mockDirectiveProcessor *mocks.MockDirectiveProcessor
	}{
		{
			name:               "it can handle invalid request",
			request:            []byte(`something useless`),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   []byte(http.StatusText(http.StatusBadRequest)),
		},
		{
			name:               "it can handle invalid json request",
			request:            []byte(`{"something":"useless"}`),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   []byte(http.StatusText(http.StatusBadRequest)),
		},
		{
			name:               "it responds with an error on unprocessable directive",
			request:            []byte(`{"directive":{"header":{"namespace":"Something.Unprocessable"}}}`),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   readFile(t, "testdata/invalid_directive_response.json"),
		},
		{
			name:                   "it uses registered processor to handle request without endpoint",
			request:                []byte(`{"directive":{"header":{"namespace":"Registered.Processor"}}}`),
			expectedStatusCode:     http.StatusOK,
			mockDirectiveProcessor: createMockDirectiveProcessor(true, nil),
			expectedResponse:       []byte(`{"event":{"header":null}}`),
		},
		{
			name:                   "it uses registered processor to handle request with endpoint",
			request:                readFile(t, "testdata/process_directive_request.json"),
			expectedStatusCode:     http.StatusOK,
			mockDirectiveProcessor: createMockDirectiveProcessor(true, nil),
			mockDeviceFactory:      createMockDeviceFactory("testing", "ABC-123"),
			expectedResponse:       []byte(`{"event":{"header":null}}`),
		},
		{
			name:                   "it responds with alexaerror returned by processor",
			request:                readFile(t, "testdata/process_directive_request.json"),
			expectedStatusCode:     http.StatusOK,
			mockDirectiveProcessor: createMockDirectiveProcessor(true, common.NewInternalError("something horrible")),
			mockDeviceFactory:      createMockDeviceFactory("testing", "ABC-123"),
			expectedResponse:       readFile(t, "testdata/internal_error_response.json"),
		},
		{
			name:                   "it responds with transformed error returned by processor",
			request:                readFile(t, "testdata/process_directive_request.json"),
			expectedStatusCode:     http.StatusOK,
			mockDirectiveProcessor: createMockDirectiveProcessor(true, fmt.Errorf("something horrible")),
			mockDeviceFactory:      createMockDeviceFactory("testing", "ABC-123"),
			expectedResponse:       readFile(t, "testdata/internal_error_response.json"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", bytes.NewReader(tc.request))

			handler := smarthome.Handler{}
			handler.Validator = &validator.Validator{}

			if tc.mockDeviceFactory != nil {
				handler.DeviceFactory = tc.mockDeviceFactory
				defer tc.mockDeviceFactory.AssertExpectations(t)
			}

			if tc.mockDirectiveProcessor != nil {
				handler.AddDirectiveProcessor(tc.mockDirectiveProcessor)
				defer tc.mockDirectiveProcessor.AssertExpectations(t)
			}

			handler.ServeHTTP(rec, req)
			defer rec.Result().Body.Close()

			responseBody, err := ioutil.ReadAll(rec.Result().Body)
			if err != nil {
				t.Fatalf("could not read response; %v", err)
			}

			assert.Equal(t, tc.expectedStatusCode, rec.Result().StatusCode, "got unexpected status code %s; %s", rec.Result().Status, string(responseBody))

			if len(tc.expectedResponse) > 0 {
				if rec.Header().Get("content-type") == "application/json" {
					assert.JSONEq(t, string(tc.expectedResponse), string(responseBody), "received unexpected response body")
				} else {
					assert.Equal(t, string(tc.expectedResponse), string(responseBody), "received unexpected response body")
				}
			}
		})
	}
}

func createMockDirectiveProcessor(isCapable bool, errorOnProcess error) *mocks.MockDirectiveProcessor {
	p := mocks.MockDirectiveProcessor{}
	p.On("IsCapable", mock.Anything).Return(isCapable)
	p.On("Process", mock.Anything, mock.Anything).Return(&common.Response{}, errorOnProcess)
	return &p
}

func createMockDeviceFactory(expectedType string, expectedID string) *mocks.MockDeviceFactory {
	f := mocks.MockDeviceFactory{}
	f.On("NewDevice", expectedType, expectedID).Return("", nil)
	return &f
}

func readFile(t *testing.T, file string) []byte {
	t.Helper()

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("could not read file; %v", err)
	}

	return bytes
}
