package api

// Event .
type Event interface {
	Type() string
}

// EventHandler .
type EventHandler interface {
	Handle(*Library, interface{})
	Type() string
}

// EventInterfaceProvider .
type EventInterfaceProvider interface {
	New() interface{}
	Type() string
}

type eventHandlerInstance struct {
	eventHandler EventHandler
}

// interface
type interfaceEventHandler func(*Library, interface{})

func (handler interfaceEventHandler) Handle(library *Library, i interface{}) {
	handler(library, i)
}

func (handler interfaceEventHandler) Type() string {
	return "__INTERFACE__"
}

// END interface
