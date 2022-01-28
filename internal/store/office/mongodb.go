package home

import (
	"context"
	"errors"
	"fmt"

	"github.com/r523/dafdardaar/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrIDNotFound = errors.New("office id does not exist")
	ErrIDNotEmpty = errors.New("office id must be empty")
)

// MongoOffice communicate with office collection in MongoDB.
type MongoOffice struct {
	DB *mongo.Database
}

const (
	// Collection is a name of the MongoDB collection for offices.
	Collection = "offices"
)

// NewMongoHome creates new Home store.
func NewMongoOffice(db *mongo.Database) *MongoOffice {
	return &MongoOffice{
		DB: db,
	}
}

// Set saves given home in database and returns its id.
func (s *MongoOffice) Set(ctx context.Context, office *model.Office) error {
	if office.ID != "" {
		return ErrIDNotEmpty
	}

	office.ID = primitive.NewObjectID().Hex()

	offices := s.DB.Collection(Collection)

	if _, err := offices.InsertOne(ctx, office); err != nil {
		return fmt.Errorf("mongodb failed: %w", err)
	}

	return nil
}

// Get retrieves home of the given id if it exists.
func (s *MongoOffice) Get(ctx context.Context, id string) (model.Office, error) {
	record := s.DB.Collection(Collection).FindOne(ctx, bson.M{
		"_id": id,
	})

	var office model.Office
	if err := record.Decode(&office); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return office, ErrIDNotFound
		}

		return office, fmt.Errorf("mongodb failed: %w", err)
	}

	return office, nil
}
