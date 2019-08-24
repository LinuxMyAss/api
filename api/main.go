package api

import (
	"api/backend"
)

// Library .
type Library struct {
	config        Config
	Mongo         backend.MongoDriver
	Redis         backend.RedisDriver
	EventManager  *EventManager
	Group         GroupService
	InternalToken InternalTokenService
	Punishment    PunishmentService
	Ticket        TicketService
	Token         TokenService
	User          UserService
}

// Config .
type Config struct {
	Secret string `json:"secret"`

	MongoDB struct {
		Active   bool   `json:"active"`
		URI      string `json:"uri"`
		Database string `json:"database"`
	} `json:"mongodb"`

	Redis struct {
		Active   bool   `json:"active"`
		URI      string `json:"uri"`
		Password string `json:"password"`
		Database int    `json:"database"`
	} `json:"redis"`
}

// New .
func New(config Config) (*Library, error) {
	mongo := backend.MongoDriver{}
	if config.MongoDB.Active {
		err := mongo.Connect(config.MongoDB.URI, config.MongoDB.Database)
		if err != nil {
			return nil, err
		}
	}

	redis := backend.RedisDriver{}
	if config.Redis.Active {
		err := redis.Connect(config.Redis.URI, config.Redis.Password, config.Redis.Database)
		if err != nil {
			return nil, err
		}
	}

	library := &Library{
		config: config,
		Mongo:  mongo,
		Redis:  redis,
	}
	library.EventManager = newEventManager(library)
	library.Group = &GroupServiceImpl{library: library}
	library.InternalToken = &InternalTokenServiceImpl{library: library}
	library.Punishment = &PunishmentServiceImpl{library: library}
	library.Ticket = &TicketServiceImpl{library: library}
	library.Token = &TokenServiceImpl{library: library}
	library.User = &UserServiceImpl{library: library}

	return library, nil
}
