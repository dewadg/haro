package haro

import (
	"context"
	"log"
	"reflect"
	"sync"

	"github.com/google/uuid"
)

type pubsub struct {
	registry registry
	mutex    sync.Mutex
}

func (p *pubsub) DeclareTopic(topicName string, payload interface{}) {
	p.mutex.Lock()
	if !p.registry.Exists(topicName) {
		p.createTopic(topicName, payload)
	}
	p.mutex.Unlock()
}

func (p *pubsub) Publish(ctx context.Context, topicName string, payload interface{}) error {
	p.mutex.Lock()
	_topic, err := p.registry.Get(topicName)
	if err != nil {
		p.mutex.Unlock()
		return err
	}
	p.mutex.Unlock()

	go func() {
		_topic.Events <- event{
			Ctx:     ctx,
			Payload: payload,
		}
	}()

	return nil
}

func (p *pubsub) validateSubscriber(handler Subscriber) error {
	handlerType := reflect.TypeOf(handler)

	// Check if subscriber is indeed a function
	if handlerType.Kind() != reflect.TypeOf(func() {}).Kind() {
		return ErrSubscriberNotAFunction
	}

	// Check if subscriber has exactly 2 parameters
	if handlerType.NumIn() != 2 {
		return ErrSubscriberInvalidParameters
	}

	// Check subscriber first param
	if handlerType.In(0).String() != "context.Context" {
		return ErrSubscriberInvalidFirstParameter
	}

	// Check if subscriber has exactly 1 return type
	if handlerType.NumOut() != 1 {
		return ErrSubscriberInvalidReturnTypes
	}

	// Check if subscriber has exactly 1 return type
	if handlerType.Out(0).String() != "error" {
		return ErrSubscriberInvalidFirstReturnType
	}

	return nil
}

func (p *pubsub) RegisterSubscriber(topicName string, handler Subscriber) error {
	if err := p.validateSubscriber(handler); err != nil {
		return err
	}

	var _topic *topic

	p.mutex.Lock()
	_topic, err := p.registry.Get(topicName)
	if err != nil {
		p.mutex.Unlock()
		return err
	}
	p.mutex.Unlock()

	_topic.Subscribers = append(_topic.Subscribers, subscriber{
		ID:      uuid.New().String(),
		Handler: handler,
	})

	return nil
}

func (p *pubsub) createTopic(topicName string, payload interface{}) *topic {
	reflectedPayload := reflect.TypeOf(payload)
	_topic := p.registry.New(topicName, reflectedPayload.String())

	go func() {
		for {
			select {
			case evt := <-_topic.Events:
				for _, subscriber := range _topic.Subscribers {
					if err := p.callSubscriber(subscriber.Handler, _topic, evt); err != nil {
						log.Println(err)
					}
				}
				break
			}
		}
	}()

	return _topic
}

func (p *pubsub) callSubscriber(handler Subscriber, _topic *topic, evt event) error {
	// Check payload type, reject if mismatch
	payloadType := reflect.TypeOf(evt.Payload)
	if payloadType.String() != _topic.PayloadType {
		return ErrPayloadTypeMismatch
	}

	args := []reflect.Value{
		reflect.ValueOf(evt.Ctx),
		reflect.ValueOf(evt.Payload),
	}

	results := reflect.ValueOf(handler).Call(args)
	err := results[0].Interface()

	if err != nil {
		return err.(error)
	}
	return nil
}
