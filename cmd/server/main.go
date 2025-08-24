package main


/*
This is the entry point of the application
it will initialize the server, load the 
configs, set up routes, and start the Fibre
app (required for HTTP requests (?))

imports needed:
will be importing the api, utils, cache, db
from the current structure, and will also
need fiber/v2 for the HTTP

functions required: 
main
*/

import (
	"fmt"
	"net/http"
	"github.com/ayush-bothra/backend-optimizer/internal/db"
	"github.com/joho/godotenv"
	"log"
	"os"
	"context"

)


func helloHandler(wrt http.ResponseWriter, req *http.Request) {
	fmt.Fprint(wrt, "Hello World!\n")
}

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

	todo := db.ToDoList_DB{Title: "Learn mongoDB", Done: false}
	res, err := db.InsertTodo(DB, todo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ID: ", res.InsertedID)

	// HandleFunc registers the handler function for the given pattern in [DefaultServeMux].
	// The documentation for [ServeMux] explains how patterns are matched.
	http.HandleFunc("/", helloHandler)

	fmt.Println("server running at http://localhost:8080/")

	// ListenAndServe listens on the TCP network address addr and then calls
	// [Serve] with handler to handle requests on incoming connections.
	// The handler is typically nil, in which case [DefaultServeMux] is used.
	http.ListenAndServe(":8080", nil)
}