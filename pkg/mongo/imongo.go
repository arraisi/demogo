package mongo

//go:generate mockgen --destination=../mocks/mock_imongo.go --package=mocks --source=imongo.go
import (
	"context"
	"demogo/config"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	bsonPrimitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// IMongoDB is the interface that defines behavior of MongoDB Client
type IMongoDB interface {
	Connect(mongoDBAccount *config.MongoDBAccount, registeredCollections map[string]func() interface{}) error
	Disconnect(ctx context.Context) error
	InsertOne(ctx context.Context, collectionName string, document interface{}) error
	UpdateOne(ctx context.Context, collectionName string, objectID bsonPrimitive.ObjectID, document interface{}) error
	FindFilter(ctx context.Context, collectionName string, filter bson.M) ([]interface{}, error)
	FindOne(ctx context.Context, collectionName, fieldName, fieldValue string) (*mongo.SingleResult, error)
	FindOneByID(ctx context.Context, collectionName string, objectID bsonPrimitive.ObjectID) (*mongo.SingleResult, error)
	FindOneFilter(ctx context.Context, collectionName string, filter bson.M) (interface{}, error)
	FindFirst(ctx context.Context, collectionName string) (*mongo.SingleResult, error)
	InsertOneReturning(ctx context.Context, collectionName string, document interface{}) (*mongo.InsertOneResult, error)
	DeleteOne(ctx context.Context, collectionName string, document interface{}) (*mongo.DeleteResult, error)
	DeleteFilter(ctx context.Context, collectionName string, document interface{}) (*mongo.DeleteResult, error)
}
