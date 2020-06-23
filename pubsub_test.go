package haro

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
)

func Test_pubsub_Publish(t *testing.T) {
	type fields struct {
		registry registry
	}
	type args struct {
		ctx        context.Context
		topicName  string
		payload    interface{}
		resultWg   sync.WaitGroup
		resultChan chan string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		configureMock func(p Pubsub, fields fields, args args)
	}{
		{
			name: "Event published",
			fields: fields{
				registry: map[string]*topic{},
			},
			args: args{
				ctx:        context.Background(),
				topicName:  "topic",
				payload:    "payload",
				resultWg:   sync.WaitGroup{},
				resultChan: make(chan string, 2), // Accept only 2 results
			},
			wantErr: false,
			configureMock: func(p Pubsub, fields fields, args args) {
				p.DeclareTopic(args.topicName, args.payload)

				p.RegisterSubscriber(args.topicName, func(ctx context.Context, payload string) error {
					return nil
				})

				p.RegisterSubscriber(args.topicName, func(ctx context.Context, payload string) error {
					return errors.New("error")
				})
			},
		},
		{
			name: "Undeclared topic",
			fields: fields{
				registry: map[string]*topic{},
			},
			args: args{
				ctx:        context.Background(),
				topicName:  "topic",
				payload:    "payload",
				resultWg:   sync.WaitGroup{},
				resultChan: make(chan string, 2), // Accept only 2 results
			},
			wantErr: true,
			configureMock: func(p Pubsub, fields fields, args args) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pubsub{
				registry: tt.fields.registry,
			}

			tt.configureMock(p, tt.fields, tt.args)

			if err := p.Publish(tt.args.ctx, tt.args.topicName, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pubsub_RegisterSubscriber(t *testing.T) {
	type fields struct {
		registry registry
		mutex    sync.Mutex
	}
	type args struct {
		topicName string
		handler   Subscriber
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		configureMock func(p Pubsub, fields fields, args args)
	}{
		{
			name: "Successfully registered a subscriber",
			fields: fields{
				registry: make(map[string]*topic),
				mutex:    sync.Mutex{},
			},
			args: args{
				topicName: "test",
				handler: func(ctx context.Context, payload string) error {
					return nil
				},
			},
			wantErr: false,
			configureMock: func(p Pubsub, fields fields, args args) {
				p.DeclareTopic(args.topicName, "")
			},
		},
		{
			name: "Undeclared topic",
			fields: fields{
				registry: make(map[string]*topic),
				mutex:    sync.Mutex{},
			},
			args: args{
				topicName: "test",
				handler: func(ctx context.Context, payload string) error {
					return nil
				},
			},
			wantErr: true,
			configureMock: func(p Pubsub, fields fields, args args) {
			},
		},
		{
			name: "Handler not a function",
			fields: fields{
				registry: make(map[string]*topic),
				mutex:    sync.Mutex{},
			},
			args: args{
				topicName: "test",
				handler:   "function",
			},
			wantErr: true,
			configureMock: func(p Pubsub, fields fields, args args) {
			},
		},
		{
			name: "Handler parameter count invalid",
			fields: fields{
				registry: make(map[string]*topic),
				mutex:    sync.Mutex{},
			},
			args: args{
				topicName: "test",
				handler:   func() {},
			},
			wantErr: true,
			configureMock: func(p Pubsub, fields fields, args args) {
			},
		},
		{
			name: "Handler first param is not a Context",
			fields: fields{
				registry: make(map[string]*topic),
				mutex:    sync.Mutex{},
			},
			args: args{
				topicName: "test",
				handler:   func(ctx int, payload int) {},
			},
			wantErr: true,
			configureMock: func(p Pubsub, fields fields, args args) {
			},
		},
		{
			name: "Handler return type invalid",
			fields: fields{
				registry: make(map[string]*topic),
				mutex:    sync.Mutex{},
			},
			args: args{
				topicName: "test",
				handler: func(ctx context.Context, payload int) {
				},
			},
			wantErr: true,
			configureMock: func(p Pubsub, fields fields, args args) {
			},
		},
		{
			name: "Handler return type is not an error",
			fields: fields{
				registry: make(map[string]*topic),
				mutex:    sync.Mutex{},
			},
			args: args{
				topicName: "test",
				handler: func(ctx context.Context, payload int) int {
					return 0
				},
			},
			wantErr: true,
			configureMock: func(p Pubsub, fields fields, args args) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pubsub{
				registry: tt.fields.registry,
				mutex:    tt.fields.mutex,
			}

			tt.configureMock(p, tt.fields, tt.args)

			if err := p.RegisterSubscriber(tt.args.topicName, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("RegisterSubscriber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pubsub_DeclareTopic(t *testing.T) {
	type fields struct {
		registry registry
		mutex    sync.Mutex
	}
	type args struct {
		topicName string
		payload   interface{}
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		configureMock func(fields fields, args args)
	}{
		{
			name: "Successfully registered a topic",
			fields: fields{
				registry: make(map[string]*topic),
				mutex:    sync.Mutex{},
			},
			args: args{
				topicName: "test",
				payload:   "string",
			},
			wantErr:       false,
			configureMock: func(fields fields, args args) {},
		},
		{
			name: "Duplicate topic with different payload type",
			fields: fields{
				registry: make(map[string]*topic),
				mutex:    sync.Mutex{},
			},
			args: args{
				topicName: "test",
				payload:   "string",
			},
			wantErr: true,
			configureMock: func(fields fields, args args) {
				fields.registry.New("test", reflect.TypeOf(0).String())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pubsub{
				registry: tt.fields.registry,
				mutex:    tt.fields.mutex,
			}

			tt.configureMock(tt.fields, tt.args)

			if err := p.DeclareTopic(tt.args.topicName, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("DeclareTopic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
