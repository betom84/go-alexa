package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockAuthority ...
type MockAuthority struct {
	mock.Mock
}

// GetClientID ...
func (a *MockAuthority) GetClientID() string {
	r := a.Called()
	return r.String(0)
}

// GetClientSecret ...
func (a *MockAuthority) GetClientSecret() string {
	r := a.Called()
	return r.String(0)
}

// AcceptGrant ...
func (a *MockAuthority) AcceptGrant(email string, bearerToken string, accessTokens map[string]interface{}) error {
	r := a.Called(email, bearerToken, accessTokens)
	return r.Error(0)
}
