package haro

import (
	"context"
	"testing"
)

func BenchmarkPublish_int(b *testing.B) {
	ps := New()

	_ = ps.DeclareTopic("int", 0)

	for i := 0; i < b.N; i++ {
		_ = ps.Publish(context.Background(), "int", 1)
	}
}

func BenchmarkPublish_string(b *testing.B) {
	ps := New()
	outputChan := make(chan string, b.N)

	_ = ps.DeclareTopic("string", "")

	_ = ps.RegisterSubscriber("string", func(ctx context.Context, p string) error {
		outputChan <- p

		return nil
	})

	for i := 0; i < b.N; i++ {
		_ = ps.Publish(context.Background(), "string", "test")

		<-outputChan
	}
}

func BenchmarkPublish_struct(b *testing.B) {
	type d struct {
		id   int64
		name string
	}

	ps := New()

	_ = ps.DeclareTopic("struct", d{})

	for i := 0; i < b.N; i++ {
		_ = ps.Publish(context.Background(), "struct", d{
			id:   1,
			name: "test",
		})
	}
}
