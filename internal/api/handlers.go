package api

import ("go.mongodb.org/mongo-driver/v2/mongo")

type Handler struct {
	db *mongo.Collection
}

func NewHandler(col *mongo.Collection) *Handler {
	// struct construction, create new models.Handler
	return &Handler{db: col}
}