package api

// PunishmentCreateEventType holds the event type string for this event.
const PunishmentCreateEventType = "punishment_create"

// PunishmentCreateEvent .
type PunishmentCreateEvent struct {
	Punishment *Punishment `json:"id"`
}

// Type returns the event's type.
func (event *PunishmentCreateEvent) Type() string {
	return PunishmentCreateEventType
}

// punishmentCreateHandler represents a PunishmentCreate event handler.
type punishmentCreateEventHandler func(*Library, *PunishmentCreateEvent)

// New .
func (handler punishmentCreateEventHandler) New() interface{} {
	return &PunishmentCreateEvent{}
}

// Handle calls the underlying handler.
func (handler punishmentCreateEventHandler) Handle(library *Library, i interface{}) {
	if event, ok := i.(*PunishmentCreateEvent); ok {
		handler(library, event)
	}
}

// Type returns the event's type.
func (handler punishmentCreateEventHandler) Type() string {
	return PunishmentCreateEventType
}
