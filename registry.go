package haro

import (
	"context"
	"fmt"
)

type event struct {
	Ctx     context.Context
	Payload interface{}
}

type subscriber struct {
	ID      string
	Handler Subscriber
	Config  *config
}

type topic struct {
	Events      chan event
	Subscribers []subscriber
	PayloadType string
}

type registry map[string]*topic

// Exists :nodoc:
func (t registry) Exists(topicNamestring) bool {
	_, exists := t[topicName]

	return exists
}

// GetChannel :nodoc:
func (t registry) Get(topicName string) (*topic, error) {
	notFoundError := fmt.Errorf("Topic `%s` does not exist", topicName)

	if value, ok := t[topicName]; ok {
		return value, nil
	}
	return nil, notFoundError
}

// New :nodoc:
func (t registry) New(topicName string, payloadType string) *topic {
	if t.Exists(topicName) {
		topic, _ := t.Get(topicName)

		return topic
	}

	value := &topic{
		Events:      make(chan event),
		Subscribers: make([]subscriber, 0),
		PayloadType: payloadType,
	}

	t[topicName] = value

	return value
}
