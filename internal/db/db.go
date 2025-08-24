package db

/*
handle SQlite3 /  MongoDB (preferred)

imports:
MongoDB imports :
"context"
"time"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
"go.mongodb.org/mongo-driver/bson"
utils package for the logging/config

functions:

initDB(path string) (*sql.DB, error)
CloseDB(db *sql.DB) to close the DB
other required mongoDB specific functions
ConnectMongo(uri string) (*mongo.Client, error)
InsertDocument(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error)
FindOneDocument(ctx context.Context, collection string, filter interface{}, result interface{}) error
FindManyDocuments(ctx context.Context, collection string, filter interface{}, results interface{}) error *(NOT DONE)
UpdateDocument(ctx context.Context, collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
DeleteDocument(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error)
(for CRUD operations)
another file called models.go may be required
*/

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// for typed IDs, use bson.ObjectID
// here any is fine, so we will use it
// any is similar to auto_inc from mysql
type ToDoList_DB struct {
	ID any `bson:"_id,omitempty"`
	Title string `bson:"title"`
	Done bool `bson:"done"`
}

func ConnectDB(uri string) (*mongo.Database, *mongo.Client) {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// Client creates a new ClientOptions instance.
	// ApplyURI parses the given URI and sets options accordingly.
	// ApplyURI is a reciever (obj.method type) on *ClientOptions
	// SetServerAPIOptions specifies a ServerAPIOptions instance 
	// used to configure the API version sent to the server
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
	panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
	panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client.Database("Todo_DB"), client
}

func InsertTodo (db *mongo.Database, todo ToDoList_DB) (*mongo.InsertOneResult, error) {
	collection := db.Collection("Todo_DB")
	return collection.InsertOne(context.TODO(), todo)
} 

func GetTodoByID (db *mongo.Database, id any) (ToDoList_DB, error) {
	collection := db.Collection("Todo_DB")
	var result ToDoList_DB
	filter := bson.M{"_id":id}

	// FindOne executes a find command and returns a SingleResult for one document in the collection.
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// returned. It cannot be nil.
	// Decode will unmarshal the document represented by this SingleResult into v.
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

func UpdateTodo (db *mongo.Database, id any, update_value ToDoList_DB) (*mongo.UpdateResult, error) {
	collection := db.Collection("Todo_DB")
	filter := bson.M{"_id":id}
	update := bson.M{"$set":update_value}

	// UpdateOne executes an update command to update at most one document in the collection.
	// The update parameter must be a document containing update operators
	return collection.UpdateOne(context.TODO(), filter, update)
}

func DeleteTodo (db *mongo.Database, id any) (*mongo.DeleteResult, error) {
	collection := db.Collection("Todo_DB")
	filter := bson.M{"_id":id}
	return collection.DeleteOne(context.TODO(), filter)
}