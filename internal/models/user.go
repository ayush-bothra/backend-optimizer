package models


// omitempty tells go to ignore or 
// remove it when encoding to BSON/JSON
type User_DB struct {
	ID any `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username"`
	Email string `bson:"email" json:"email"`
	Password string `bson:"password,omitempty" json:"password"`
}