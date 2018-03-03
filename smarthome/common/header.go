package common

import (
	"log"

	"github.com/nu7hatch/gouuid"
)

// A Header has a set of expected fields that are the same across message types.
// These provide different types of identifying information.
type Header struct {
	Namespace        string `json:"namespace"`
	Name             string `json:"name"`
	MessageID        string `json:"messageId"`
	CorrelationToken string `json:"correlationToken,omitempty"`
	PayloadVersion   string `json:"payloadVersion"`
}

// NewHeader creates a new header for given name and namespace
func NewHeader(name string, namespace string) *Header {
	return &Header{
		Namespace:      namespace,
		Name:           name,
		PayloadVersion: "3",
		MessageID:      createUUID()}
}

// ConstMessageID is used to ensure a constant messageId for tests
var ConstMessageID = ""

func createUUID() string {
	if len(ConstMessageID) > 0 {
		return ConstMessageID
	}

	res, err := uuid.NewV4()
	if err != nil {
		log.Fatal("Could not create uuid: ", err)
		panic(err)
	}
	return res.String()
}
