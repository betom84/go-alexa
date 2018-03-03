package validator

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// ValidationError encapsulate gojsonschema result errors
type ValidationError struct {
	errors []gojsonschema.ResultError
}

func (e ValidationError) Error() (message string) {
	for _, desc := range e.errors {
		message = fmt.Sprintf("%s\n%s", message, desc)
	}
	return
}

// Validator can validate []byte response against his json schema
type Validator struct {

	// SchemaReference points to a file or an url
	SchemaReference string

	schemaBytes []byte
}

// Validate the given response against the json schema
func (v *Validator) Validate(resp []byte) error {

	if len(v.schemaBytes) == 0 {
		err := v.loadSchemaBytes()
		if err != nil {
			return err
		}
	}

	schemaLoader := gojsonschema.NewBytesLoader(v.schemaBytes)
	responseLoader := gojsonschema.NewBytesLoader(resp)

	result, err := gojsonschema.Validate(schemaLoader, responseLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		return ValidationError{result.Errors()}
	}

	return nil
}

func (v *Validator) loadSchemaBytes() error {
	if v.SchemaReference == "" {
		v.SchemaReference = "https://raw.githubusercontent.com/alexa/alexa-smarthome/master/validation_schemas/alexa_smart_home_message_schema.json"
	}

	var err error
	if strings.HasPrefix(v.SchemaReference, "http") {
		var r *http.Response
		r, err = http.Get(v.SchemaReference)
		if err != nil {
			return err
		}
		v.schemaBytes, err = ioutil.ReadAll(r.Body)
	} else {
		v.schemaBytes, err = ioutil.ReadFile(v.SchemaReference)
	}

	return err
}
