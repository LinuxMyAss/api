package api

// UserDeleteEventType holds the event type string for this event.
const UserDeleteEventType = "user_delete"

// UserDeleteEvent .
type UserDeleteEvent struct {
	ID string `json:"id"`
}

// Type returns the event's type.
func (event *UserDeleteEvent) Type() string {
	return UserDeleteEventType
}

// userDeleteEventHandler represents a UserDelete event handler.
type userDeleteEventHandler func(*Library, *UserDeleteEvent)

// New .
func (handler userDeleteEventHandler) New() interface{} {
	return &UserDeleteEvent{}
}

// Handle calls the underlying handler.
func (handler userDeleteEventHandler) Handle(library *Library, i interface{}) {
	if event, ok := i.(*UserDeleteEvent); ok {
		handler(library, event)
	}
}

// Type returns the event's type.
func (handler userDeleteEventHandler) Type() string {
	return UserDeleteEventType
}
