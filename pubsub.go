package haro

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
)

type pubsub struct {
	registry registry
	mutex    sync.Mutex
}

func (p *pubsub) DeclareTopic(topicName string, payload interface{}) error {
	p.mutex.Lock()
	if !p.registry.Exists(topicName) {
		p.createTopic(topicName, payload)
	} else {
		_topic, err := p.registry.Get(topicName)
		if err != nil {
			return err
		}

		reflectedPayload := reflect.TypeOf(payload)
		if reflectedPayload.String() != _topic.PayloadType {
			return fmt.Errorf("Topic `%s` is already declared with payload type %s", topicName, _topic.PayloadType)
		}
	}
	p.mutex.Unlock()

	return nil
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

func (p *pubsub) RegisterSubscriber(topicName string, handler Subscriber, configs ...ConfigFunc) error {
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

	cfg := &config{}
	for _, configure := range configs {
		configure(cfg)
	}

	_topic.Subscribers = append(_topic.Subscribers, subscriber{
		ID:      uuid.New().String(),
		Handler: handler,
		Config:  cfg,
	})

	return nil
}

func (p *pubsub) createTopic(topicName string, payload interface{}) *topic {
	reflectedPayload := reflect.TypeOf(payload)
	_topic := p.registry.New(topicName, reflectedPayload.String())

	go p.runSubscribers(_topic)

	return _topic
}

func (p *pubsub) runSubscribers(_topic *topic) {
	for {
		select {
		case evt := <-_topic.Events:
			for _, subscriber := range _topic.Subscribers {
				var err error
				// If subscriber is set to retry
				if subscriber.Config.Retry > 0 {
					err = p.callSubscriberWithRetry(subscriber, _topic, evt)
				} else {
					err = p.callSubscriber(subscriber.Handler, _topic, evt)
				}

				if subscriber.Config.OnError != nil {
					subscriber.Config.OnError(err)
				}
			}
			break
		}
	}
}

func (p *pubsub) callSubscriberWithRetry(_subscriber subscriber, _topic *topic, evt event) error {
	var retryCount int
	var lastError error

	for retryCount < _subscriber.Config.Retry {
		if err := p.callSubscriber(_subscriber.Handler, _topic, evt); err != nil {
			lastError = err
			retryCount++

			if _subscriber.Config.RetryDelay > 0 {
				time.Sleep(_subscriber.Config.RetryDelay)
			}
		} else {
			break
		}
	}

	return lastError
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
