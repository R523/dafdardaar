package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/r523/dafdardaar/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrIDNotFound indicates that given id does not exist on database.
	ErrIDNotFound = errors.New("given email does not exist")
	// ErrIDDuplicate indicates that given id is exists on database.
	ErrIDDuplicate = errors.New("given email exists")
)

// MongoURL communicate with users collection in MongoDB.
type MongoUser struct {
	DB *mongo.Database
}

// Collection is a name of the MongoDB collection for Users.
const Collection = "users"

// NewMongoUser creates new User store.
func NewMongoUser(db *mongo.Database) *MongoUser {
	return &MongoUser{
		DB: db,
	}
}

// Set saves given user in database.
func (s *MongoUser) Set(ctx context.Context, user model.User) error {
	users := s.DB.Collection(Collection)

	if _, err := users.InsertOne(ctx, user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrIDDuplicate
		}

		return fmt.Errorf("mongodb failed: %w", err)
	}

	return nil
}

// Get retrieves user of the given id if it exists.
func (s *MongoUser) Get(ctx context.Context, id string) (model.User, error) {
	record := s.DB.Collection(Collection).FindOne(ctx, bson.M{
		"id": id,
	})

	var user model.User
	if err := record.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, ErrIDNotFound
		}

		return user, fmt.Errorf("mongodb failed: %w", err)
	}

	return user, nil
}
