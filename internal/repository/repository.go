package repository

import (
	"errors"
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

var errEmptyID = errors.New("передан пустой айди")

func (s *StorageRepository) Create(task *models.Todo) error {
	if task == nil {
		return errors.New("передана пустая задача")
	}

	if task.ID == "" {
		return errEmptyID
	}

	if _, exists := s.todos[task.ID]; exists {
		return errors.New("задача с таким айди уже существует")
	}

	s.todos[task.ID] = task

	return nil
}

func (s *StorageRepository) GetById(id string) (*models.Todo, error) {
	if id == "" {
		return nil, errEmptyID
	}

	if task, exists := s.todos[id]; exists {
		return task, nil
	}

	return nil, errors.New("задача с таким айди не найдена")
}

func (s *StorageRepository) Update(id string, updateData *models.UpdateTodoRequest) error {
	if id == "" {
		return errEmptyID
	}

	if updateData == nil {
		return errors.New("необходимо передать данные для обновления")
	}

	if task, exists := s.todos[id]; exists {
		if updateData.Completed != nil {
			task.Completed = *updateData.Completed
		}
		if updateData.Description != nil {
			task.Description = updateData.Description
		}
		if updateData.TaskName != nil {
			task.TaskName = *updateData.TaskName
		}
	} else {
		return errors.New("задача с таким айди не найдена")
	}

	return nil
}

func (s *StorageRepository) Delete(id string) error {
	if id == "" {
		return errEmptyID
	}

	if _, exists := s.todos[id]; exists {
		delete(s.todos, id)
		return nil
	}

	return errors.New("задача для удаления с таким айди не найдена")
}

func (s *StorageRepository) GetAllTask() ([]*models.Todo, error) {
	var result []*models.Todo

	for _, task := range s.todos {
		result = append(result, task)
	}

	return result, nil
}
