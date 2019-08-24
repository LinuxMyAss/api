package api

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"api/utils"
	"strings"
	"time"
)

// UserService is an interface for interfacing with Users.
type UserService interface {
	New(context.Context, string, string, string, bson.ObjectId) (*User, error)
	GetByID(context.Context, string) (*User, error)
	GetByUniqueID(context.Context, string) (*User, error)
	GetByEmail(context.Context, string) (*User, error)
	GetByToken(context.Context, string) (*User, error)
	List(context.Context, map[string]interface{}) ([]User, error)
	Create(context.Context, *User) error
	Update(context.Context, *User) error
	Delete(context.Context, string) error
	Paginate(context.Context, int, int, map[string]interface{}) ([]User, error)
	Count(context.Context, map[string]interface{}) (int, error)
}

// UserServiceImpl is an implementation for the UserService interface.
type UserServiceImpl struct {
	library *Library
}

// New attempts to create a new User object.
func (service *UserServiceImpl) New(ctx context.Context, uniqueID string, email string, password string, group bson.ObjectId) (*User, error) {
	user := &User{
		ID:               bson.NewObjectId(),
		UniqueID:         uniqueID,
		Email:            email,
		MessagingEnabled: true,
		MessagingSounds:  true,
		Group:            group,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if len(password) > 0 {
		err := user.SetPassword(password)
		if err != nil {
			return nil, err
		}
	} else {
		user.Password = ""
	}

	return user, nil
}

// GetByID attempts to get a user by using an id.
func (service *UserServiceImpl) GetByID(ctx context.Context, id string) (*User, error) {
	var user *User
	err := service.library.Mongo.User.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&user)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, err
	}

	return user, nil
}

// GetByUniqueID attempts to get a user by using an email address.
func (service *UserServiceImpl) GetByUniqueID(ctx context.Context, id string) (*User, error) {
	var user *User
	err := service.library.Mongo.User.Find(bson.M{"uniqueId": id}).One(&user)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, err
	}

	return user, nil
}

// GetByEmail attempts to get a user by using an email address.
func (service *UserServiceImpl) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user *User
	err := service.library.Mongo.User.Find(bson.M{"email": email}).One(&user)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, err
	}

	return user, nil
}

// GetByToken attempts to get a user by using a token.
func (service *UserServiceImpl) GetByToken(ctx context.Context, token string) (*User, error) {
	var user *User
	err := service.library.Mongo.User.Find(bson.M{"token": token}).One(&user)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, err
	}

	return user, nil
}

// List users
func (service *UserServiceImpl) List(ctx context.Context, filter map[string]interface{}) ([]User, error) {
	var users []User

	err := service.library.Mongo.User.Find(filter).All(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Create a user
func (service *UserServiceImpl) Create(ctx context.Context, user *User) error {
	service.library.EventManager.Call(&UserCreateEvent{
		User: user,
	})
	return service.library.Mongo.User.Insert(&user)
}

// Update a user
func (service *UserServiceImpl) Update(ctx context.Context, user *User) error {
	service.library.EventManager.Call(&UserUpdateEvent{
		User: user,
	})
	return service.library.Mongo.User.UpdateId(user.ID, &user)
}

// Delete a user
func (service *UserServiceImpl) Delete(ctx context.Context, id string) error {
	service.library.EventManager.Call(&UserDeleteEvent{
		ID: id,
	})
	return service.library.Mongo.User.RemoveId(bson.ObjectIdHex(id))
}

// Paginate a list of users
func (service *UserServiceImpl) Paginate(ctx context.Context, page int, perPage int, filter map[string]interface{}) ([]User, error) {
	var users []User

	err := service.library.Mongo.User.Find(filter).Skip(perPage * (page - 1)).Limit(perPage).All(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Count all users
func (service *UserServiceImpl) Count(ctx context.Context, filter map[string]interface{}) (int, error) {
	return service.library.Mongo.User.Find(filter).Count()
}

// User represents a "egirls.me" user
type User struct {
	ID               bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UniqueID         string        `json:"uniqueId" bson:"uniqueId"`
	Email            string        `json:"email" bson:"email"`
	Password         string        `json:"-" bson:"password"`
	Name             string        `json:"name" bson:"name"`
	Address          string        `json:"address" bson:"address"`
	Prefix           string        `json:"prefix" bson:"prefix"`
	Suffix           string        `json:"suffix" bson:"suffix"`
	Alts             []string      `json:"alts" bson:"alts"`
	Addresses        []string      `json:"addresses" bson:"addresses"`
	Permissions      []string      `json:"permissions" bson:"permissions"`
	Notes            []string      `json:"notes" bson:"notes"`
	Friends          []string      `json:"friends" bson:"friends"`
	Ignored          []string      `json:"ignored" bson:"ignored"`
	MessagingEnabled bool          `json:"messagingEnabled" bson:"messagingEnabled"`
	MessagingSounds  bool          `json:"messagingSounds" bson:"messagingSounds"`
	Group            bson.ObjectId `json:"group" bson:"group"`
	Token            string        `json:"-" bson:"token"`
	CreatedAt        time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// SetPassword hashes a raw password and updates the user's password.
func (user *User) SetPassword(password string) error {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user.Password = hash
	return nil
}

// VerifyPassword takes a password and the user's hashed password and verifies the password against the hashed password.
func (user *User) VerifyPassword(password string) (bool, error) {
	return utils.VerifyPassword(password, user.Password)
}

// IsRegistered returns a boolean based off of if the user is registered.
func (user *User) IsRegistered() bool {
	return len(user.Email) > 0 && len(user.Password) > 0
}

// IsRegistering returns a boolean based off of if the user is registering.
func (user *User) IsRegistering() bool {
	return !user.IsRegistered() && len(user.Token) > 0
}

// IsResetting returns a boolean based off of if the user is resetting their password.
func (user *User) IsResetting() bool {
	return user.IsRegistered() && len(user.Token) > 0
}
