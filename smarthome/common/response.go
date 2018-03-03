package common

// Response is used to respond to a request send by alexa
type Response struct {
	Context *Context `json:"context,omitempty"`
	Event   struct {
		Header   *Header     `json:"header"`
		Endpoint *Endpoint   `json:"endpoint,omitempty"`
		Payload  interface{} `json:"payload,omitempty"`
	} `json:"event"`
}
