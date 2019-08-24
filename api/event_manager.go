package api

import (
	"encoding/json"
	"api/logger"
	"sync"
)

type redisEvent struct {
	Type  string      `json:"type"`
	Event interface{} `json:"event"`
}

// EventManager represents an "egirls.me" event manager.
type EventManager struct {
	Library    *Library
	handleLock sync.RWMutex
	handlers   map[string][]*eventHandlerInstance
}

// newEventManager will create a new event manager.
func newEventManager(library *Library) *EventManager {
	manager := &EventManager{
		Library: library,
	}

	return manager
}

// Register registers an event handler.
func (manager *EventManager) Register(i interface{}) {
	// Get the proper event handler.
	handler := getHandlerForInterface(i)

	// Check if the handler is nil
	if handler == nil {
		return
	}

	manager.handleLock.Lock()
	defer manager.handleLock.Unlock()

	// Create a new map if the existing one is nil.
	if manager.handlers == nil {
		manager.handlers = map[string][]*eventHandlerInstance{}
	}

	eventHandler := &eventHandlerInstance{handler}
	manager.handlers[handler.Type()] = append(manager.handlers[handler.Type()], eventHandler)
}

// Call calls registered event handlers for the event passed into the function.
func (manager *EventManager) Call(i interface{}) {
	eventType := getTypeFromInterface(i)

	// Check if there was no event type found.
	if len(eventType) == 0 {
		return
	}

	go func() {
		data, err := json.Marshal(redisEvent{
			Type:  eventType,
			Event: i,
		})
		if err != nil {
			logger.Errorw("[Events] Failed to json#Marshal event data.")
			return
		}

		manager.Library.Redis.Client.Publish("ikuta:access:events", data)
	}()

	// Check if no handlers for the event type are registered.
	if manager.handlers[eventType] == nil {
		return
	}

	manager.handleLock.RLock()
	defer manager.handleLock.RUnlock()

	// Loop through registered event handlers and call it's function.
	for _, handler := range manager.handlers[eventType] {
		handler.eventHandler.Handle(manager.Library, i)
	}
}

func getHandlerForInterface(handler interface{}) EventHandler {
	switch params := handler.(type) {

	case func(*Library, interface{}):
		return interfaceEventHandler(params)

	case func(*Library, *GroupCreateEvent):
		return groupCreateEventHandler(params)

	case func(*Library, *GroupDeleteEvent):
		return groupDeleteEventHandler(params)

	case func(*Library, *GroupUpdateEvent):
		return groupUpdateEventHandler(params)

	case func(*Library, *PunishmentCreateEvent):
		return punishmentCreateEventHandler(params)

	case func(*Library, *PunishmentDeleteEvent):
		return punishmentDeleteEventHandler(params)

	case func(*Library, *PunishmentUpdateEvent):
		return punishmentUpdateEventHandler(params)

	case func(*Library, *UserCreateEvent):
		return userCreateEventHandler(params)

	case func(*Library, *UserDeleteEvent):
		return userDeleteEventHandler(params)

	case func(*Library, *UserLoginEvent):
		return userLoginEventHandler(params)

	case func(*Library, *UserUpdateEvent):
		return userUpdateEventHandler(params)

	}
	return nil
}

func getTypeFromInterface(event interface{}) string {
	switch event.(type) {

	case *GroupCreateEvent:
		return GroupCreateEventType

	case *GroupDeleteEvent:
		return GroupDeleteEventType

	case *GroupUpdateEvent:
		return GroupUpdateEventType

	case *PunishmentCreateEvent:
		return PunishmentCreateEventType

	case *PunishmentDeleteEvent:
		return PunishmentDeleteEventType

	case *PunishmentUpdateEvent:
		return PunishmentUpdateEventType

	case *UserCreateEvent:
		return UserCreateEventType

	case *UserDeleteEvent:
		return UserDeleteEventType

	case *UserLoginEvent:
		return UserLoginEventType

	case *UserUpdateEvent:
		return UserUpdateEventType

	}

	return ""
}
