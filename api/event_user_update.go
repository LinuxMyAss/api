package api

// UserUpdateEventType holds the event type string for this event.
const UserUpdateEventType = "user_update"

// UserUpdateEvent .
type UserUpdateEvent struct {
	User *User `json:"user"`
}

// Type returns the event's type.
func (event *UserUpdateEvent) Type() string {
	return UserUpdateEventType
}

// userUpdateEventHandler represents a UserUpdate event handler.
type userUpdateEventHandler func(*Library, *UserUpdateEvent)

// New .
func (handler userUpdateEventHandler) New() interface{} {
	return &UserUpdateEvent{}
}

// Handle calls the underlying handler.
func (handler userUpdateEventHandler) Handle(library *Library, i interface{}) {
	if event, ok := i.(*UserUpdateEvent); ok {
		handler(library, event)
	}
}

// Type returns the event's type.
func (handler userUpdateEventHandler) Type() string {
	return UserUpdateEventType
}
