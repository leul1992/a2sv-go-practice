package data

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"task_manager/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
    Client          *mongo.Client
    TasksCollection *mongo.Collection
    UsersCollection *mongo.Collection
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


func RegisterUser(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	username = strings.TrimSpace(username)
	
	var existing models.Users

	err := UsersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&existing)
	if err == nil {
		return errors.New("username already taken")
	}
	if err != mongo.ErrNoDocuments {
		return err
	}

	// Check if this is the first user → make admin
	count, err := UsersCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	role := "user"
	if count == 0 {
		role = "admin"
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.Users{
		ID:       primitive.NewObjectID(),
		UserName: username,
		Password: string(hashed),
		Role:     role,
	}

	_, err = UsersCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func LoginUser(username, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	username = strings.TrimSpace(username)

	var user models.Users
	err := UsersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT secret not set in environment")
	}
	
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID.Hex(),
		"username": user.UserName,
		"role": user.Role,
		"exp":  time.Now().Add(72 * time.Hour).Unix(),
	})

	jwtTokenString, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return jwtTokenString, nil
}

func PromoteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user id")
	}

	var user models.Users
	err = UsersCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
		return err
	}

	if user.Role == "admin" {
		return errors.New("user is already an admin")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"role": "admin"}}

	result, err := UsersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
