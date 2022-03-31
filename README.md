# Haro

Topic-based pubsub library.

[![go report card](https://goreportcard.com/badge/github.com/dewadg/haro "go report card")](https://goreportcard.com/report/github.com/dewadg/haro)
[![test](https://github.com/dewadg/haro/workflows/Tests/badge.svg?branch=master "test")](https://github.com/dewadg/haro/actions)
[![Coverage Status](https://coveralls.io/repos/github/dewadg/haro/badge.svg?branch=master)](https://coveralls.io/github/dewadg/haro?branch=master)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/dewadg/haro?tab=doc)

## Usage

Install with:
```
go get github.com/dewadg/haro
```


### Declaring topic
```go
intTopic := haro.DeclareTopic[int]()

stringTopic := haro.DeclareTopic[string]()
```

### Registering subscribers
```go
intTopic.Subscribe(func(ctx context.Context, p int) error {
    return nil
})

stringTopic.Subscribe(func(ctx context.Context, p string) error {
return nil
})
```

### Publishing
```go
var err error

err = intTopic.Publish(context.Background(), 5)

err = stringTopic.Publish(context.Background(), "test")
```


## Limitations

- It's not a distributed pubsub mechanism, it means pubsub across different instances won't work (eg: across docker container replication)
- No recovery since it doesn't utilize any external storage, it's purely using Channel to store the published messages 
