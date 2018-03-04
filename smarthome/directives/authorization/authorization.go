// Package authorization contains the directive processor to handle directives with namepace "Alexa.Authorization"
package authorization

import (
	"encoding/json"
	"fmt"
	"github.com/betom84/go-alexa/smarthome/common"
	"io/ioutil"
	"net/http"
	"net/url"
)

// RequestTokenURL to change for tests
var RequestTokenURL = "https://api.amazon.com/auth/o2/token"

// RequestUserProfileURL to change for tests
var RequestUserProfileURL = "https://api.amazon.com/user/profile"

// Authority represents the instance where an alexa user gets access granted
type Authority interface {
	GetClientID() string
	GetClientSecret() string
	AcceptGrant(email string, bearerToken string, accessTokens map[string]interface{}) error
}

// Authorization checks the grantee and requests the access_token (along with renewal information)
// for async directive responses
type Authorization struct {
	Authority Authority
}

// IsCapable checks if an common.Directive is an authorization directive
func (a Authorization) IsCapable(dir *common.Directive) bool {
	return dir.Header.Namespace == "Alexa.Authorization"
}

// Process the authorization directive, device should be nil
func (a Authorization) Process(dir *common.Directive, device interface{}) (*common.Response, error) {

	if !a.IsCapable(dir) {
		return nil, fmt.Errorf("incompatible directive")
	}

	if a.Authority == nil {
		return nil, fmt.Errorf("authority is missing")
	}

	grantee := dir.Payload["grantee"].(map[string]interface{})
	profile, err := a.retrieveGranteeProfile(grantee["token"].(string))
	if err != nil {
		return nil, err
	}

	grant := dir.Payload["grant"].(map[string]interface{})
	tokens, err := a.retrieveAccessTokens(grant["code"].(string))
	if err != nil {
		return nil, err
	}

	err = a.Authority.AcceptGrant(profile["email"].(string), grantee["token"].(string), tokens)
	if err != nil {
		return nil, common.NewAcceptGrantFailedError(err.Error())
	}

	return a.createResponse(), nil
}

func (a Authorization) retrieveGranteeProfile(granteeToken string) (map[string]interface{}, error) {
	request, err := http.NewRequest("GET", RequestUserProfileURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", granteeToken))

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var granteeProfile map[string]interface{}
	err = json.Unmarshal(body, &granteeProfile)

	return granteeProfile, err
}

func (a Authorization) retrieveAccessTokens(code string) (map[string]interface{}, error) {
	resp, err := http.PostForm(RequestTokenURL, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {a.Authority.GetClientID()},
		"client_secret": {a.Authority.GetClientSecret()}})

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokens map[string]interface{}
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (a Authorization) createResponse() *common.Response {
	r := new(common.Response)
	r.Event.Header = common.NewHeader("AcceptGrant.Response", "Alexa.Authorization")
	r.Event.Payload = struct{}{}

	return r
}
