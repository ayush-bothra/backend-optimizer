package api

/*
this file will define all the HTTP
routes and endpoints required

imports needed:
gin (HTTP service)
handlers (for route handling)

will have JWT service as well
using  go-jwt-middleware/v2/validator and jwks

::suggestions -> appleboy/gin-jwt is gin specific
consider using that if time permits

functions needed:
setUpRoutes(app *gin.Engine)
*/

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	// jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// NewAuthValidator creates a validator for incoming Auth0 JWTs.
func NewAuthValidator() (*validator.Validator, error) {
	issuer, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		return nil, err
	}
	audience := os.Getenv("AUTH0_AUDIENCE")

	// Create JWKS (JSON Web Key Set) provider
	provider := jwks.NewCachingProvider(issuer, 5*time.Minute)

	// Parse the JWT:
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuer.String(),
		[]string{audience},
	)
	if err != nil {
		return nil, err
	}
	return jwtValidator, nil
}

// below func will return a gin handlerfunc and requires
// a jwtvalidator as its argument
func Auth0MiddleWare(jwtValidator *validator.Validator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "missing authorization error."})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "missing token."})
			return
		}

		validated, err := jwtValidator.ValidateToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}

		ctx.Set("jwt", validated)
		ctx.Next()
	}
}

// SetUpRoutes configures Gin routes with JWT middleware.
func SetUpRoutes(db *mongo.Collection) *gin.Engine {
	r := gin.Default()

	jwtValidator, err := NewAuthValidator()
	if err != nil {
		log.Fatalf("[JWT] Failed to set up validator: %v", err)
	}

	// jwtMw := jwtmiddleware.New(
	// 	jwtValidator.ValidateToken,
	// 	jwtmiddleware.WithCredentialsOptional(false),
	// )

	gjwt := Auth0MiddleWare(jwtValidator)

	// Handlers
	h := NewHandler(db)
	uh := NewUserHandler(db)

	// Routes
	r.GET("/todos/", gjwt, h.GetTodos)
	r.GET("/todos/:id", gjwt, h.GetTodobyIDHandler)
	r.POST("/Register", uh.Register)
	r.POST("/Login", uh.Login)
	r.POST("/todos", gjwt, h.CreateTodoByID)
	r.PUT("/todos/:id", gjwt, h.UpdateTodoByID)
	r.DELETE("/todos/:id", gjwt, h.DeleteTodoByID)
	r.DELETE("/todos/", gjwt, h.DeleteTodoByType)

	return r
}
