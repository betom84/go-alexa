package discoverable

// A Capability describes the functionality an endpoint is capable of for Alexa.
// Capabilities are corresponding with properties from context used to respond to an Alexa directive.
type Capability struct {
	Type       string     `json:"type"`
	Interface  string     `json:"interface"`
	Version    string     `json:"version"`
	Properties Properties `json:"properties"`
}

// Properties of an Capability
type Properties struct {
	Supported           []Supported `json:"supported"`
	ProactivelyReported bool        `json:"proactivelyReported"`
	Retrievable         bool        `json:"retrievable"`
}

// Supported property names
type Supported struct {
	Name string `json:"name"`
}

// NewCapability to create Capability with default values
func NewCapability(interfacE string, supportedPropertyNames []string) Capability {
	supported := []Supported{}

	for _, propName := range supportedPropertyNames {
		supported = append(supported, Supported{
			Name: propName,
		})
	}

	return Capability{
		Type:      "AlexaInterface",
		Interface: interfacE,
		Version:   "3",
		Properties: Properties{
			Supported:           supported,
			ProactivelyReported: false,
			Retrievable:         true,
		},
	}
}
