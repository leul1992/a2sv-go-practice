package Usecases

import (
	"errors"
	"task-manager/Domain"
	"task-manager/Infrastructure"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	GetAll() ([]Domain.Task, error)
	GetByID(id int) (Domain.Task, error)
	Create(t Domain.Task) (Domain.Task, error)
	Update(id int, updated Domain.Task) (Domain.Task, error)
	Delete(id int) error
}

type TaskUseCase interface {
	GetAllTasks() ([]Domain.Task, error)
	GetTaskByID(id int) (Domain.Task, error)
	CreateTask(t Domain.Task) (Domain.Task, error)
	UpdateTask(id int, updated Domain.Task) (Domain.Task, error)
	DeleteTask(id int) error
}

type taskUseCase struct {
	repo TaskRepository
}

func NewTaskUseCase(repo TaskRepository) TaskUseCase {
	return &taskUseCase{repo: repo}
}

func (tuc *taskUseCase) GetAllTasks() ([]Domain.Task, error) {
	return tuc.repo.GetAll()
}

func (tuc *taskUseCase) GetTaskByID(id int) (Domain.Task, error) {
	task, err := tuc.repo.GetByID(id)
	if errors.Is(err, Infrastructure.ErrNoDocuments) {
		return Domain.Task{}, ErrTaskNotFound
	}
	if err != nil {
		return Domain.Task{}, err
	}
	return task, nil
}

func (tuc *taskUseCase) CreateTask(t Domain.Task) (Domain.Task, error) {
	return tuc.repo.Create(t)
}

func (tuc *taskUseCase) UpdateTask(id int, updated Domain.Task) (Domain.Task, error) {
	task, err := tuc.repo.Update(id, updated)
	if err != nil && errors.Is(err, errors.New("task not found")) {
		return Domain.Task{}, ErrTaskNotFound
	}
	return task, err
}

func (tuc *taskUseCase) DeleteTask(id int) error {
	err := tuc.repo.Delete(id)
	if err != nil && errors.Is(err, errors.New("task not found")) {
		return ErrTaskNotFound
	}
	return err
}
