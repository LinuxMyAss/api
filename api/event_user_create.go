package api

// UserCreateEventType holds the event type string for this event.
const UserCreateEventType = "user_create"

// UserCreateEvent .
type UserCreateEvent struct {
	User *User `json:"user"`
}

// Type returns the event's type.
func (event *UserCreateEvent) Type() string {
	return UserCreateEventType
}

// userCreateEventHandler represents a UserCreate event handler.
type userCreateEventHandler func(*Library, *UserCreateEvent)

// New .
func (handler userCreateEventHandler) New() interface{} {
	return &UserCreateEvent{}
}

// Handle calls the underlying handler.
func (handler userCreateEventHandler) Handle(library *Library, i interface{}) {
	if event, ok := i.(*UserCreateEvent); ok {
		handler(library, event)
	}
}

// Type returns the event's type.
func (handler userCreateEventHandler) Type() string {
	return UserCreateEventType
}
