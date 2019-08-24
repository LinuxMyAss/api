package api

// PunishmentUpdateEventType holds the event type string for this event.
const PunishmentUpdateEventType = "punishment_update"

// PunishmentUpdateEvent .
type PunishmentUpdateEvent struct {
	Punishment *Punishment `json:"id"`
}

// Type returns the event's type.
func (event *PunishmentUpdateEvent) Type() string {
	return PunishmentUpdateEventType
}

// punishmentUpdateHandler represents a PunishmentUpdate event handler.
type punishmentUpdateEventHandler func(*Library, *PunishmentUpdateEvent)

// New .
func (handler punishmentUpdateEventHandler) New() interface{} {
	return &PunishmentUpdateEvent{}
}

// Handle calls the underlying handler.
func (handler punishmentUpdateEventHandler) Handle(library *Library, i interface{}) {
	if event, ok := i.(*PunishmentUpdateEvent); ok {
		handler(library, event)
	}
}

// Type returns the event's type.
func (handler punishmentUpdateEventHandler) Type() string {
	return PunishmentUpdateEventType
}
