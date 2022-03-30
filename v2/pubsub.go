package haro

import (
	"context"
	"sync"
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
	mtx    sync.Mutex
	subs   []Subscriber[Payload]
	stream chan payloadPair[Payload]
}

func DeclareTopic[Payload any]() Topic[Payload] {
	return &topic[Payload]{
		mtx:    sync.Mutex{},
		subs:   make([]Subscriber[Payload], 0),
		stream: make(chan payloadPair[Payload], 0),
	}
}

func (t *topic[any]) run() {
	go func() {
		for {
			select {
			case p := <-t.stream:
				for _, sub := range t.subs {
					_ = sub(p.ctx, p.payload)
				}
			}
		}
	}()
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
