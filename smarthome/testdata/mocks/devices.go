package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockDevice ...
type MockDevice struct {
	mock.Mock
}

// MockPowerDevice ...
type MockPowerDevice struct {
	MockDevice
}

// SetState satisfies powerState capability
func (t *MockPowerDevice) SetState(value bool) (bool, error) {
	r := t.Called(value)
	return r.Bool(0), r.Error(1)
}

// State satisfies powerState capability
func (t *MockPowerDevice) State() (bool, error) {
	r := t.Called()
	return r.Bool(0), r.Error(1)
}
