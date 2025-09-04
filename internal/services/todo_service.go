package services

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"todo-api/internal/models"
	"todo-api/internal/repository"
)

type TodoService interface {
	CreateTodo(request *models.CreateTodoRequest) (*models.Todo, error)
	GetById(id string) (*models.Todo, error)
	GetAllTodos() ([]*models.Todo, error)
	UpdateTodo(id string, request *models.UpdateTodoRequest) (*models.Todo, error)
	DeleteTodo(id string) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) CreateTodo(request *models.CreateTodoRequest) (*models.Todo, error) {
	if request.TaskName == "" {
		return nil, errors.New("необходимо передать наименование задачи")
	}

	name := strings.TrimSpace(request.TaskName)

	task := models.Todo{
		ID:          uuid.NewString(),
		TaskName:    name,
		Description: request.Description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}

	err := s.repo.Create(&task)

	return &task, err
}

func (s *todoService) GetById(id string) (*models.Todo, error) {
	task, err := s.repo.GetById(id)
	return task, err
}

func (s *todoService) GetAllTodos() ([]*models.Todo, error) {
	return s.repo.GetAllTask()
}

func (s *todoService) UpdateTodo(id string, request *models.UpdateTodoRequest) (*models.Todo, error) {
	err := s.repo.Update(id, request)
	if err != nil {
		return nil, err
	}

	return s.repo.GetById(id)
}

func (s *todoService) DeleteTodo(id string) error {
	return s.repo.Delete(id)
}
