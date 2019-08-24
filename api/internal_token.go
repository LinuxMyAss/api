package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"api/logger"
	"strings"
	"time"
)

// InternalTokenService is an interface for interfacing with Tokens.
type InternalTokenService interface {
	New(context.Context, string, map[string]bool) *InternalToken
	FromJWT(context.Context, string) (*InternalToken, error)
	GetByID(context.Context, string) (*InternalToken, error)
	List(context.Context, map[string]interface{}) ([]InternalToken, error)
	Create(context.Context, *InternalToken) error
	Delete(context.Context, string) error
	Paginate(context.Context, int, int, map[string]interface{}) ([]InternalToken, error)
	Count(context.Context, map[string]interface{}) (int, error)
}

// InternalTokenServiceImpl is an implementation for the InternalTokenService interface.
type InternalTokenServiceImpl struct {
	library *Library
}

// New attempts to create a new Token object.
func (service *InternalTokenServiceImpl) New(ctx context.Context, description string, permissions map[string]bool) *InternalToken {
	token := &InternalToken{
		ID:          bson.NewObjectId(),
		Description: description,
		Permissions: permissions,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(time.Hour * 720), // 1 month
	}

	return token
}

// FromJWT converts a "Json Web Token" string to a Token object.
func (service *InternalTokenServiceImpl) FromJWT(ctx context.Context, rawJwt string) (*InternalToken, error) {
	parsedJwt, err := jwt.Parse(rawJwt, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}

		return []byte(service.library.config.Secret), nil
	})

	if parsedJwt == nil {
		return nil, errors.New("parsed jwt is nil")
	}

	claims, ok := parsedJwt.Claims.(jwt.MapClaims)
	if !ok || !parsedJwt.Valid {
		return nil, err
	}

	if parsedJwt.Header["id"] == nil {
		return nil, errors.New("jwt is missing the \"id\" field")
	}

	if claims["aud"] == nil {
		return nil, errors.New("jwt is missing the \"aud\" field")
	}

	if parsedJwt.Header["permissions"] == nil {
		return nil, errors.New("jwt is missing the \"permissions\" field")
	}

	if claims["iat"] == nil {
		return nil, errors.New("jwt is missing the \"iat\" field")
	}

	/*if claims["exp"] == nil {
		return nil, errors.New("jwt is missing the \"exp\" field")
	}*/

	permissions := make(map[string]bool)

	for key, value := range parsedJwt.Header["permissions"].(map[string]interface{}) {
		switch value := value.(type) {
		case bool:
			permissions[key] = value
		}
	}

	token := &InternalToken{
		ID:          bson.ObjectIdHex(parsedJwt.Header["id"].(string)),
		Description: claims["aud"].(string),
		Permissions: permissions,
		CreatedAt:   time.Unix(int64(claims["iat"].(float64)), 0),
		//ExpiresAt:   time.Unix(int64(claims["exp"].(float64)), 0),
	}

	return token, nil
}

// GetByID attempts to get a token by using an id.
func (service *InternalTokenServiceImpl) GetByID(ctx context.Context, id string) (*InternalToken, error) {
	var token *InternalToken

	// Attempt to get the token from redis.
	result, err := service.library.Redis.Client.Get(fmt.Sprintf("ikuta:access:token:%s", id)).Result()
	if err != nil {
		// Check if the error is telling us the value doesn't exist.
		if err.Error() == "redis: nil" {
			// Pull data from mongo.
			err := service.library.Mongo.InternalToken.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&token)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					return nil, nil
				}
				return nil, err
			}

			go func() {
				// Convert the token to a JSON string.
				data, err := json.Marshal(token)
				if err != nil {
					logger.Errorw("[Redis] (token.go) Failed to json#Marshal object.", logger.Err(err))
					return
				}

				// Insert the token into Redis.
				err = service.library.Redis.Client.Set(fmt.Sprintf("ikuta:access:token:%s", token.ID), data, 10*time.Minute).Err()
				if err != nil {
					logger.Errorw("[Redis] (token.go) Failed to insert object.", logger.Err(err))
					return
				}
			}()

			// Return the token from MongoDB.
			return token, nil
		}
		return nil, err
	}

	if len(result) < 1 {
		return nil, errors.New("empty result returned from cache")
	}

	err = json.Unmarshal([]byte(result), token)
	return token, err
}

