package mocks

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockHTTPHandler ...
type MockHTTPHandler struct {
	T *testing.T
	mock.Mock
}

// ServeHTTP ...
func (h *MockHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.T.Fatalf("could not read request body; %v", err)
	}
	defer r.Body.Close()

	resp, status := h.RequestURI(r.Method, r.RequestURI, body)
	w.Write(resp)
	w.WriteHeader(status)
}

// RequestURI is used to expect requests made to the mocked handler
func (h *MockHTTPHandler) RequestURI(method string, uri string, body []byte) ([]byte, int) {
	r := h.Called(method, uri, body)
	return r.Get(0).([]byte), r.Int(1)
}
