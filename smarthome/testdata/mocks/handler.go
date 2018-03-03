package mocks

import (
	"github.com/betom84/go-alexa/smarthome/common"

	"github.com/stretchr/testify/mock"
)

// MockDeviceFactory ...
type MockDeviceFactory struct {
	mock.Mock
}

// NewDevice ...
func (f *MockDeviceFactory) NewDevice(epType string, id string) (interface{}, error) {
	r := f.Called(epType, id)
	return r.Get(0), r.Error(1)
}

// MockDirectiveProcessor ...
type MockDirectiveProcessor struct {
	mock.Mock
}

// IsCapable ...
func (p *MockDirectiveProcessor) IsCapable(dir *common.Directive) bool {
	r := p.Called(dir)
	return r.Bool(0)
}

// Process ...
func (p *MockDirectiveProcessor) Process(dir *common.Directive, device interface{}) (*common.Response, error) {
	r := p.Called(dir, device)
	return r.Get(0).(*common.Response), r.Error(1)
}
