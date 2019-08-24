package api

// UserLoginEventType holds the event type string for this event.
const UserLoginEventType = "user_login"

// UserLoginEvent .
type UserLoginEvent struct {
	User  *User  `json:"user"`
	Token *Token `json:"token"`
}

// Type returns the event's type.
func (event *UserLoginEvent) Type() string {
	return UserLoginEventType
}

// userLoginEventHandler represents a UserLogin event handler.
type userLoginEventHandler func(*Library, *UserLoginEvent)

// New .
func (handler userLoginEventHandler) New() interface{} {
	return &UserLoginEvent{}
}

// Handle calls the underlying handler.
func (handler userLoginEventHandler) Handle(library *Library, i interface{}) {
	if event, ok := i.(*UserLoginEvent); ok {
		handler(library, event)
	}
}

// Type returns the event's type.
func (handler userLoginEventHandler) Type() string {
	return UserLoginEventType
}
