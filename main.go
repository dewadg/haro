package haro

// New returns new Pubsub instance
func New() Pubsub {
	return &pubsub{
		registry: make(map[string]*topic),
	}
}
