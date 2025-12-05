package Repositories

import (
	"context"
	"errors"
	"task-manager/Domain"
	"task-manager/Infrastructure"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository interface {
	GetAll() ([]Domain.Task, error)
	GetByID(id int) (Domain.Task, error)
	Create(t Domain.Task) (Domain.Task, error)
	Update(id int, updated Domain.Task) (Domain.Task, error)
	Delete(id int) error
}

type taskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository() TaskRepository {
	return &taskRepository{collection: Infrastructure.GetTasksCollection()}
}

func (tr *taskRepository) GetAll() ([]Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := tr.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var tasks []Domain.Task
	if err = cur.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (tr *taskRepository) GetByID(id int) (Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var task Domain.Task

	err := tr.collection.FindOne(ctx, bson.M{"id": id}).Decode(&task)
	if err != nil {
		return Domain.Task{}, err
	}
	return task, nil
}

func (tr *taskRepository) Create(t Domain.Task) (Domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "maxID", Value: bson.D{{Key: "$max", Value: "$id"}}},
		}}},
	}
	cur, err := tr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return Domain.Task{}, err
	}
	var results []bson.M
	if err = cur.All(ctx, &results); err != nil {
		return Domain.Task{}, err
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
	_, err = tr.collection.InsertOne(ctx, t)
	if err != nil {
		return Domain.Task{}, err
	}
	return t, nil
}

func (tr *taskRepository) Update(id int, updated Domain.Task) (Domain.Task, error) {
	updated.ID = id
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.Update().SetUpsert(false)
	res, err := tr.collection.UpdateOne(ctx, bson.D{{Key: "id", Value: id}}, bson.D{{Key: "$set", Value: updated}}, opts)
	if err != nil {
		return Domain.Task{}, err
	}
	if res.MatchedCount == 0 {
		return Domain.Task{}, errors.New("task not found")
	}
	return tr.GetByID(id)
}

func (tr *taskRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := tr.collection.DeleteOne(ctx, bson.D{{Key: "id", Value: id}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
