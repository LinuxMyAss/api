package api

// GroupUpdateEventType holds the event type string for this event.
const GroupUpdateEventType = "group_update"

// GroupUpdateEvent .
type GroupUpdateEvent struct {
	Group *Group `json:"group"`
}

// Type returns the event's type.
func (event *GroupUpdateEvent) Type() string {
	return GroupUpdateEventType
}

// groupUpdateEventHandler represents a GroupUpdate event handler.
type groupUpdateEventHandler func(*Library, *GroupUpdateEvent)

// New .
func (handler groupUpdateEventHandler) New() interface{} {
	return &GroupUpdateEvent{}
}

// Handle calls the underlying handler.
func (handler groupUpdateEventHandler) Handle(library *Library, i interface{}) {
	if event, ok := i.(*GroupUpdateEvent); ok {
		handler(library, event)
	}
}

// Type returns the event's type.
func (handler groupUpdateEventHandler) Type() string {
	return GroupUpdateEventType
}
