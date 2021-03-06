package haro

import "context"

// Subscriber represents subscriber function
type Subscriber interface{}

// Pubsub exposes functions for pubsub operations
type Pubsub interface {
	// DeclareTopic registers a topic with desired payload type
	DeclareTopic(topicName string, payload interface{}) error

	// Publish dispatches an event to a topic
	Publish(ctx context.Context, topicName string, payload interface{}) error

	// RegisterSubscriber registers a subscriber for a topic. You can register more than 1 subscriber
	RegisterSubscriber(topicName string, callback Subscriber, configs ...ConfigFunc) error
}
