package mongo

import (
	"context"
	"errors"
	"github.com/arraisi/demogo/config"
	"github.com/arraisi/demogo/pkg/utils"
	"log"
	"time"

	mongoTracer "gopkg.in/DataDog/dd-trace-go.v1/contrib/go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	bsonPrimitive "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBHandler struct {
	Client                 *mongo.Client
	registeredCollections  map[string]*mongo.Collection
	registeredConstructors map[string]func() interface{}
	Database               *mongo.Database
}

const (
	DBTimeout int = 5
)

var (
	collectionIsNotRegistered = errors.New("Collection is not registered")
)

// Connect connects to a mongoDB instance and save the client connection
// example registeredCollections : m["audit_log"]=func() interface{} {return &mongoEntity.AuditLog{}}
func (m *DBHandler) Connect(mongoDBAccount *config.MongoDBAccount, registeredCollections map[string]func() interface{}) error {
	uri := utils.GenerateMongoDBURL(mongoDBAccount)
	opts := options.Client().ApplyURI(uri)
	// local
	//opts := options.Client().ApplyURI(uri).SetDirect(true)

	opts.SetMonitor(mongoTracer.NewMonitor(mongoTracer.WithServiceName("ims-mongo")))
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(DBTimeout)*time.Second)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to open connection to MongoDB! Error : %s\n", err.Error())
		return err
	}

	m.Client = client
	m.Database = client.Database(mongoDBAccount.DBName)

	collectionMap := make(map[string]*mongo.Collection)
	constructorMap := make(map[string]func() interface{})

	for collectionName, constructor := range registeredCollections {
		collectionMap[collectionName] = m.Database.Collection(collectionName)
		constructorMap[collectionName] = constructor
	}

	m.registeredCollections = collectionMap
	m.registeredConstructors = constructorMap

	//manually set constructor for each collection
	//constructorMap := make(map[string]func() interface{})
	// constructorMap["audit_log"] = func() interface{} {
	// 	return &mongoEntity.AuditLog{}
	// }

	// constructorMap["request_log"] = func() interface{} {
	// 	return &mongoEntity.RequestLog{}
	// }

	return nil
}

// Disconnect disconnects client mongo db
func (m *DBHandler) Disconnect(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

// InsertOne insert a document (struct type with bson tag) into a collection
func (m *DBHandler) InsertOne(ctx context.Context, collectionName string, document interface{}) error {
	if !m.isCollectionAlreadyRegistered(collectionName) {
		return collectionIsNotRegistered
	}

	collection := m.getCollection(collectionName)

	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}

	return nil
}

// UpdateOne update a document (struct type with bson tag) of a collection
func (m *DBHandler) UpdateOne(ctx context.Context, collectionName string, objectID bsonPrimitive.ObjectID, document interface{}) error {
	if !m.isCollectionAlreadyRegistered(collectionName) {
		return collectionIsNotRegistered
	}

	collection := m.getCollection(collectionName)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, document)
	if err != nil {
		return err
	}

	return nil
}

// FindFilter finds documents that matches the filter (bson.M)
func (m *DBHandler) FindFilter(ctx context.Context, collectionName string, filter bson.M) ([]interface{}, error) {
	if !m.isCollectionAlreadyRegistered(collectionName) {
		return nil, collectionIsNotRegistered
	}

	collection := m.getCollection(collectionName)
	constructor := m.getConstructor(collectionName)

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var result []interface{}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		targetDecode := constructor()
		err := cursor.Decode(&targetDecode)
		if err != nil {
			return nil, err
		}

		result = append(result, targetDecode)
	}

	return result, nil
}

