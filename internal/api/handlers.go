package api

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Handler struct {
	db *mongo.Collection
	rdb *redis.Client
}

func NewHandler(col *mongo.Collection, rdb* redis.Client) *Handler {
	// struct construction, create new models.Handler
	return &Handler{
		db: col,
		rdb: rdb,
	}
}