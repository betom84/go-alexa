package smarthome

import (
	"errors"
	"log"
)

// Authority to handle alexa user authorization
type Authority struct {

	// Amazon client ID used for oauth authorization
	ClientID string

	// Amazon client secret used for oauth authorization
	ClientSecret string

	// E-Mail addresses of users granted access
	RestrictedUsers []string
}

// AcceptGrant is used to grant access to an alexa user and store the according access tokens
func (a Authority) AcceptGrant(email string, bearerToken string, accessTokens map[string]interface{}) error {

	var granted = (len(a.RestrictedUsers) == 0)
	for _, restricted := range a.RestrictedUsers {
		granted = (restricted == email)
		if granted {
			break
		}
	}

	if !granted {
		return errors.New("Restricted users only")
	}

	// todo, hold the token for async responses and refresh after the given time
	//log.Printf("Granted access to %s.\nAccess tokens:%s", email, accessTokens)
	log.Printf("Granted access to %s.", email)

	return nil
}

// GetClientID gets the ClientID property, this function is needed to satisfy interfaces
func (a Authority) GetClientID() string {
	return a.ClientID
}

// GetClientSecret gets the ClientSecret property, this function is needed to satisfy interfaces
func (a Authority) GetClientSecret() string {
	return a.ClientSecret
}
