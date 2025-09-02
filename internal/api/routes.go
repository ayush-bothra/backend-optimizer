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
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func NewAuthMiddleWare() (*validator.Validator, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// returns a url.URL and the error during parsing, if any
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

// HandlerFunc defines the handler used by gin middleware as return value.
// type HandlerFunc func(*Context)
func GinJWTMiddleware(m *jwtmiddleware.JWTMiddleware) gin.HandlerFunc {
	// Middleware must wrap the handler instead of running once,
	// so we return a func that validates JWT before calling the next handler.

	return func(c *gin.Context) {
		// call net/http middleware:

		// note: this is gin.__ and shld be handled with care
		w := c.Writer
		// this is *http.request
		r := c.Request

		// jwtmiddleware works on http.Handler
		handler := m.CheckJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//  continue only if JWT is valid
			if !c.IsAborted() {
				c.Next()
			}
		}))

		handler.ServeHTTP(w, r)

		// If we got here without c.Next(), it means auth failed
		if c.IsAborted() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}
	}
}

func SetUpRoutes(db *mongo.Collection) *gin.Engine {
	// initialize the gin Engine
	r := gin.Default()

	// creating the validator using above func:
	jwtValidator, err := NewAuthMiddleWare()
	if err != nil {
		log.Fatalf("failed to set up the middleware: %v", err)
	}

	// The jwtValidator is a *validator, which contains the ValidateToken func
	// so it just catches that func in here, no need to define separately
	// ValidateToken validates the passed in JWT using the jose v2 package.
	// WithCredentialsOptional sets up if credentials are optional or not.
	// If set to true then an empty token will be considered valid.
	jwtMw := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithCredentialsOptional(false),
	)

	// creates a new Handler
	h := NewHandler(db)
	uh := NewUserHandler(db)
	// note -> gin accepts multiple Handlerfuncs.
	gjwt := GinJWTMiddleware(jwtMw)

	// call the GET API
	// --- r.GET("/get-token", uh.LoginUser) --- {manual login, use oauth} 
	r.GET("/todos/", gjwt, h.GetTodos)
	r.GET("/todos/:id", gjwt, h.GetTodobyIDHandler)
	// call the POST API
	r.POST("/Register", uh.Register)
	r.POST("/Login", uh.Login)
	r.POST("/todos", gjwt, h.CreateTodoByID)
	// call the PUT API
	r.PUT("/todos/:id", gjwt, h.UpdateTodoByID)
	// call the DELETE API
	r.DELETE("/todos/:id", gjwt, h.DeleteTodoByID)
	r.DELETE("/todos/", gjwt, h.DeleteTodoByType)

	return r
}
