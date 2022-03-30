package haro

import (
	"context"
	"sync"
	"testing"
)

func Test_topic_Publish_int(t1 *testing.T) {
	type fields struct {
		mtx    sync.Mutex
		subs   []Subscriber[int]
		stream chan payloadPair[int]
	}
	type args struct {
		ctx context.Context
		a   int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				mtx:    sync.Mutex{},
				subs:   make([]Subscriber[int], 0),
				stream: make(chan payloadPair[int]),
			},
			args: args{
				ctx: context.Background(),
				a:   2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &topic[int]{
				mtx:    tt.fields.mtx,
				subs:   tt.fields.subs,
				stream: tt.fields.stream,
			}
			if err := t.Publish(tt.args.ctx, tt.args.a); (err != nil) != tt.wantErr {
				t1.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := <-t.stream
			if got.payload != tt.args.a {
				t1.Errorf("Publish() error, got = %v, want %v", got, tt.args.a)
			}
		})
	}
}

func Test_topic_Subscribe_int(t1 *testing.T) {
	type fields struct {
		mtx    sync.Mutex
		subs   []Subscriber[int]
		stream chan payloadPair[int]
		wg     sync.WaitGroup
	}
	type args struct {
		sub func(wg *sync.WaitGroup) Subscriber[int]
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				mtx:    sync.Mutex{},
				subs:   make([]Subscriber[int], 0),
				stream: make(chan payloadPair[int]),
				wg:     sync.WaitGroup{},
			},
			args: args{
				sub: func(wg *sync.WaitGroup) Subscriber[int] {
					return func(ctx context.Context, p int) error {
						wg.Done()

						return nil
					}
				},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &topic[int]{
				mtx:    tt.fields.mtx,
				subs:   tt.fields.subs,
				stream: tt.fields.stream,
			}

			t.run()
			t.Subscribe(tt.args.sub(&tt.fields.wg))

			tt.fields.wg.Add(1)
			_ = t.Publish(context.Background(), 1)
			tt.fields.wg.Wait()
		})
	}
}

func Test_topic_Subscribe_struct(t1 *testing.T) {
	type example struct {
		val int
	}
	type fields struct {
		mtx    sync.Mutex
		subs   []Subscriber[example]
		stream chan payloadPair[example]
		wg     sync.WaitGroup
	}
	type args struct {
		sub func(wg *sync.WaitGroup) Subscriber[example]
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				mtx:    sync.Mutex{},
				subs:   make([]Subscriber[example], 0),
				stream: make(chan payloadPair[example]),
				wg:     sync.WaitGroup{},
			},
			args: args{
				sub: func(wg *sync.WaitGroup) Subscriber[example] {
					return func(ctx context.Context, p example) error {
						wg.Done()

						return nil
					}
				},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &topic[example]{
				mtx:    tt.fields.mtx,
				subs:   tt.fields.subs,
				stream: tt.fields.stream,
			}

			t.run()
			t.Subscribe(tt.args.sub(&tt.fields.wg))

			tt.fields.wg.Add(1)
			_ = t.Publish(context.Background(), example{
				val: 1,
			})
			tt.fields.wg.Wait()
		})
	}
}
