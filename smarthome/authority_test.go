package smarthome_test

import (
	"testing"

	"github.com/betom84/go-alexa/smarthome"

	"github.com/stretchr/testify/assert"
)

func TestAuthority(t *testing.T) {
	authority := smarthome.Authority{
		ClientID:        "clientID",
		ClientSecret:    "clientSecret",
		RestrictedUsers: []string{"somebody@mail.com"},
	}

	assert.Equal(t, "clientID", authority.GetClientID())
	assert.Equal(t, "clientSecret", authority.GetClientSecret())

	assert.Nil(t, authority.AcceptGrant("somebody@mail.com", "", nil))
	assert.Errorf(t, authority.AcceptGrant("anybody@mail.com", "", nil), "Restricted users only")
}
