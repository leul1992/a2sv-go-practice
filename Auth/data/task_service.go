package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"errors"
	"task_manager/models"
)

func GetAll() []models.Task {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := TasksCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error fetching tasks:", err)
		return []models.Task{}
	}
	var tasks []models.Task
	if err = cur.All(ctx, &tasks); err != nil {
		log.Println("Error decoding tasks:", err)
		return []models.Task{}
	}
	return tasks
}

func GetByID(id int) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var task models.Task

	err := TasksCollection.FindOne(ctx, bson.M{"id": id}).Decode(&task)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func Create(t models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "maxID", Value: bson.D{{Key: "$max", Value: "$id"}}},
		}}},
	}
	cur, err := TasksCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return models.Task{}, err
	}
	var results []bson.M
	if err = cur.All(ctx, &results); err != nil {
		return models.Task{}, err
	}
	var nextID int
	if len(results) == 0 {
		nextID = 1
	} else {
		maxID := results[0]["maxID"]
		switch v := maxID.(type) {
		case int32:
			nextID = int(v) + 1
		case int64:
			nextID = int(v) + 1
		default:
			nextID = 1
		}
	}
	t.ID = nextID
	_, err = TasksCollection.InsertOne(ctx, t)
	if err != nil {
		return models.Task{}, err
	}
	return t, nil
}

func Update(id int, updated models.Task) (models.Task, error) {
	updated.ID = id
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.Update().SetUpsert(false)
	res, err := TasksCollection.UpdateOne(ctx, bson.D{{Key: "id", Value: id}}, bson.D{{Key: "$set", Value: updated}}, opts)
	if err != nil {
		return models.Task{}, err
	}
	if res.MatchedCount == 0 {
		return models.Task{}, errors.New("task not found")
	}
	return GetByID(id)
}

func Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := TasksCollection.DeleteOne(ctx, bson.D{{Key: "id", Value: id}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
