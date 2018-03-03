package capabilities

// HealthConscious specifies the ability to report health (like connectivity capability)
type HealthConscious interface {
	IsConnected() bool
}
