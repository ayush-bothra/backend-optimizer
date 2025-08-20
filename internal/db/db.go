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
FindManyDocuments(ctx context.Context, collection string, filter interface{}, results interface{}) error
UpdateDocument(ctx context.Context, collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
DeleteDocument(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error)
(for CRUD operations)
another file called models.go may be required
*/