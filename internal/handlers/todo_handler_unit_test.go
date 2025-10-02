package handlers

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"todo-api/internal/models"
	"todo-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockService struct {
	createTodoFunc  func(req *models.CreateTodoRequest) (*models.Todo, error)
	getByIdFunc     func(id string) (*models.Todo, error)
	getAllTodosFunc func() ([]*models.Todo, error)
	updateTodoFunc  func(id string, req *models.UpdateTodoRequest) (*models.Todo, error)
	deleteTodoFunc  func(id string) error
}

func (m *MockService) CreateTodo(req *models.CreateTodoRequest) (*models.Todo, error) {
	return m.createTodoFunc(req)
}

func (m *MockService) GetById(id string) (*models.Todo, error) {
	return m.getByIdFunc(id)
}

func (m *MockService) GetAllTodos() ([]*models.Todo, error) {
	return m.getAllTodosFunc()
}

func (m *MockService) UpdateTodo(id string, req *models.UpdateTodoRequest) (*models.Todo, error) {
	return m.updateTodoFunc(id, req)
}

func (m *MockService) DeleteTodo(id string) error {
	return m.deleteTodoFunc(id)
}

func TestTodoHadler_Create(t *testing.T) {
	mock := &MockService{
		createTodoFunc: func(req *models.CreateTodoRequest) (*models.Todo, error) {
			return &models.Todo{ID: "1", TaskName: req.TaskName, Description: req.Description}, nil
		},
	}

	handler := NewTodoHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	reqBody := `{"taskName":"test","description":"test text"}`
	c.Request = httptest.NewRequest("POST", "/todos", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateTodo(c)

	assert.Equal(t, 201, w.Code)

	var response *models.Todo
	description := "test text"
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1", response.ID)
	assert.Equal(t, "test", response.TaskName)
	assert.Equal(t, &description, response.Description)
}

func TestTodoHadler_Create_ErrEmptyName(t *testing.T) {
	mock := &MockService{
		createTodoFunc: func(req *models.CreateTodoRequest) (*models.Todo, error) {
			return nil, repository.ErrEmptyName
		},
	}

	handler := NewTodoHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	reqBody := `{"taskName":""}`
	c.Request = httptest.NewRequest("POST", "/todos", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateTodo(c)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"необходимо передать наименование задачи"}`, w.Body.String())
}

func TestTodoHadler_Create_ErrAlreadyExist(t *testing.T) {
	mock := &MockService{
		createTodoFunc: func(req *models.CreateTodoRequest) (*models.Todo, error) {
			return nil, repository.ErrAlreadyExist
		},
	}

	handler := NewTodoHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	reqBody := `{"taskName":""}`
	c.Request = httptest.NewRequest("POST", "/todos", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateTodo(c)

	assert.Equal(t, 409, w.Code)
	assert.JSONEq(t, `{"error":"задача с таким айди уже существует"}`, w.Body.String())
}

func TestTodoHandler_Create_InvalidJSON(t *testing.T) {
	mock := &MockService{
		createTodoFunc: func(req *models.CreateTodoRequest) (*models.Todo, error) {
			return nil, nil
		},
	}

	handler := NewTodoHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := `{"taskName":"test`
	c.Request = httptest.NewRequest("POST", "/todos", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateTodo(c)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"неверный JSON"}`, w.Body.String())
}

func TestTodoHandler_Create_InternalServerError(t *testing.T) {
	mock := &MockService{
		createTodoFunc: func(req *models.CreateTodoRequest) (*models.Todo, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := NewTodoHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	reqBody := `{"taskName":"test"}`
	c.Request = httptest.NewRequest("POST", "/todos", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateTodo(c)

	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"error":"внутренняя ошибка сервера"}`, w.Body.String())
}

func TestTodoHadler_GetById(t *testing.T) {
	description := "test description"
	expectedTodo := &models.Todo{
		ID:          "1",
		TaskName:    "test task",
		Description: &description,
	}

	mock := &MockService{
		getByIdFunc: func(id string) (*models.Todo, error) {
			assert.Equal(t, "1", id)
			return expectedTodo, nil
		},
	}

	handler := NewTodoHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/todos/1", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.GetById(c)

	var response models.Todo

	assert.Equal(t, 200, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTodo.ID, response.ID)
	assert.Equal(t, expectedTodo.TaskName, response.TaskName)
	assert.Equal(t, expectedTodo.Description, response.Description)
}

func TestTodoHandler_GetById_ErrEmptyID(t *testing.T) {
	mock := &MockService{
		getByIdFunc: func(id string) (*models.Todo, error) {
			return nil, repository.ErrEmptyID
		},
	}

	handler := NewTodoHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/todo/", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

	handler.GetById(c)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"передан пустой айди"}`, w.Body.String())
}

func TestTodoHandler_GetById_ErrInvalidID(t *testing.T) {
	mock := &MockService{
		getByIdFunc: func(id string) (*models.Todo, error) {
			return nil, repository.ErrInvalidID
		},
	}

	handler := NewTodoHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/todo/1", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.GetById(c)

	assert.Equal(t, 404, w.Code)
	assert.JSONEq(t, `{"error":"задача с таким айди не найдена"}`, w.Body.String())
}

func TestTodoHandler_GetById_InternalServerError(t *testing.T) {
	mock := &MockService{
		getByIdFunc: func(id string) (*models.Todo, error) {
			return nil, errors.New("внутренняя ошибка сервера")
		},
	}

	handler := NewTodoHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/todo/1", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.GetById(c)

	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"error":"внутренняя ошибка сервера"}`, w.Body.String())
}

