package common

// PowerDevice specifies an device with powerState capabilities
type PowerDevice interface {
	SetState(bool) (bool, error)
	State() (bool, error)
}

// TemperatureSensor specifies an device with temperature capabilities
type TemperatureSensor interface {
	Temperature() float32
}

// HealthConscious specifies the ability to report health (like connectivity capability)
type HealthConscious interface {
	IsConnected() bool
}
