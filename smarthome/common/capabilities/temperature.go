package capabilities

// TemperatureSensor specifies an device with temperature capabilities
type TemperatureSensor interface {
	Temperature() float32
}
