package data

import (
	"errors"
	"task_manager/models"
)

var (
	tasks  = []models.Task{}
	nextID = 1
)

func GetAll() []models.Task {
	return tasks
}

func GetByID(id int) (models.Task, error) {
	for _, t := range tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return models.Task{}, errors.New("task not found")
}

func Create(t models.Task) models.Task {
	t.ID = nextID
	nextID++
	tasks = append(tasks, t)
	return t
}

func Update(id int, updated models.Task) (models.Task, error) {
	for i, t := range tasks {
		if t.ID == id {
			updated.ID = id
			tasks[i] = updated
			return updated, nil
		}
	}
	return models.Task{}, errors.New("task not found")
}

func Delete(id int) error {
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
