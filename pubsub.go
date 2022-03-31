package haro

import (
	"context"
	"sync"
	"time"
)

type Subscriber[Payload any] func(context.Context, Payload) error

type Topic[Payload any] interface {
	Publish(context.Context, Payload) error

	Subscribe(sub Subscriber[Payload])
}

type payloadPair[Payload any] struct {
	ctx     context.Context
	payload Payload
}

type topic[Payload any] struct {
	cfg    *config
	mtx    sync.Mutex
	subs   []Subscriber[Payload]
	stream chan payloadPair[Payload]
}

func DeclareTopic[Payload any](cfgFuncs ...ConfigFunc) Topic[Payload] {
	cfg := config{}
	for _, f := range cfgFuncs {
		f(&cfg)
	}

	t := &topic[Payload]{
		cfg:    &cfg,
		mtx:    sync.Mutex{},
		subs:   make([]Subscriber[Payload], 0),
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
				for _, sub := range t.subs {
					callSubscriber[any](t.cfg, sub, p)
				}
			}
		}
	}()
}

func callSubscriber[Payload any](
	cfg *config,
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

func (t *topic[any]) Subscribe(sub Subscriber[any]) {
	t.mtx.Lock()
	if t.subs == nil {
		t.subs = make([]Subscriber[any], 0)
	}
	t.subs = append(t.subs, sub)
	t.mtx.Unlock()
}
