package api

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"strings"
	"time"
)

// TicketService is an interface for interfacing with Tickets.
type TicketService interface {
	New(context.Context, bson.ObjectId, string, bson.ObjectId) *Ticket
	GetByID(context.Context, string) (*Ticket, error)
	List(context.Context, map[string]interface{}) ([]Ticket, error)
	Create(context.Context, *Ticket) error
	Update(context.Context, *Ticket) error
	Delete(context.Context, string) error
	Paginate(context.Context, int, int, map[string]interface{}) ([]Ticket, error)
	Count(context.Context, map[string]interface{}) (int, error)
}

// TicketServiceImpl is an implementation for the TicketService interface.
type TicketServiceImpl struct {
	library *Library
}

// New attempts to create a new Ticket object.
func (service *TicketServiceImpl) New(ctx context.Context, user bson.ObjectId, body string, category bson.ObjectId) *Ticket {
	ticket := &Ticket{
		ID:        bson.NewObjectId(),
		User:      user,
		Body:      body,
		Category:  category,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return ticket
}

// GetByID attempts to get a ticket by using an id.
func (service *TicketServiceImpl) GetByID(ctx context.Context, id string) (*Ticket, error) {
	var ticket *Ticket
	err := service.library.Mongo.Ticket.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&ticket)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, err
	}

	return ticket, nil
}

// List tickets
func (service *TicketServiceImpl) List(ctx context.Context, filter map[string]interface{}) ([]Ticket, error) {
	var tickets []Ticket

	err := service.library.Mongo.Ticket.Find(filter).All(&tickets)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

// Create a ticket
func (service *TicketServiceImpl) Create(ctx context.Context, ticket *Ticket) error {
	return service.library.Mongo.Ticket.Insert(&ticket)
}

// Update a ticket
func (service *TicketServiceImpl) Update(ctx context.Context, ticket *Ticket) error {
	return service.library.Mongo.Ticket.UpdateId(ticket.ID, &ticket)
}

// Delete a ticket
func (service *TicketServiceImpl) Delete(ctx context.Context, id string) error {
	return service.library.Mongo.Ticket.RemoveId(bson.ObjectIdHex(id))
}

// Paginate a list of tickets
func (service *TicketServiceImpl) Paginate(ctx context.Context, page int, perPage int, filter map[string]interface{}) ([]Ticket, error) {
	var tickets []Ticket

	err := service.library.Mongo.Ticket.Find(filter).Skip(perPage * (page - 1)).Limit(perPage).All(&tickets)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

// Count all tickets
func (service *TicketServiceImpl) Count(ctx context.Context, filter map[string]interface{}) (int, error) {
	return service.library.Mongo.Ticket.Find(filter).Count()
}

// Ticket represents a "egirls.me" ticket
type Ticket struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User      bson.ObjectId `json:"user" bson:"user"`
	Body      string        `json:"body" bson:"body"`
	Category  bson.ObjectId `json:"category" bson:"category"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}
