package api

// GroupDeleteEventType holds the event type string for this event.
const GroupDeleteEventType = "group_delete"

// GroupDeleteEvent .
type GroupDeleteEvent struct {
	ID string `json:"id"`
}

// Type returns the event's type.
func (event *GroupDeleteEvent) Type() string {
	return GroupDeleteEventType
}

// groupDeleteEventHandler represents a GroupDelete event handler.
type groupDeleteEventHandler func(*Library, *GroupDeleteEvent)

// New .
func (handler groupDeleteEventHandler) New() interface{} {
	return &GroupDeleteEvent{}
}

// Handle calls the underlying handler.
func (handler groupDeleteEventHandler) Handle(library *Library, i interface{}) {
	if event, ok := i.(*GroupDeleteEvent); ok {
		handler(library, event)
	}
}

// Type returns the event's type.
func (handler groupDeleteEventHandler) Type() string {
	return GroupDeleteEventType
}
