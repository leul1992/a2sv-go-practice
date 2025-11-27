package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserName string             `json:"username" bson:"username"`
	Password string             `json:"-" bson:"password"`
	Role     string             `json:"role" bson:"role"`
}
