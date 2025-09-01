package models

// for typed IDs, use bson.ObjectID
// here any is fine, so we will use it
// any is similar to auto_inc from mysql
type ToDoList_DB struct {
	ID    any    `bson:"_id,omitempty"`
	Title string `bson:"title"`
	Done  bool   `bson:"done"`
}