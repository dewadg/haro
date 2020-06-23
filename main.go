package haro

// New :nodoc:
func New() Pubsub {
	return &pubsub{
		registry: make(map[string]*topic),
	}
}
