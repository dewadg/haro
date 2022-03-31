package haro

import (
	"context"
	"testing"
)

func BenchmarkPublish_int(b *testing.B) {
	t := DeclareTopic[int]()
	outputChan := make(chan int, b.N)

	t.Subscribe(func(ctx context.Context, p int) error {
		outputChan <- p
		return nil
	})

	for i := 0; i < b.N; i++ {
		_ = t.Publish(context.Background(), 5)

		<-outputChan
	}
}

func BenchmarkPublish_string(b *testing.B) {
	t := DeclareTopic[string]()
	outputChan := make(chan string, b.N)

	t.Subscribe(func(ctx context.Context, p string) error {
		outputChan <- p
		return nil
	})

	for i := 0; i < b.N; i++ {
		_ = t.Publish(context.Background(), "test")

		<-outputChan
	}
}

func BenchmarkPublish_struct(b *testing.B) {
	type d struct {
		id   int64
		name string
	}

	t := DeclareTopic[d]()
	outputChan := make(chan d, b.N)

	t.Subscribe(func(ctx context.Context, p d) error {
		outputChan <- p
		return nil
	})

	for i := 0; i < b.N; i++ {
		_ = t.Publish(context.Background(), d{
			id:   1,
			name: "test",
		})

		<-outputChan
	}
}
