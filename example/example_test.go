package example

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/betom84/go-alexa/smarthome"
	"github.com/betom84/go-alexa/smarthome/directives"
	"github.com/betom84/go-alexa/smarthome/validator"
)

type MyDeviceFactory struct{}

func (f *MyDeviceFactory) NewDevice(epType string, id string) (interface{}, error) {
	return &MyDevice{epType, id}, nil
}

type MyDevice struct {
	Type string
	ID   string
}

func (d *MyDevice) SetState(state bool) (bool, error) {
	fmt.Printf("%s device %s state changed to %t", d.Type, d.ID, state)
	return state, nil
}

func (d *MyDevice) State() (bool, error) {
	return true, nil
}

func Example() {
	// First we need to create an alexa handler.
	handler := smarthome.Handler{

		// The device factory is needed to transform an alexa endpoint into a device with
		// the needed capability. In this example the devices created by the factory will
		// have the powerState capability.
		DeviceFactory: &MyDeviceFactory{},

		// Optionally we set a validator instance. Because alexa wont work with non well-formed
		// responses, this is especially helpful for debugging. The validation result will be
		// logged as warning.
		Validator: &validator.Validator{},
	}

	// Then we add the directive processor to handle alexas powerController directive.
	// When creating an alexa handler calling the NewDefaultHandler method, all supported
	// directive processors are added by default.
	// You can also create new processor by yourself to be able to handle even more alexa
	// directives.
	handler.AddDirectiveProcessor(directives.CreatePowerControllerDirectiveProcessor())

	// Finally we create a http server to handle all incoming requests with the previously
	// created alexa handler.
	server := http.Server{
		Addr:    ":8181",
		Handler: &handler,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()
	defer server.Close()

	// Now we, or of course alexa, can post directives to that server.
	_, err := http.Post("http://localhost:8181", "application/json", strings.NewReader(
		`{
			"directive": {
			  "header": {
				"namespace": "Alexa.PowerController",
				"name": "TurnOn",
				"payloadVersion": "3",
				"messageId": "1bd5d003-31b9-476f-ad03-71d471922820",
				"correlationToken": "dFMb0z+PgpgdDmluhJ1LddFvSqZ/jCc8ptlAKulUj90jSqg=="
			  },
			  "endpoint": {
				"scope": {
				  "type": "BearerToken",
				  "token": "access-token-from-skill"
				},
				"endpointId": "appliance-001",
				"cookie": {
				  "type": "homematic",
				  "id": "4711"
				}
			  },
			  "payload": {}
			}
		  }`))

	if err != nil {
		panic(err)
	}

	// Output: homematic device 4711 state changed to true
}
