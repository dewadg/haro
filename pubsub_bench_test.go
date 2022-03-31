package haro

import (
	"context"
	"testing"
)

func BenchmarkPublish_int(b *testing.B) {
	t := DeclareTopic[int]()

	for i := 0; i < b.N; i++ {
		_ = t.Publish(context.Background(), 5)
	}
}

func BenchmarkPublish_string(b *testing.B) {
	t := DeclareTopic[string]()

	for i := 0; i < b.N; i++ {
		_ = t.Publish(context.Background(), "test")
	}
}

func BenchmarkPublish_struct(b *testing.B) {
	type d struct {
		id   int64
		name string
	}

	t := DeclareTopic[d]()

	for i := 0; i < b.N; i++ {
		_ = t.Publish(context.Background(), d{
			id:   1,
			name: "test",
		})
	}
}
