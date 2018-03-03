package capabilities

// PowerDevice specifies an device with powerState capabilities
type PowerDevice interface {
	SetState(bool) (bool, error)
	State() (bool, error)
}
