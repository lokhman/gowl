package events

type Emitter map[EventType][]func(EventInterface)

func (e Emitter) On(eventType EventType, listener func(event EventInterface)) {
	e[eventType] = append(e.Listeners(eventType), listener)
}

func (e Emitter) Emit(eventType EventType, event EventInterface) {
	for _, listener := range e.Listeners(eventType) {
		listener(event)

		if event.isPropagationStopped() {
			break
		}
	}
}

func (e Emitter) Listeners(eventType EventType) []func(EventInterface) {
	return e[eventType]
}

func (e Emitter) HasListeners(eventType EventType) bool {
	return len(e.Listeners(eventType)) > 0
}

func (e Emitter) RemoveAllListeners(eventType EventType) {
	delete(e, eventType)
}

func (e Emitter) Copy() Emitter {
	emitter := make(Emitter)
	for eventType, listeners := range e {
		var _listeners []func(EventInterface)
		copy(_listeners, listeners)
		emitter[eventType] = _listeners
	}
	return emitter
}
