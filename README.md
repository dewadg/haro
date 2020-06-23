# Haro

Basic callback-based pubsub library with topic capability.

[![test](https://github.com/dewadg/haro/workflows/Tests/badge.svg?branch=master "test")](https://github.com/dewadg/haro/actions)

## Usage

Install with:
```
go get github.com/dewadg/haro
```

```go
pubsub := haro.New()

pubsub.DeclareTopic("topicName", "") // 2nd param will be used to check the payload type

if err := pubsub.RegisterSubscriber("topicName", func(ctx context.Context, payload string) {
    fmt.Println("Received " + payload)
}); err != nil {
    panic(err)
}

if err := pubsub.Publish(context.Background(), "topicName", "message"); err != nil {
    panic(err)
}
```
