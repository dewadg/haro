package haro

import (
	"context"
	"sync"
	"time"
)

type Subscriber[Payload any] func(context.Context, Payload) error

type Topic[Payload any] interface {
	Publish(context.Context, Payload) error

	Subscribe(sub Subscriber[Payload], configurators ...SubscriberConfigFunc)
}

type payloadPair[Payload any] struct {
	ctx     context.Context
	payload Payload
}

type topic[Payload any] struct {
	mtx    sync.Mutex
	subs   []subscriberPair[Payload]
	stream chan payloadPair[Payload]
}

type subscriberPair[Payload any] struct {
	subscriber Subscriber[Payload]
	cfg        *subscriberConfig
}

func DeclareTopic[Payload any]() Topic[Payload] {
	t := &topic[Payload]{
		mtx:    sync.Mutex{},
		subs:   make([]subscriberPair[Payload], 0),
		stream: make(chan payloadPair[Payload], 0),
	}

	t.run()

	return t
}

func (t *topic[any]) run() {
	go func() {
		for {
			select {
			case p := <-t.stream:
				for _, pair := range t.subs {
					callSubscriber[any](pair.cfg, pair.subscriber, p)
				}
			}
		}
	}()
}

func callSubscriber[Payload any](
	cfg *subscriberConfig,
	sub Subscriber[Payload],
	p payloadPair[Payload],
) {
	if cfg.retry <= 0 {
		err := sub(p.ctx, p.payload)
		if err != nil && cfg.onError != nil {
			cfg.onError(err)
		}
		if err == nil && cfg.onSuccess != nil {
			cfg.onSuccess()
		}
	} else {
		for i := 0; i < cfg.retry; i++ {
			err := sub(p.ctx, p.payload)
			if err != nil && cfg.onError != nil {
				cfg.onError(err)
			}

			if err != nil {
				time.Sleep(cfg.retryDelay)
			} else {
				if cfg.onSuccess != nil {
					cfg.onSuccess()
				}

				break
			}
		}
	}
}

func (t *topic[any]) Publish(ctx context.Context, a any) error {
	go func() {
		t.stream <- payloadPair[any]{
			ctx:     ctx,
			payload: a,
		}
	}()

	return nil
}

func (t *topic[any]) Subscribe(sub Subscriber[any], configurators ...SubscriberConfigFunc) {
	t.mtx.Lock()
	if t.subs == nil {
		t.subs = make([]subscriberPair[any], 0)
	}

	var cfg subscriberConfig
	for _, c := range configurators {
		c(&cfg)
	}

	t.subs = append(t.subs, subscriberPair[any]{
		subscriber: sub,
		cfg:        &cfg,
	})

	t.mtx.Unlock()
}
