package gowl

type EventType string

const (
	EventRequest  EventType = "request"
	EventResponse EventType = "response"
	EventPanic    EventType = "panic"
)

// EventEmitterInterface
type EventEmitterInterface interface {
	On(eventType EventType, listener func(event EventInterface))
	Emit(eventType EventType, event EventInterface)
	Listeners(eventType EventType) []func(EventInterface)
	HasListeners(eventType EventType) bool
	RemoveAllListeners(eventType EventType)
	Copy() EventEmitter
}

// EventEmitter
type EventEmitter map[EventType][]func(EventInterface)

func (e EventEmitter) On(eventType EventType, listener func(event EventInterface)) {
	e[eventType] = append(e.Listeners(eventType), listener)
}

func (e EventEmitter) Emit(eventType EventType, event EventInterface) {
	for _, listener := range e.Listeners(eventType) {
		listener(event)

		if event.isPropagationStopped() {
			break
		}
	}
}

func (e EventEmitter) Listeners(eventType EventType) []func(EventInterface) {
	return e[eventType]
}

func (e EventEmitter) HasListeners(eventType EventType) bool {
	return len(e.Listeners(eventType)) > 0
}

func (e EventEmitter) RemoveAllListeners(eventType EventType) {
	delete(e, eventType)
}

func (e EventEmitter) Copy() EventEmitter {
	emitter := make(EventEmitter)
	for eventType, listeners := range e {
		var _listeners []func(EventInterface)
		copy(_listeners, listeners)
		emitter[eventType] = _listeners
	}
	return emitter
}

// EventInterface
type EventInterface interface {
	StopPropagation()
	isPropagationStopped() bool
}

// Event
type Event struct {
	Data Data

	propagationStopped bool
}

func (e *Event) StopPropagation() {
	e.propagationStopped = true
}

func (e *Event) isPropagationStopped() bool {
	return e.propagationStopped
}

// RequestEvent
type RequestEvent struct {
	Event

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
	Event

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
	Event

	error error
}

func (e *PanicEvent) Error() error {
	return e.error
}