// FindOne finds a document that matches by a key (field name) and a value (field value)
func (m *DBHandler) FindOne(ctx context.Context, collectionName, fieldName, fieldValue string) (*mongo.SingleResult, error) {
	if !m.isCollectionAlreadyRegistered(collectionName) {
		return nil, collectionIsNotRegistered
	}

	collection := m.getCollection(collectionName)
	opts := options.FindOneOptions{}
	opts.SetSort(bson.D{{"_id", -1}})
	return collection.FindOne(ctx, bson.M{fieldName: fieldValue}, &opts), nil
}

// FindOneByID finds a document that matches by object ID
func (m *DBHandler) FindOneByID(ctx context.Context, collectionName string, objectID bsonPrimitive.ObjectID) (*mongo.SingleResult, error) {
	if !m.isCollectionAlreadyRegistered(collectionName) {
		return nil, collectionIsNotRegistered
	}

	collection := m.getCollection(collectionName)
	return collection.FindOne(ctx, bson.M{"_id": objectID}), nil
}

// FindOne finds a document that matches the filter (bson.M)
func (m *DBHandler) FindOneFilter(ctx context.Context, collectionName string, filter bson.M) (interface{}, error) {
	if !m.isCollectionAlreadyRegistered(collectionName) {
		return nil, collectionIsNotRegistered
	}

	collection := m.getCollection(collectionName)
	constructor := m.getConstructor(collectionName)

	rawResult := collection.FindOne(ctx, filter)
	targetDecode := constructor()
	err := rawResult.Decode(&targetDecode)
	if err != nil {
		return nil, err
	}
	return targetDecode, nil
}

// FindFirst finds the first document in a collection sorted by _id
func (m *DBHandler) FindFirst(ctx context.Context, collectionName string) (*mongo.SingleResult, error) {
	if !m.isCollectionAlreadyRegistered(collectionName) {
		return nil, collectionIsNotRegistered
	}

	collection := m.getCollection(collectionName)
	opts := options.FindOneOptions{}
	opts.SetSort(bson.D{{"_id", -1}})
	return collection.FindOne(ctx, bson.M{}, &opts), nil
}

// isCollectionAlreadyRegistered checks if a collection is already registered or not
func (m *DBHandler) isCollectionAlreadyRegistered(collectionName string) bool {
	_, ok := m.registeredCollections[collectionName]
	return ok
}

// getCollection from a registered collection (returns nil if collection is not found)
func (m *DBHandler) getCollection(collectionName string) *mongo.Collection {
	collection, ok := m.registeredCollections[collectionName]
	if ok {
		return collection
	}
	return nil
}

// getConstructor from a registered constructor
func (m *DBHandler) getConstructor(collectionName string) func() interface{} {
	constructor, ok := m.registeredConstructors[collectionName]
	if ok {
		return constructor
	}
	return nil
}

// InsertOneReturning insert a document (struct type with bson tag) into a collection returning result
func (m *DBHandler) InsertOneReturning(ctx context.Context, collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	if !m.isCollectionAlreadyRegistered(collectionName) {
		return &mongo.InsertOneResult{}, collectionIsNotRegistered
	}

	collection := m.getCollection(collectionName)

	return collection.InsertOne(ctx, document)
}

// DeleteOne delete a document (struct type with bson tag) into a collection returning result
func (m *DBHandler) DeleteOne(ctx context.Context, collectionName string, document interface{}) (*mongo.DeleteResult, error) {
	if !m.isCollectionAlreadyRegistered(collectionName) {
		return &mongo.DeleteResult{}, collectionIsNotRegistered
	}

	collection := m.getCollection(collectionName)

	return collection.DeleteOne(ctx, document)
}

// DeleteFilter delete filter document (struct type with bson tag) into a collection returning result
func (m *DBHandler) DeleteFilter(ctx context.Context, collectionName string, document interface{}) (*mongo.DeleteResult, error) {
	if !m.isCollectionAlreadyRegistered(collectionName) {
		return &mongo.DeleteResult{}, collectionIsNotRegistered
	}

	collection := m.getCollection(collectionName)
	return collection.DeleteMany(ctx, document)
}
