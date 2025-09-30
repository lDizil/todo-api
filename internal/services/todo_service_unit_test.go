package services

import (
	"testing"
	"todo-api/internal/models"
	"todo-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	createErr  error
	getByIdErr error
	updateErr  error
	deleteErr  error
}

func (m *mockRepo) Create(task *models.Todo) error {
	return m.createErr
}

func (m *mockRepo) GetById(id string) (*models.Todo, error) {
	return nil, m.getByIdErr

}
func (m *mockRepo) GetAllTask() ([]*models.Todo, error) {
	return nil, nil
}

func (m *mockRepo) Update(id string, req *models.UpdateTodoRequest) error {
	return m.updateErr
}

func (m *mockRepo) Delete(id string) error {
	return m.deleteErr
}

func TestTodoService_CreateTodo(t *testing.T) {
	repo := repository.Constructor()
	services := NewTodoService(repo)
	description := "test text"

	req := &models.CreateTodoRequest{TaskName: "test", Description: &description}

	todo, err := services.CreateTodo(req)

	assert.NoError(t, err)
	assert.Equal(t, "test", todo.TaskName)
	assert.Equal(t, "test text", *todo.Description)
}

func TestTodoService_CreateTodo_ErrEmptyName(t *testing.T) {
	repo := repository.Constructor()
	services := NewTodoService(repo)
	description := "test text"

	req := &models.CreateTodoRequest{TaskName: "", Description: &description}

	_, err := services.CreateTodo(req)

	assert.ErrorIs(t, err, repository.ErrEmptyName)
}

func TestTodoService_CreateTodo_ErrRepo(t *testing.T) {
	req := &models.CreateTodoRequest{TaskName: "test"}

	repo := &mockRepo{createErr: repository.ErrAlreadyExist}
	services := NewTodoService(repo)
	_, err := services.CreateTodo(req)
	assert.ErrorIs(t, err, repository.ErrAlreadyExist)

	repo = &mockRepo{createErr: repository.ErrEmptyTask}
	services = NewTodoService(repo)
	_, err = services.CreateTodo(req)
	assert.ErrorIs(t, err, repository.ErrEmptyTask)
}

func TestTodoService_GetById(t *testing.T) {
	repo := repository.Constructor()
	services := NewTodoService(repo)

	description := "test text"
	req := &models.CreateTodoRequest{TaskName: "test", Description: &description}
	todo, _ := services.CreateTodo(req)

	todo, err := services.GetById(todo.ID)

	assert.NoError(t, err)
	assert.Equal(t, "test", todo.TaskName)
	assert.Equal(t, "test text", *todo.Description)
}

func TestTodoService_GetById_ErrRepo(t *testing.T) {
	repo := &mockRepo{getByIdErr: repository.ErrEmptyID}
	services := NewTodoService(repo)

	id := ""
	_, err := services.GetById(id)
	assert.Error(t, err, repository.ErrEmptyID)

	repo = &mockRepo{getByIdErr: repository.ErrInvalidID}
	services = NewTodoService(repo)

	id = "test"
	_, err = services.GetById(id)
	assert.Error(t, err, repository.ErrInvalidID)
}

func TestTodoService_GetAllTodos(t *testing.T) {
	repo := repository.Constructor()
	services := NewTodoService(repo)

	description1 := "test text 1"
	req := &models.CreateTodoRequest{TaskName: "test1", Description: &description1}
	_, _ = services.CreateTodo(req)

	description2 := "test text 2"
	req = &models.CreateTodoRequest{TaskName: "test2", Description: &description2}
	_, _ = services.CreateTodo(req)

	todos, err := services.GetAllTodos()

	assert.NoError(t, err)
	assert.Len(t, todos, 2)
	assert.Equal(t, "test1", todos[0].TaskName)
	assert.Equal(t, "test2", todos[1].TaskName)
	assert.Equal(t, "test text 1", *todos[0].Description)
	assert.Equal(t, "test text 2", *todos[1].Description)
}

func TestTodoService_Update(t *testing.T) {
	repo := repository.Constructor()
	services := NewTodoService(repo)

	description := "test text 1"
	req := &models.CreateTodoRequest{TaskName: "test1", Description: &description}
	todo, _ := services.CreateTodo(req)

	descriptionNew := "test text new"
	taskNameNew := "testNew"
	updReq := &models.UpdateTodoRequest{TaskName: &taskNameNew, Description: &descriptionNew}
	newTodo, err := services.UpdateTodo(todo.ID, updReq)

	assert.NoError(t, err)
	assert.Equal(t, "testNew", newTodo.TaskName)
	assert.Equal(t, "test text new", *newTodo.Description)
}

func TestTodoService_Update_ErrRepo(t *testing.T) {
	repo := &mockRepo{updateErr: repository.ErrEmptyTask}
	services := NewTodoService(repo)

	_, err := services.UpdateTodo("1", nil)

	assert.ErrorIs(t, err, repository.ErrEmptyTask)

	repo = &mockRepo{updateErr: repository.ErrEmptyID}
	services = NewTodoService(repo)

	_, err = services.UpdateTodo("", nil)

	assert.ErrorIs(t, err, repository.ErrEmptyID)

	repo = &mockRepo{updateErr: repository.ErrEmptyName}
	services = NewTodoService(repo)

	description := "test text 1"
	req := &models.CreateTodoRequest{TaskName: "test1", Description: &description}
	todo, _ := services.CreateTodo(req)

	descriptionNew := "test text new"
	taskNameNew := ""
	updReq := &models.UpdateTodoRequest{TaskName: &taskNameNew, Description: &descriptionNew}
	_, err = services.UpdateTodo(todo.ID, updReq)

	assert.ErrorIs(t, err, repository.ErrEmptyName)

	repo = &mockRepo{updateErr: repository.ErrInvalidID}
	services = NewTodoService(repo)

	updReq = &models.UpdateTodoRequest{TaskName: &taskNameNew, Description: &descriptionNew}
	_, err = services.UpdateTodo("231", updReq)

	assert.ErrorIs(t, err, repository.ErrInvalidID)
}

func TestTodoService_Delete(t *testing.T) {
	repo := repository.Constructor()
	services := NewTodoService(repo)
	description := "test text"

	req := &models.CreateTodoRequest{TaskName: "test", Description: &description}

	todo, _ := services.CreateTodo(req)

	err := services.DeleteTodo(todo.ID)
	todos, _ := services.GetAllTodos()

	assert.NoError(t, err)
	assert.Len(t, todos, 0)
}

func TestTodoService_Delete_ErrRepo(t *testing.T) {
	repo := &mockRepo{deleteErr: repository.ErrEmptyID}
	services := NewTodoService(repo)

	err := services.DeleteTodo("")
	assert.ErrorIs(t, err, repository.ErrEmptyID)

	repo = &mockRepo{deleteErr: repository.ErrInvalidID}
	services = NewTodoService(repo)

	err = services.DeleteTodo("213")
	assert.ErrorIs(t, err, repository.ErrInvalidID)
}
