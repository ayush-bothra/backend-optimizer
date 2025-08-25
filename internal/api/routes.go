package api

/*
this file will define all the HTTP
routes and endpoints required

imports needed:
gin (HTTP service)
handlers (for route handling)

functions needed:
setUpRoutes(app *fiber.App)
*/

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetUpRoutes(db *mongo.Collection) *gin.Engine {
	// initialize the gin Engine
	r := gin.Default()

	// creates a new Handler
	h := NewHandler(db)

	// call the GET API
	r.GET("/todos", h.GetTodos)
	return r
}