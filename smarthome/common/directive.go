package common

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// A Directive got sent by alexa
type Directive struct {
	Header   *Header                `json:"header"`
	Endpoint *Endpoint              `json:"endpoint,omitempty"`
	Payload  map[string]interface{} `json:"payload,omitempty"`
}

func (d Directive) String() string {
	var buffer bytes.Buffer

	if d.Header != nil {
		buffer.WriteString(d.Header.Namespace)
		buffer.WriteString(("."))
		buffer.WriteString(d.Header.Name)
	}

	if d.Endpoint != nil {
		buffer.WriteString(" (")
		buffer.WriteString(d.Endpoint.Cookie.Name)
		buffer.WriteString(")")
	}

	return buffer.String()
}

// NewDirective creates a new directive from json
func NewDirective(data []byte) (dir *Directive, err error) {
	dir = new(Directive)
	err = json.Unmarshal(data, dir)

	if err == nil && dir.Header == nil {
		err = fmt.Errorf("directive does not contain a header")
	}

	return
}
