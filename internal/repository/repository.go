package repository

import (
	"errors"
	"strings"
	"todo-api/internal/models"
)

type TodoRepository interface {
	Create(task *models.Todo) error
	GetById(id string) (*models.Todo, error)
	Update(id string, updateData *models.UpdateTodoRequest) error
	Delete(id string) error
	GetAllTask() ([]*models.Todo, error)
}

type StorageRepository struct {
	todos map[string]*models.Todo
}

func Constructor() *StorageRepository {
	return &StorageRepository{
		todos: make(map[string]*models.Todo),
	}
}

var ErrEmptyID = errors.New("передан пустой айди")
var ErrInvalidID = errors.New("задача с таким айди не найдена")
var ErrEmptyTask = errors.New("передана пустая задача")
var ErrEmptyData = errors.New("переданы пустые данные")
var ErrAlreadyExist = errors.New("задача с таким айди уже существует")
var ErrEmptyName = errors.New("необходимо передать наименование задачи")

func (s *StorageRepository) Create(task *models.Todo) error {
	if task == nil {
		return ErrEmptyTask
	}

	if task.ID == "" {
		return ErrEmptyID
	}

	if _, exists := s.todos[task.ID]; exists {
		return ErrAlreadyExist
	}

	s.todos[task.ID] = task

	return nil
}

func (s *StorageRepository) GetById(id string) (*models.Todo, error) {
	if id == "" {
		return nil, ErrEmptyID
	}

	if task, exists := s.todos[id]; exists {
		return task, nil
	}

	return nil, ErrInvalidID
}

func (s *StorageRepository) Update(id string, updateData *models.UpdateTodoRequest) error {
	if id == "" {
		return ErrEmptyID
	}

	if updateData == nil {
		return ErrEmptyTask
	}

	if task, exists := s.todos[id]; exists {
		if updateData.Completed != nil {
			task.Completed = *updateData.Completed
		}
		if updateData.Description != nil {
			task.Description = updateData.Description
		}
		if updateData.TaskName != nil {
			name := strings.TrimSpace(*updateData.TaskName)
			if name == "" {
				return ErrEmptyName
			}
			task.TaskName = name
		}
	} else {
		return ErrInvalidID
	}

	return nil
}

func (s *StorageRepository) Delete(id string) error {
	if id == "" {
		return ErrEmptyID
	}

	if _, exists := s.todos[id]; exists {
		delete(s.todos, id)
		return nil
	}

	return ErrInvalidID
}

func (s *StorageRepository) GetAllTask() ([]*models.Todo, error) {
	var result []*models.Todo

	for _, task := range s.todos {
		result = append(result, task)
	}

	return result, nil
}
