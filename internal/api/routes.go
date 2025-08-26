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
	r.GET("/todos/", h.GetTodos)
	r.GET("/todos/:id", h.GetTodobyIDHandler)
	// call the POST API
	r.POST("/todos", h.CreateTodoByID)
	// call the PUT API
	r.PUT("/todos/:id", h.UpdateTodoByID)
	// call the DELETE API
	r.DELETE("/todos/:id", h.DeleteTodoByID)
	r.DELETE("/todos/", h.DeleteTodoByType)
	return r
}