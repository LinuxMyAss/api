package api

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"strings"
	"time"
)

// PunishmentService is an interface for interfacing with Punishments.
type PunishmentService interface {
	New(context.Context, string, bson.ObjectId, string, bson.ObjectId, string, bool, string, int64) *Punishment
	GetByID(context.Context, string) (*Punishment, error)
	List(context.Context, map[string]interface{}) ([]Punishment, error)
	Create(context.Context, *Punishment) error
	Update(context.Context, *Punishment) error
	Delete(context.Context, string) error
	Paginate(context.Context, int, int, map[string]interface{}) ([]Punishment, error)
	Count(context.Context, map[string]interface{}) (int, error)
}

// PunishmentServiceImpl is an implementation for the PunishmentService interface.
type PunishmentServiceImpl struct {
	library *Library
}

// New attempts to create a new Punishment object.
func (service *PunishmentServiceImpl) New(ctx context.Context, server string, userID bson.ObjectId, address string, punisherID bson.ObjectId, reason string, silent bool, _type string, expiresAt int64) *Punishment {
	punishment := &Punishment{
		ID:         bson.NewObjectId(),
		Server:     server,
		UserID:     userID,
		Address:    address,
		PunisherID: punisherID,
		Reason:     reason,
		Silent:     silent,
		Type:       _type,
		ExpiresAt:  time.Unix(0, expiresAt*int64(time.Millisecond)),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return punishment
}

// GetByID attempts to get a punishment by using an id.
func (service *PunishmentServiceImpl) GetByID(ctx context.Context, id string) (*Punishment, error) {
	var punishment *Punishment
	err := service.library.Mongo.Punishment.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&punishment)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, err
	}

	return punishment, nil
}

// List punishments
func (service *PunishmentServiceImpl) List(ctx context.Context, filter map[string]interface{}) ([]Punishment, error) {
	var punishments []Punishment

	err := service.library.Mongo.Punishment.Find(filter).All(&punishments)
	if err != nil {
		return nil, err
	}

	return punishments, nil
}

// Create a punishment
func (service *PunishmentServiceImpl) Create(ctx context.Context, punishment *Punishment) error {
	service.library.EventManager.Call(&PunishmentCreateEvent{
		Punishment: punishment,
	})
	return service.library.Mongo.Punishment.Insert(&punishment)
}

// Update a punishment
func (service *PunishmentServiceImpl) Update(ctx context.Context, punishment *Punishment) error {
	service.library.EventManager.Call(&PunishmentUpdateEvent{
		Punishment: punishment,
	})
	return service.library.Mongo.Punishment.UpdateId(punishment.ID, &punishment)
}

// Delete a punishment
func (service *PunishmentServiceImpl) Delete(ctx context.Context, id string) error {
	service.library.EventManager.Call(&PunishmentDeleteEvent{
		ID: id,
	})
	return service.library.Mongo.Punishment.RemoveId(bson.ObjectIdHex(id))
}

// Paginate a list of punishments
func (service *PunishmentServiceImpl) Paginate(ctx context.Context, page int, perPage int, filter map[string]interface{}) ([]Punishment, error) {
	var punishments []Punishment

	err := service.library.Mongo.Punishment.Find(filter).Skip(perPage * (page - 1)).Limit(perPage).All(&punishments)
	if err != nil {
		return nil, err
	}

	return punishments, nil
}

// Count all punishments
func (service *PunishmentServiceImpl) Count(ctx context.Context, filter map[string]interface{}) (int, error) {
	return service.library.Mongo.Punishment.Find(filter).Count()
}

// Punishment represents a "egirls.me" punishment
type Punishment struct {
	ID           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Server       string        `json:"server" bson:"server"`
	UserID       bson.ObjectId `json:"userId" bson:"userId"`
	Address      string        `json:"address" bson:"address"`
	PunisherID   bson.ObjectId `json:"punisherId" bson:"punisherId"`
	Reason       string        `json:"reason" bson:"reason"`
	Silent       bool          `json:"silent" bson:"silent"`
	Type         string        `json:"type" bson:"type"`
	ExpiresAt    time.Time     `json:"expiresAt" bson:"expiresAt"`
	RemovedAt    time.Time     `json:"removedAt" bson:"removedAt"`
	RemovedBy    string        `json:"removedBy" bson:"removedBy"`
	RemoveReason string        `json:"removeReason" bson:"removeReason"`
	CreatedAt    time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt" bson:"updatedAt"`
}
