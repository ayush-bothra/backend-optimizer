package models

import "time"

type CachedDB struct {
	ID any `bson:"_id, omitempty"`
	UserID string `bson:"userid"`
	Message string `bson:"message"`
	CreatedAt time.Time `bson:"created_at"`
	Seq int `bson:"seq"`
}