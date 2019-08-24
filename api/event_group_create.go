package api

// GroupCreateEventType holds the event type string for this event.
const GroupCreateEventType = "group_create"

// GroupCreateEvent .
type GroupCreateEvent struct {
	Group *Group `json:"group"`
}

// Type returns the event's type.
func (event *GroupCreateEvent) Type() string {
	return GroupCreateEventType
}

// groupCreateEventHandler represents a GroupCreate event handler.
type groupCreateEventHandler func(*Library, *GroupCreateEvent)

// New .
func (handler groupCreateEventHandler) New() interface{} {
	return &GroupCreateEvent{}
}

// Handle calls the underlying handler.
func (handler groupCreateEventHandler) Handle(library *Library, i interface{}) {
	if event, ok := i.(*GroupCreateEvent); ok {
		handler(library, event)
	}
}

// Type returns the event's type.
func (handler groupCreateEventHandler) Type() string {
	return GroupCreateEventType
}
