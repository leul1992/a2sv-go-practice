package Infrastructure

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client          *mongo.Client
	TasksCollection *mongo.Collection
	UsersCollection *mongo.Collection
	ErrNoDocuments  = mongo.ErrNoDocuments
)

func InitDB() error {
	godotenv.Load()

	uri := os.Getenv("MONGODB_URI")
	log.Println("Connecting to MongoDB at", uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	Client = client
	TasksCollection = client.Database("taskdb").Collection("tasks")
	UsersCollection = client.Database("taskdb").Collection("users")

	log.Println("MongoDB Connection Established")

	return nil
}

func GetTasksCollection() *mongo.Collection {
	return TasksCollection
}

func GetUsersCollection() *mongo.Collection {
	return UsersCollection
}
