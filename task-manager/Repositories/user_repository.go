package Repositories

import (
	"context"
	"errors"
	"task-manager/Domain"
	"task-manager/Infrastructure"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindByUsername(username string) (*Domain.User, error)
	Count() (int64, error)
	Create(user Domain.User) error
	FindByID(id string) (*Domain.User, error)
	UpdateRole(id primitive.ObjectID, role string) error
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository() UserRepository {
	return &userRepository{collection: Infrastructure.GetUsersCollection()}
}

func (ur *userRepository) FindByUsername(username string) (*Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user Domain.User
	err := ur.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return ur.collection.CountDocuments(ctx, bson.M{})
}

func (ur *userRepository) Create(user Domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := ur.collection.InsertOne(ctx, user)
	return err
}

func (ur *userRepository) FindByID(id string) (*Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	var user Domain.User
	err = ur.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) UpdateRole(id primitive.ObjectID, role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"role": role}}

	result, err := ur.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
