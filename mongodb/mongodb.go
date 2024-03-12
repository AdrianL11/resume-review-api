package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"resume-review-api/util/resume_ai_env"
	"time"
)

var (
	// This should be using sync.Once
	globalClient *mongo.Client
)

// getMongoClient – creates a session for mongoDB
func getMongoClient() (*mongo.Client, error) {
	serverSettings := resume_ai_env.GetSettingsForEnv()
	var err error

	opts := options.Client()

	if resume_ai_env.IsProd() {
		opts.ApplyURI("mongodb+srv://" + serverSettings.DBUsername + ":" + serverSettings.DBPassword + "@" + serverSettings.DBURL + "/?retryWrites=true&w=majority")
	} else {
		opts.ApplyURI("mongodb://localhost/?retryWrites=true&w=majority")
	}

	opts.SetConnectTimeout(30 * time.Second)

	// Check If Global is Set
	if globalClient != nil {
		// Check If Already Connected
		if err = globalClient.Ping(context.Background(), nil); err != nil {
			// Not Connected
			if globalClient, err = mongo.Connect(context.Background(), opts); err != nil {
				return globalClient, err
			}
		}
	} else {
		// No Global Set, Connect
		if globalClient, err = mongo.Connect(context.Background(), opts); err != nil {
			return globalClient, err
		}
	}

	if err = globalClient.Ping(context.Background(), nil); err != nil {
		return globalClient, err
	}

	return globalClient, err
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
	//defer func(client *mongo.Client, ctx context.Context) {
	//	_ = client.Disconnect(ctx)
	//}(client, ctx)

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
	//defer func(client *mongo.Client, ctx context.Context) {
	//	_ = client.Disconnect(ctx)
	//}(client, ctx)

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
	//defer func(client *mongo.Client, ctx context.Context) {
	//	_ = client.Disconnect(ctx)
	//}(client, ctx)

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
	//defer func(client *mongo.Client, ctx context.Context) {
	//	_ = client.Disconnect(ctx)
	//}(client, ctx)

	coll = client.Database(database).Collection(collection)
	if _, err = coll.UpdateOne(ctx, filter, bson.D{{"$set", update}}); err != nil {
		return err
	}

	return nil
}

// UpdateMany – updates multiple documents based on filter in mongoDB
// after retrieving client connection from getMongoClient
func UpdateMany(database string, collection string, filter interface{}, update interface{}) error {

	var err error
	var client *mongo.Client
	var ctx = context.Background()
	var coll *mongo.Collection

	if client, err = getMongoClient(); err != nil {
		return err
	}
	//defer func(client *mongo.Client, ctx context.Context) {
	//	_ = client.Disconnect(ctx)
	//}(client, ctx)

	coll = client.Database(database).Collection(collection)
	if _, err = coll.UpdateMany(ctx, filter, bson.D{{"$set", update}}); err != nil {
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
	//defer func(client *mongo.Client, ctx context.Context) {
	//	_ = client.Disconnect(ctx)
	//}(client, ctx)

	coll = client.Database(database).Collection(collection)

	if cursor, err = coll.Aggregate(ctx, pipeline); err != nil {
		return err
	}

	if err = cursor.All(ctx, results); err != nil {
		return err
	}

	return err
}
