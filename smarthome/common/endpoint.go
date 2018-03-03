package common

// An Endpoint object identifies the target for a directive and the origin of an event.
type Endpoint struct {
	Scope struct {
		Type  string `json:"type"`
		Token string `json:"token"`
	} `json:"scope"`
	EndpointID string `json:"endpointId"`
	Cookie     struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"cookie"`
}
