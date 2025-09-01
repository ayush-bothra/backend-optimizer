package api

import (
	"context"
	"net/http"
	"time"
	"github.com/ayush-bothra/backend-optimizer/internal/db"
	"github.com/ayush-bothra/backend-optimizer/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db *mongo.Collection
}

func NewUserHandler(col *mongo.Collection) *UserHandler {
	// struct construction, create new UserHandler
	return &UserHandler{db: col}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"requried"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return 
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	// Create a new user variable
	user := models.User_DB{
		Username: request.Username,
		Password: string(hashed),
	}

	res, err := db.InsertUser(ctx, h.db.Database(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"id": res.InsertedID})
}