package db

/*
handle  MongoDB (preferred)

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
	"github.com/ayush-bothra/backend-optimizer/internal/models"
)

// Connection specific
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

	return client.Database("app_db"), client
}


// Todo_DB specific
func InsertTodo(ctx context.Context, db *mongo.Database, todo models.ToDoList_DB) (*mongo.InsertOneResult, error) {
	collection := db.Collection("todos")
	return collection.InsertOne(ctx, todo)
}

func GetTodoByID(ctx context.Context, db *mongo.Database, id any) (models.ToDoList_DB, error) {
	collection := db.Collection("todos")
	var result models.ToDoList_DB
	filter := bson.M{"_id": id}

	// FindOne executes a find command and returns a SingleResult for one document in the collection.
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// returned. It cannot be nil.
	// Decode will unmarshal the document represented by this SingleResult into v.
	err := collection.FindOne(ctx, filter).Decode(&result)
	return result, err
}

func UpdateTodo(ctx context.Context, db *mongo.Database, id any, update_value models.ToDoList_DB) (*mongo.UpdateResult, error) {
	collection := db.Collection("todos")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": update_value}

	// UpdateOne executes an update command to update at most one document in the collection.
	// The update parameter must be a document containing update operators
	return collection.UpdateOne(ctx, filter, update)
}

func DeleteTodo(ctx context.Context, db *mongo.Database, id any) (*mongo.DeleteResult, error) {
	collection := db.Collection("todos")
	filter := bson.M{"_id": id}
	return collection.DeleteOne(ctx, filter)
}

func DeleteTodoByFilter(ctx context.Context, db *mongo.Database, filter bson.M) (*mongo.DeleteResult, error) {
	collection := db.Collection("todos")
	return collection.DeleteMany(ctx, filter)
}


// User_DB specific
func InsertUser(ctx context.Context, db *mongo.Database, user models.User_DB) (*mongo.InsertOneResult, error) {
	collection := db.Collection("users")
	return collection.InsertOne(ctx, user)
}

func FindUserByUsername(ctx context.Context, db *mongo.Database, username string) (models.User_DB, error) {
	var user models.User_DB
	err := db.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return user, err
}