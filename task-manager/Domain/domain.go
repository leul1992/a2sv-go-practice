package Domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          int    `json:"id" bson:"id"`
	Title       string `json:"title" binding:"required" bson:"title"`
	Description string `json:"description" bson:"description"`
	DueDate     string `json:"due_date" bson:"due_date"`
	Status      string `json:"status" bson:"status"`
}

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserName string             `json:"username" bson:"username"`
	Password string             `json:"-" bson:"password"`
	Role     string             `json:"role" bson:"role"`
}
