package smarthome

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/betom84/go-alexa/smarthome/common"
	"github.com/betom84/go-alexa/smarthome/directives"
	"github.com/betom84/go-alexa/smarthome/validator"
)

// DeviceFactory creates the devices to handle the alexa endpoint capability
type DeviceFactory interface {

	// NewDevice creates a device for the given type and id, the created device should support
	// the capabilities defined by the endpoint
	NewDevice(epType string, id string) (interface{}, error)
}

// Handler is a http server to handle alexa directives
type Handler struct {
	BasicAuth struct {
		Username string
		Password string
	}

	// DeviceFactory to creates the devices
	DeviceFactory DeviceFactory

	// Validator to ensure correct response formats, optional
	Validator *validator.Validator

	// Processors to handle directives
	directiveProcessors []directives.DirectiveProcessor
}

// NewDefaultHandler creates an instance to handle all supported alexa directives.
func NewDefaultHandler(authority *Authority, endpoints io.Reader) *Handler {
	handler := new(Handler)

	handler.AddDirectiveProcessor(directives.CreateAuthorizeDirectiveProcessor(authority))
	handler.AddDirectiveProcessor(directives.CreateDiscoveryDirectiveProcessor(endpoints))
	handler.AddDirectiveProcessor(directives.CreatePowerControllerDirectiveProcessor())
	handler.AddDirectiveProcessor(directives.CreateReportAlexaDirectiveProcessor())

	return handler
}

// AddDirectiveProcessor is used to add an directive processor to this instance
func (h *Handler) AddDirectiveProcessor(processor directives.DirectiveProcessor) {
	h.directiveProcessors = append(h.directiveProcessors, processor)
}

// ServeHTTP is needed to satisfy the net/http.Handler interface. Therefore alexa.Handler can be used as http handler.
func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if ok := h.verifyBasicAuth(request); !ok {
		h.writeUnauthorizedHTTPResponse(writer)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}

	dir, err := h.getDirectiveFromRequestBody(body)
	if err != nil {
		h.writeBadRequestHTTPResponse(writer, err)
		return
	}

	resp, err := json.Marshal(h.handleDirective(dir))
	if err != nil {
		h.writeBadRequestHTTPResponse(writer, err)
		return
	}

	h.validateOnDemand(resp)

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(resp)
}

func (h *Handler) validateOnDemand(payload []byte) {
	if h.Validator == nil {
		return
	}

	validationStart := time.Now()
	if err := h.Validator.Validate(payload); err != nil {
		log.Printf("WARNING Response validation failed\npayload:\n%s\nerrors:%s\n", string(payload), err)
	} else {
		log.Printf("Response validated without errors within %.3fs (Schema: %s)", time.Now().Sub(validationStart).Seconds(), h.Validator.SchemaReference)
	}
}

func (h *Handler) verifyBasicAuth(request *http.Request) bool {
	if h.BasicAuth.Username != "" {
		user, pass, ok := request.BasicAuth()
		return ok && user == h.BasicAuth.Username && pass == h.BasicAuth.Password
	}

	return true
}

func (h *Handler) writeUnauthorizedHTTPResponse(writer http.ResponseWriter) {
	log.Println("Unauthorized request rejected")

	writer.WriteHeader(http.StatusUnauthorized)
	writer.Header().Set("Content-Type", "text/plain")
	writer.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func (h *Handler) writeBadRequestHTTPResponse(writer http.ResponseWriter, err error) {
	log.Println("Unable to handle request,", err)

	writer.WriteHeader(http.StatusBadRequest)
	writer.Header().Set("Content-Type", "text/plain")
	writer.Write([]byte(http.StatusText(http.StatusBadRequest)))
}

func (h *Handler) getDirectiveFromRequestBody(body []byte) (dir *common.Directive, err error) {
	var raw map[string]*json.RawMessage
	err = json.Unmarshal(body, &raw)
	if err != nil {
		log.Printf("%s; %s", err, body)
		return nil, fmt.Errorf("Failed to unmarshal request body")
	}

	if rawDir, ok := raw["directive"]; ok {
		dir, err = common.NewDirective(*rawDir)
	} else {
		return nil, fmt.Errorf("Request does not contain valid alexa directive")
	}

	return
}

func (h *Handler) handleDirective(dir *common.Directive) (r interface{}) {
	var err error

	startTime := time.Now()
	log.Printf("Received directive %s", dir)

	for _, processor := range h.directiveProcessors {
		if !processor.IsCapable(dir) {
			continue
		}

		var device interface{}
		if dir.Endpoint != nil {
			device, err = h.DeviceFactory.NewDevice(dir.Endpoint.Cookie.Type, dir.Endpoint.Cookie.ID)
			if err != nil {
				log.Printf("Unable to create endpoint device (%v)", err)
			}
		}

		r, err = processor.Process(dir, device)
		if err != nil {
			r = h.createErrorResponse(dir, h.transformError(err))
			log.Print(err)
		}

		log.Printf("Processed %s in %.3fs", dir, time.Now().Sub(startTime).Seconds())

		return
	}

	r = h.createErrorResponse(dir, common.NewInvalidDirectiveError("Directive not supported"))
	return
}

func (h *Handler) createErrorResponse(dir *common.Directive, err common.AlexaError) (resp *common.Response) {
	resp = new(common.Response)
	resp.Event.Header = common.NewHeader("ErrorResponse", err.Namespace)
	resp.Event.Header.CorrelationToken = dir.Header.CorrelationToken
	resp.Event.Endpoint = dir.Endpoint

	resp.Event.Payload = struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}{
		Type:    err.Type,
		Message: err.Message,
	}

	return
}

func (h *Handler) transformError(err error) common.AlexaError {
	switch e := err.(type) {
	case common.AlexaError:
		return e
	default:
		return common.NewInternalError(e.Error())
	}
}
