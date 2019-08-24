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

// TokenService is an interface for interfacing with Tokens.
type TokenService interface {
	New(context.Context, bson.ObjectId, string, string, map[string]bool) *Token
	FromJWT(context.Context, string) (*Token, error)
	GetByID(context.Context, string) (*Token, error)
	List(context.Context, map[string]interface{}) ([]Token, error)
	Create(context.Context, *Token) error
	Delete(context.Context, string) error
	Paginate(context.Context, int, int, map[string]interface{}) ([]Token, error)
	Count(context.Context, map[string]interface{}) (int, error)
}

// TokenServiceImpl is an implementation for the TokenService interface.
type TokenServiceImpl struct {
	library *Library
}

// New attempts to create a new Token object.
func (service *TokenServiceImpl) New(ctx context.Context, user bson.ObjectId, address string, userAgent string, permissions map[string]bool) *Token {
	token := &Token{
		ID:          bson.NewObjectId(),
		User:        user,
		Address:     address,
		UserAgent:   userAgent,
		Permissions: permissions,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(time.Hour * 720), // 1 month
	}

	return token
}

// FromJWT converts a "Json Web Token" string to a Token object.
func (service *TokenServiceImpl) FromJWT(ctx context.Context, rawJwt string) (*Token, error) {
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

	if parsedJwt.Header["address"] == nil {
		return nil, errors.New("jwt is missing the \"address\" field")
	}

	if parsedJwt.Header["userAgent"] == nil {
		return nil, errors.New("jwt is missing the \"userAgent\" field")
	}

	if parsedJwt.Header["permissions"] == nil {
		return nil, errors.New("jwt is missing the \"permissions\" field")
	}

	if claims["iat"] == nil {
		return nil, errors.New("jwt is missing the \"iat\" field")
	}

	if claims["exp"] == nil {
		return nil, errors.New("jwt is missing the \"exp\" field")
	}

	permissions := make(map[string]bool)

	for key, value := range parsedJwt.Header["permissions"].(map[string]interface{}) {
		switch value := value.(type) {
		case bool:
			permissions[key] = value
		}
	}

	token := &Token{
		ID:          bson.ObjectIdHex(parsedJwt.Header["id"].(string)),
		User:        bson.ObjectIdHex(claims["aud"].(string)),
		Address:     parsedJwt.Header["address"].(string),
		UserAgent:   parsedJwt.Header["userAgent"].(string),
		Permissions: permissions,
		CreatedAt:   time.Unix(int64(claims["iat"].(float64)), 0),
		ExpiresAt:   time.Unix(int64(claims["exp"].(float64)), 0),
	}

	return token, nil
}

// GetByID attempts to get a token by using an id.
func (service *TokenServiceImpl) GetByID(ctx context.Context, id string) (*Token, error) {
	var token *Token

	// Attempt to get the token from redis.
	result, err := service.library.Redis.Client.Get(fmt.Sprintf("ikuta:access:token:%s", id)).Result()
	if err != nil {
		// Check if the error is telling us the value doesn't exist.
		if err.Error() == "redis: nil" {
			// Pull data from mongo.
			err := service.library.Mongo.Token.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&token)
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
func (service *TokenServiceImpl) List(ctx context.Context, filter map[string]interface{}) ([]Token, error) {
	var tokens []Token

	err := service.library.Mongo.Token.Find(filter).All(&tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Create a token
func (service *TokenServiceImpl) Create(ctx context.Context, token *Token) error {
	/*service.library.EventManager.Call(&TokenCreateEvent{
		Token: token,
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
	return service.library.Mongo.Token.Insert(&token)
}

// Delete a token
func (service *TokenServiceImpl) Delete(ctx context.Context, id string) error {
	/*service.library.EventManager.Call(&TokenDeleteEvent{
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

	return service.library.Mongo.Token.RemoveId(bson.ObjectIdHex(id))
}

// Paginate a list of tokens
func (service *TokenServiceImpl) Paginate(ctx context.Context, page int, perPage int, filter map[string]interface{}) ([]Token, error) {
	var tokens []Token

	err := service.library.Mongo.Token.Find(filter).Skip(perPage * (page - 1)).Limit(perPage).All(&tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Count all tokens
func (service *TokenServiceImpl) Count(ctx context.Context, filter map[string]interface{}) (int, error) {
	return service.library.Mongo.Token.Find(filter).Count()
}

// Token represents a "egirls.me" token
type Token struct {
	ID          bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	User        bson.ObjectId   `json:"user" bson:"user"`
	Address     string          `json:"address" bson:"address"`
	UserAgent   string          `json:"userAgent" bson:"userAgent"`
	Permissions map[string]bool `json:"permissions" bson:"permissions"`
	CreatedAt   time.Time       `json:"createdAt" bson:"createdAt"`
	ExpiresAt   time.Time       `json:"expiresAt" bson:"expiresAt"`
}

// JWT generates a Json Web Token using the data from the Token object.
func (token *Token) JWT(lib *Library) (string, error) {
	jsonToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Audience:  token.User.Hex(),
		ExpiresAt: token.ExpiresAt.Unix(),
		Issuer:    "egirls.me",
		IssuedAt:  token.CreatedAt.Unix(),
		Subject:   "access_token",
	})
	jsonToken.Header["id"] = token.ID.Hex()
	jsonToken.Header["address"] = token.Address
	jsonToken.Header["userAgent"] = token.UserAgent
	jsonToken.Header["permissions"] = token.Permissions

	signedJwt, err := jsonToken.SignedString([]byte(lib.config.Secret))
	if err != nil {
		logger.Errorw("[Backend] Failed to sign JWT.", logger.Err(err))
	}

	return signedJwt, err
}
