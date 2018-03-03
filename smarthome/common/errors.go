package common

// AlexaError is an error which additional holds the alexa error type
// https://developer.amazon.com/de/docs/device-apis/alexa-errorresponse.html#example-errorresponse
type AlexaError struct {
	Type      string
	Message   string
	Namespace string
}

func (e AlexaError) Error() string {
	return e.Message
}

// NewInvalidDirectiveError creates an AlexaError to indicate a directive is not valid for this skill or is malformed.
func NewInvalidDirectiveError(message string) AlexaError {
	return AlexaError{"INVALID_DIRECTIVE", message, "Alexa"}
}

// NewInternalError creates an AlexaError to indicate an error that cannot be accurately described as one of the
// other error types occurred while you were handling the directive.
// For example, a generic runtime exception occurred while handling a directive.
// Ideally, you will never send this error event, but instead send a more specific error type.
func NewInternalError(message string) AlexaError {
	return AlexaError{"INTERNAL_ERROR", message, "Alexa"}
}

// NewAcceptGrantFailedError creates an AlexaError to indicate that user authentication failed
func NewAcceptGrantFailedError(message string) AlexaError {
	return AlexaError{"ACCEPT_GRANT_FAILED", message, "Alexa.Authorization"}
}
