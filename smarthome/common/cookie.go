package common

// The Cookie represents custom properties of an alexa endpoint and
// is used to create an device with the DeviceFactory which can be processed by the
// according directive processor.
type Cookie struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}