// List all tokens matching a filter
func (service *InternalTokenServiceImpl) List(ctx context.Context, filter map[string]interface{}) ([]InternalToken, error) {
	var tokens []InternalToken

	err := service.library.Mongo.InternalToken.Find(filter).All(&tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Create a token
func (service *InternalTokenServiceImpl) Create(ctx context.Context, token *InternalToken) error {
	/*service.library.EventManager.Call(&InternalTokenCreateEvent{
		InternalToken: token,
	})*/

	go func() {
		// Convert the token to a JSON string.
		data, err := json.Marshal(token)
		if err != nil {
			logger.Errorw("[Redis] (token.go) Failed to json#Marshal object.", logger.Err(err))
			return
		}

		// Insert the token into Redis.
		err = service.library.Redis.Client.Set(fmt.Sprintf("ikuta:access:token:%s", token.ID), data, 10*time.Minute).Err()
		if err != nil {
			logger.Errorw("[Redis] (token.go) Failed to insert object.", logger.Err(err))
			return
		}
	}()

	// Insert the token into MongoDB.
	return service.library.Mongo.InternalToken.Insert(&token)
}

// Delete a token
func (service *InternalTokenServiceImpl) Delete(ctx context.Context, id string) error {
	/*service.library.EventManager.Call(&InternalTokenDeleteEvent{
		ID: id,
	})*/

	go func() {
		// Delete the token from Redis.
		err := service.library.Redis.Client.Del(fmt.Sprintf("ikuta:access:token:%s", id)).Err()
		if err != nil {
			logger.Errorw("[Redis] (token.go) Failed to delete object.", logger.Err(err))
			return
		}
	}()

	return service.library.Mongo.InternalToken.RemoveId(bson.ObjectIdHex(id))
}

// Paginate a list of tokens
func (service *InternalTokenServiceImpl) Paginate(ctx context.Context, page int, perPage int, filter map[string]interface{}) ([]InternalToken, error) {
	var tokens []InternalToken

	err := service.library.Mongo.InternalToken.Find(filter).Skip(perPage * (page - 1)).Limit(perPage).All(&tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Count all tokens
func (service *InternalTokenServiceImpl) Count(ctx context.Context, filter map[string]interface{}) (int, error) {
	return service.library.Mongo.InternalToken.Find(filter).Count()
}

// InternalToken represents a "egirls.me" token
type InternalToken struct {
	ID          bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	Description string          `json:"description" bson:"description"`
	Permissions map[string]bool `json:"permissions" bson:"permissions"`
	CreatedAt   time.Time       `json:"createdAt" bson:"createdAt"`
	ExpiresAt   time.Time       `json:"expiresAt" bson:"expiresAt"`
}

// JWT generates a Json Web Token using the data from the Token object.
func (token *InternalToken) JWT(lib *Library) (string, error) {
	jsonToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Audience: token.Description,
		//ExpiresAt: token.ExpiresAt.Unix(), DISABLED: We want the token to last forever.
		Issuer:   "egirls.me",
		IssuedAt: token.CreatedAt.Unix(),
		Subject:  "access_token",
	})
	jsonToken.Header["id"] = token.ID.Hex()
	jsonToken.Header["permissions"] = token.Permissions

	signedJwt, err := jsonToken.SignedString([]byte(lib.config.Secret))
	if err != nil {
		logger.Errorw("[Backend] Failed to sign JWT.", logger.Err(err))
	}

	return signedJwt, err
}
