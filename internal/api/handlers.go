package api


/*
this file will implement the logic
when a route is called (control layer)

imports required:
this file will import the service,
utils, gin,

functions needed:
RiskCheckHandler(c *fiber.Ctx)
PortfolioHandler(c *fiber.Ctx) (optional)
*/

import (
	"net/http"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type Handler struct {
	db *mongo.Collection
}

func NewHandler(col *mongo.Collection) *Handler {
	// struct construction, create new Handler
	return &Handler{db: col}
}

func (h *Handler) GetTodos(c *gin.Context) (){

	// WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)).
	//
	// Canceling this context releases resources associated with it, so code should
	// call cancel as soon as the operations running in this [Context] complete:
	// this creates a context ctx, which will be used everywhere
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	// Find executes a find command and returns a Cursor over the matching documents in the collection.
	// kind of similar to map.Find()
	cursor, err := h.db.Find(ctx, bson.M{})

	if err != nil {
		// return the error as JSON
		// H is a shortcut for map[string]any
		c.JSON(http.StatusInternalServerError, gin.H{"error: ": err.Error()})
		return
	}
	// cleanup after cursor has completed run
	defer cursor.Close(ctx)


	var todos []bson.M 
	if err := cursor.All(ctx, &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error: ": err.Error()})
		return 
	}

	c.IndentedJSON(http.StatusOK, todos)
}