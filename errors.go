package haro

import "errors"

// ErrSubscriberNotAFunction indicates if the subscriber is not a function
var ErrSubscriberNotAFunction = errors.New("Subscriber should be a function")

// ErrSubscriberInvalidParameters indicates if the subscriber doesn't have 2 parameters
var ErrSubscriberInvalidParameters = errors.New("Subscriber should have exactly 2 parameters")

// ErrSubscriberInvalidFirstParameter indicates if the subscriber's first parameter should be context.Context
var ErrSubscriberInvalidFirstParameter = errors.New("Subscriber should have `context.Context` as the first parameter")

// ErrSubscriberInvalidReturnTypes indicates if the subscriber should only have 1 return type
var ErrSubscriberInvalidReturnTypes = errors.New("Subscriber should have 1 return type")

// ErrSubscriberInvalidFirstReturnType indicates if the subscriber's return type should be error
var ErrSubscriberInvalidFirstReturnType = errors.New("Subscriber should have return type `error`")

// ErrPayloadTypeMismatch indicates if the payload type between subscriber and event not match
var ErrPayloadTypeMismatch = errors.New("Mismatch parameter between subscriber and event payload")
