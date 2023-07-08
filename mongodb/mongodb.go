package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

// getMongoClient – creates a session for mongoDB
func getMongoClient() (*mongo.Client, error) {

	var err error
	var client *mongo.Client

	opts := options.Client()
	opts.ApplyURI("mongodb+srv://" + os.Getenv("mongodb_username") + ":" + os.Getenv("mongodb_password") + "@" + os.Getenv("mongodb_url") + "/?retryWrites=true&w=majority")
	opts.SetConnectTimeout(30 * time.Second)

	if client, err = mongo.Connect(context.Background(), opts); err != nil {
		return client, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return client, err
	}

	return client, err
}

// NewDocument – creates a new document in the mongoDB
// after retrieving client connection from getMongoClient
func NewDocument(database string, collection string, document interface{}) (*mongo.InsertOneResult, error) {

	var err error
	var client *mongo.Client
	var ctx = context.Background()
	var coll *mongo.Collection
	var result *mongo.InsertOneResult

	if client, err = getMongoClient(); err != nil {
		return result, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		_ = client.Disconnect(ctx)
	}(client, ctx)

	coll = client.Database(database).Collection(collection)
	if result, err = coll.InsertOne(ctx, document); err != nil {
		return result, err
	}

	return result, err
}

// FindOne – searches for document based on filter in mongoDB
// after retrieving client connection from getMongoClient
func FindOne(database string, collection string, filter interface{}, decode interface{}) error {

	var err error
	var client *mongo.Client
	var ctx = context.Background()
	var coll *mongo.Collection

	if client, err = getMongoClient(); err != nil {
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		_ = client.Disconnect(ctx)
	}(client, ctx)

	coll = client.Database(database).Collection(collection)
	if err = coll.FindOne(ctx, filter).Decode(decode); err != nil {
		return err
	}

	return nil
}

func FindMany(database string, collection string, filter interface{}, decode interface{}) error {

	var err error
	var client *mongo.Client
	var ctx = context.Background()
	var coll *mongo.Collection

	if client, err = getMongoClient(); err != nil {
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		_ = client.Disconnect(ctx)
	}(client, ctx)

	// Find
	coll = client.Database(database).Collection(collection)
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return err
	}

	// Decode
	if err = cursor.All(ctx, decode); err != nil {
		return err
	}

	return err
}

// UpdateOne – updates a document based on filter in mongoDB
// after retrieving client connection from getMongoClient
func UpdateOne(database string, collection string, filter interface{}, update interface{}) error {

	var err error
	var client *mongo.Client
	var ctx = context.Background()
	var coll *mongo.Collection

	if client, err = getMongoClient(); err != nil {
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		_ = client.Disconnect(ctx)
	}(client, ctx)

	coll = client.Database(database).Collection(collection)
	if _, err = coll.UpdateOne(ctx, filter, bson.D{{"$set", update}}); err != nil {
		return err
	}

	return nil
}

// Aggregate – updates a document based on filter in mongoDB
// after retrieving client connection from getMongoClient
func Aggregate(database string, collection string, pipeline interface{}, results interface{}) error {

	var err error
	var client *mongo.Client
	var ctx = context.Background()
	var coll *mongo.Collection
	var cursor *mongo.Cursor

	if client, err = getMongoClient(); err != nil {
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		_ = client.Disconnect(ctx)
	}(client, ctx)

	coll = client.Database(database).Collection(collection)

	if cursor, err = coll.Aggregate(ctx, pipeline); err != nil {
		return err
	}

	if err = cursor.All(context.TODO(), results); err != nil {
		return err
	}

	return err
}
