package api

/*
this file will implement the logic
when a route is called (control layer)

imports required:
this file will import the service,
utils, gin,

*/

import (
	"context"
	"github.com/ayush-bothra/backend-optimizer/internal/db"
	"github.com/ayush-bothra/backend-optimizer/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
	"time"
)

func (h *Handler) GetTodobyIDHandler(c *gin.Context) {
	// Param returns the value of the URL param. It is a shortcut for c.Params.ByName(key)
	id := c.Param("id")

	// ObjectIDFromHex creates a new ObjectID from a hex string.
	// It returns an error if the hex string is not a valid ObjectID.
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	todo, err := db.GetTodoByID(ctx, h.db.Database(), objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error: ": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, todo)

}

// suggestion: add pagination
func (h *Handler) GetTodos(c *gin.Context) {
	// WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)).
	//
	// Canceling this context releases resources associated with it, so code should
	// call cancel as soon as the operations running in this [Context] complete:
	// this creates a context ctx, which will be used everywhere
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func (h *Handler) CreateTodoByID(c *gin.Context) {
	var new_todos models.ToDoList_DB

	// ShouldBindJSON is a shortcut for c.ShouldBindWith(obj, binding.JSON).
	// ShouldBindWith binds the passed struct pointer using the specified binding engine.
	// Binding describes the interface which needs to be implemented for binding the
	// data present in the request such as JSON request body, query parameters or
	// the form POST.
	if err := c.ShouldBindJSON(&new_todos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	res, err := db.InsertTodo(ctx, h.db.Database(), new_todos)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"insertedID ": res.InsertedID})
}

func (h *Handler) UpdateTodoByID(c *gin.Context) {
	var updated_todos models.ToDoList_DB
	id := c.Param("id")

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&updated_todos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	res, err := db.UpdateTodo(ctx, h.db.Database(), objID, updated_todos)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matchedCount ":  res.MatchedCount,
		"modifiedCount ": res.ModifiedCount,
	})
}

func (h *Handler) DeleteTodoByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	res, err := db.DeleteTodo(ctx, h.db.Database(), objID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error ": "no matching documents found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deleteCount ":  res.DeletedCount,
		"Acknowledged ": res.Acknowledged,
	})
}

func (h *Handler) DeleteTodoByType(c *gin.Context) {
	var filter bson.M

	// Bind JSON into BSON.M
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	res, err := db.DeleteTodoByFilter(ctx, h.db.Database(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error ": "no matching documents found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deleteCount ":  res.DeletedCount,
		"Acknowledged ": res.Acknowledged,
	})
}
