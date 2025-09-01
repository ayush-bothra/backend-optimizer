package main

/*
This is the entry point of the application
it will initialize the server, load the
configs, set up routes, and start the Fibre
app (required for HTTP requests (?))

imports needed:
will be importing the api, utils, cache, db
from the current structure, and will also
need gin/v2 for the HTTP

functions required:
main
*/

import (
	"context"
	"github.com/ayush-bothra/backend-optimizer/internal/api"
	"github.com/ayush-bothra/backend-optimizer/internal/db"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	DB, client := db.ConnectDB(mongoURI)

	// This function will run after connect
	// is done running, similar to the
	// destructor (other sections of the code will function as normal)
	// this one here will disconnect from the db
	// and if there is a problem, it will PANIK
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	r := api.SetUpRoutes(DB.Collection("Todo_DB"))

	r.Run(":8080")
}