func TestTodoHandler_Update(t *testing.T) {
	description := "test text new"
	expectedTodo := &models.Todo{
		ID:          "1",
		TaskName:    "test",
		Description: &description,
	}

	mock := &MockService{
		updateTodoFunc: func(id string, req *models.UpdateTodoRequest) (*models.Todo, error) {
			assert.Equal(t, "1", id)
			return expectedTodo, nil
		},
	}

	handler := NewTodoHandler(mock)

	reqBody := `{"taskName":"test","description":"test text new"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("PATCH", "/todo/1", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.Update(c)

	assert.Equal(t, 200, w.Code)

	var response models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1", response.ID)
	assert.Equal(t, "test", response.TaskName)
	assert.Equal(t, "test text new", *response.Description)
}

func TestTodoHandler_Update_ErrEmptyData(t *testing.T) {
	mock := &MockService{
		updateTodoFunc: func(id string, req *models.UpdateTodoRequest) (*models.Todo, error) {
			return nil, repository.ErrEmptyData
		},
	}

	handler := NewTodoHandler(mock)

	reqBody := `{}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("PATCH", "/todo/1", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.Update(c)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"переданы пустые данные"}`, w.Body.String())
}

func TestTodoHandler_Update_InvalidJSON(t *testing.T) {
	mock := &MockService{}
	handler := NewTodoHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := `{"taskName":"test`
	c.Request = httptest.NewRequest("PATCH", "/todos/1", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.Update(c)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"неверный JSON"}`, w.Body.String())
}

func TestTodoHandler_Update_ErrEmptyName(t *testing.T) {
	mock := &MockService{
		updateTodoFunc: func(id string, req *models.UpdateTodoRequest) (*models.Todo, error) {
			return nil, repository.ErrEmptyName
		},
	}
	handler := NewTodoHandler(mock)

	reqBody := `{"taskName":""}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PATCH", "/todos/1", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.Update(c)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"необходимо передать наименование задачи"}`, w.Body.String())
}

func TestTodoHandler_Update_ErrInvalidID(t *testing.T) {
	mock := &MockService{
		updateTodoFunc: func(id string, req *models.UpdateTodoRequest) (*models.Todo, error) {
			return nil, repository.ErrInvalidID
		},
	}

	handler := NewTodoHandler(mock)

	reqBody := `{"taskName":"test","description":"test text"}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("PATCH", "/todo/1", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.Update(c)

	assert.Equal(t, 404, w.Code)
	assert.JSONEq(t, `{"error":"задача с таким айди не найдена"}`, w.Body.String())
}

func TestTodoHandler_Update_InternalServerError(t *testing.T) {
	mock := &MockService{
		updateTodoFunc: func(id string, req *models.UpdateTodoRequest) (*models.Todo, error) {
			return nil, errors.New("внутренняя ошибка сервера")
		},
	}

	handler := NewTodoHandler(mock)

	reqBody := `{"taskName":"test","description":""}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("PATCH", "/todo/1", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.Update(c)

	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"error":"внутренняя ошибка сервера"}`, w.Body.String())
}

/*
var ErrEmptyID = errors.New("передан пустой айди")
var ErrInvalidID = errors.New("задача с таким айди не найдена")
var ErrEmptyTask = errors.New("передана пустая задача")
var ErrEmptyData = errors.New("переданы пустые данные")
var ErrAlreadyExist = errors.New("задача с таким айди уже существует")
var ErrEmptyName = errors.New("необходимо передать наименование задачи")
*/
