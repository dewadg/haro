# Haro

Basic callback-based pubsub library with topic capability.

[![test](https://github.com/dewadg/haro/workflows/Tests/badge.svg?branch=master "test")](https://github.com/dewadg/haro/actions)

## Usage

Install with:
```
go get github.com/dewadg/haro
```

Create a new instance:
```go
p := haro.New()
```

### Declaring topic
`Pubsub.DeclareTopic` will receive 2 arguments; topic name and payload type. The payload type will be used to validate the pubsub processes. For example if you give string to it, the topic will be declared as a string-type topic, publishing integer to the topic will make error occured.
```go
var err error

// Register a string-type topic
err = p.DeclareTopic("stringTopic", "")

// Register a int-type topic
err = p.DeclareTopic("intTopic", 0)

// Register a uint-type topic
err = p.DeclareTopic("uintTopic", uint(0))
```

### Registering subscribers
`Pubub.RegisterSubscriber` can be used to register a subscriber (callback).

- A topic should be declared before registering a subscriber
- A topic can have more than 1 subscriber
- Registering subscriber with different payload type with the topic will return error

```go
// Prints the payload to uppercase
p.RegisterSubscriber(
    "stringTopic",
    func(ctx context.Context, payload string) error {
        fmt.Println(strings.ToUppercase(payload))

        return nil
    },,
)

// Prints the payload to lowercase
p.RegisterSubscriber(
    "stringTopic",
    func(ctx context.Context, payload string) error {
        fmt.Println(strings.ToLowercase(payload))
        
        return nil
    },,
)

// Prints the payload length
p.RegisterSubscriber(
    "stringTopic",
    func(ctx context.Context, payload string) error {
        fmt.Println(len(payload))
        
        return nil
    },,
)
```

Note that callback should have signature as a function with 2 parameters; `context.Context` and payload (payload type should match topic type).

We can also specify custom behavior for the registered subscriber. Currently, there are following possible behavior configurations:
- `haro.Retry(count int)` for configuring number of maximum retry when error occurred
- `haro.DelayRetry(duration time.Duration)` for configuring delay between each retry
- `haro.OnError(func(err error))` for specifying callback when error is occured

```go
// Set retry
p.RegisterSubscriber(
    "stringTopic",
    func(ctx context.Context, payload string) error {
        // Do something, and if error this callback will be called 2 times more
    },
    haro.Retry(2), // Will retry 2 times before returning the error
    haro.DelayRetry(2 * time.Second), // Set 2 seconds delay between retry
    haro.OnError(func(err error) {
        // Do something when error
    }),
)
```

### Publishing
Publishing is pretty straightforward:

```go
var err error

err = p.Publish(context.TODO(), "stringTopic", "message")

err = p.Publish(context.TODO(), "intTopic", 1)

err = p.Publish(context.TODO(), "uintTopic", uint(1))
```

`Pubsub.Publish` will return error if payload type between publishing and topic are mismatch

## Limitations

- It's not a distributed pubsub mechanism, it means pubsub across different instances won't work (eg: across docker container replication)
- No recovery since it doesn't utilize any external storage, it's purely using Channel to store the published events 
