package api

// PunishmentDeleteEventType holds the event type string for this event.
const PunishmentDeleteEventType = "punishment_delete"

// PunishmentDeleteEvent .
type PunishmentDeleteEvent struct {
	ID string `json:"id"`
}

// Type returns the event's type.
func (event *PunishmentDeleteEvent) Type() string {
	return PunishmentDeleteEventType
}

// punishmentDeleteEventHandler represents a PunishmentDelete event handler.
type punishmentDeleteEventHandler func(*Library, *PunishmentDeleteEvent)

// New .
func (handler punishmentDeleteEventHandler) New() interface{} {
	return &PunishmentDeleteEvent{}
}

// Handle calls the underlying handler.
func (handler punishmentDeleteEventHandler) Handle(library *Library, i interface{}) {
	if event, ok := i.(*PunishmentDeleteEvent); ok {
		handler(library, event)
	}
}

// Type returns the event's type.
func (handler punishmentDeleteEventHandler) Type() string {
	return PunishmentDeleteEventType
}
