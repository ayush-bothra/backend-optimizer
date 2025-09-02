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


// This code is for manual sign-in and log-in
// no longer required since oauth will be used
func (h *UserHandler) RegisterUser(c *gin.Context) {
	// basically created a set of variables that will help
	// in getting the values from the web server
	// we will operate on these values, just line cin , cout
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"requried"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	// GenerateFromPassword returns the bcrypt hash of the password at the given
	// cost. If the cost given is less than MinCost, the cost will be set to
	// DefaultCost instead. 
	// Use CompareHashAndPassword, as defined in this package,
	// to compare the returned hashed password with its cleartext version.
	// GenerateFromPassword does not accept passwords longer than 72 bytes, which
	// is the longest password bcrypt will operate on.
	// byte is alias for uint8
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

func (h *UserHandler) LoginUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	} 

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	s_user, err := db.FindUserByUsername(ctx, h.db.Database(), request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	if err := bcrypt.CompareHashAndPassword([]byte(s_user.Password), 
	[]byte(request.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := RegisterUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	c.JSON(http.StatusOK, resp)
}