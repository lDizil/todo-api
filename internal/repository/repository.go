package repository

import (
	"todo-api/internal/models"
	"errors"
)

type TodoRepository interface {
	Create(task *models.Todo) error
	GetById(id string) (*models.Todo, error)
	Update(task *models.Todo) error
	Delete(id string) error
	GetAllTask() ([]*models.Todo, error)
}

type StorageRepository struct {
	todos map[string] *models.Todo
}

func Constructor() *StorageRepository {
	return $StorageRepository{
		todos: make(map[string] *models.Todo),
	}
}