package haro

import "errors"

// ErrSubscriberNotAFunction :nodoc:
var ErrSubscriberNotAFunction = errors.New("Subscriber should be a function")

// ErrSubscriberInvalidParameters :nodoc:
var ErrSubscriberInvalidParameters = errors.New("Subscriber should have exactly 2 parameters")

// ErrSubscriberInvalidFirstParameter :nodoc:
var ErrSubscriberInvalidFirstParameter = errors.New("Subscriber should have `context.Context` as the first parameter")

// ErrSubscriberInvalidReturnTypes :nodoc:
var ErrSubscriberInvalidReturnTypes = errors.New("Subscriber should have 1 return type")

// ErrSubscriberInvalidFirstReturnType :nodoc:
var ErrSubscriberInvalidFirstReturnType = errors.New("Subscriber should have return type `error`")

// ErrPayloadTypeMismatch :nodoc:
var ErrPayloadTypeMismatch = errors.New("Mismatch parameter between subscriber and event payload")
