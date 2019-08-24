package api

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"strings"
)

// GroupService is an interface for interfacing with Groups.
type GroupService interface {
	New(context.Context, string, bool) *Group
	GetByID(context.Context, string) (*Group, error)
	List(context.Context, map[string]interface{}) ([]Group, error)
	Create(context.Context, *Group) error
	Update(context.Context, *Group) error
	Delete(context.Context, string) error
	Count(context.Context) (int, error)
}

// GroupServiceImpl is an implementation for the GroupService interface.
type GroupServiceImpl struct {
	library *Library
}

// New attempts to create a new User object.
func (service *GroupServiceImpl) New(ctx context.Context, name string, protected bool) *Group {
	group := &Group{
		ID:        bson.NewObjectId(),
		Name:      name,
		Protected: protected,
	}
	return group
}

// GetByID attempts to get a group by using an id.
func (service *GroupServiceImpl) GetByID(ctx context.Context, id string) (*Group, error) {
	var group *Group
	err := service.library.Mongo.Group.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&group)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, err
	}

	return group, nil
}

// List groups
func (service *GroupServiceImpl) List(ctx context.Context, filter map[string]interface{}) ([]Group, error) {
	var groups []Group

	err := service.library.Mongo.Group.Find(filter).All(&groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// Create a group
func (service *GroupServiceImpl) Create(ctx context.Context, group *Group) error {
	service.library.EventManager.Call(&GroupCreateEvent{
		Group: group,
	})
	return service.library.Mongo.Group.Insert(&group)
}

// Update a group
func (service *GroupServiceImpl) Update(ctx context.Context, group *Group) error {
	service.library.EventManager.Call(&GroupUpdateEvent{
		Group: group,
	})
	return service.library.Mongo.Group.UpdateId(group.ID, &group)
}

// Delete a group
func (service *GroupServiceImpl) Delete(ctx context.Context, id string) error {
	service.library.EventManager.Call(&GroupDeleteEvent{
		ID: id,
	})
	return service.library.Mongo.Group.RemoveId(bson.ObjectIdHex(id))
}

// Count all groups
func (service *GroupServiceImpl) Count(ctx context.Context) (int, error) {
	return service.library.Mongo.Group.Count()
}

// Group represents a "egirls.me" group
type Group struct {
	ID             bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	Name           string          `json:"name" bson:"name"`
	Prefix         string          `json:"prefix" bson:"prefix"`
	Suffix         string          `json:"suffix" bson:"suffix"`
	Color          string          `json:"color" bson:"color"`
	Permissions    []string        `json:"permissions" bson:"permissions"`
	WebPermissions map[string]bool `json:"webPermissions" bson:"webPermissions"`
	SortID         int             `json:"sortId" bson:"sortId"`
	Protected      bool            `json:"protected" bson:"protected"`
}

// HasWebPermission returns true if the group has the specified permission.
func (group *Group) HasWebPermission(permission string) bool {
	if group.WebPermissions == nil {
		return false
	}

	if len(group.WebPermissions) < 1 {
		return false
	}

	if group.WebPermissions["root"] == true {
		return true
	}

	return group.WebPermissions[permission] == true
}

// IsProtected returns true if the group is protected.
func (group *Group) IsProtected() bool {
	return group.Protected
}
