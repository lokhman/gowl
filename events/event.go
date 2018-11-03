package events

import (
	"github.com/lokhman/gowl/types"
)

// EventType
type EventType string

// EventInterface
type EventInterface interface {
	StopPropagation()
	isPropagationStopped() bool
}

// Event
type Event struct {
	Data types.Data

	propagationStopped bool
}

func (e *Event) StopPropagation() {
	e.propagationStopped = true
}

func (e *Event) isPropagationStopped() bool {
	return e.propagationStopped
}
