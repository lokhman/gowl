package gowl

import (
	"github.com/lokhman/gowl/events"
)

type EventInterface = events.EventInterface

const (
	EventRequest  events.EventType = "request"
	EventResponse events.EventType = "response"
	EventPanic    events.EventType = "panic"
)

// RequestEvent
type RequestEvent struct {
	events.Event

	request  *Request
	response ResponseInterface
}

func (e *RequestEvent) Request() *Request {
	return e.request
}

func (e *RequestEvent) SetResponse(response ResponseInterface) {
	e.response = response
}

// ResponseEvent
type ResponseEvent struct {
	events.Event

	request  *Request
	response ResponseInterface
}

func (e *ResponseEvent) Request() *Request {
	return e.request
}

func (e *ResponseEvent) Response() ResponseInterface {
	return e.response
}

func (e *ResponseEvent) SetResponse(response ResponseInterface) {
	e.response = response
}

// PanicEvent
type PanicEvent struct {
	events.Event

	error error
}

func (e *PanicEvent) Error() error {
	return e.error
}
