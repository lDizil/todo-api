package integration_tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"todo-api/internal/handlers"
	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setUpRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	todoRepo := repository.NewPostgresRepository(db)
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	router.POST("/todos", todoHandler.CreateTodo)
	router.GET("/todos/:id", todoHandler.GetById)
	router.GET("/todos", todoHandler.GetAllTask)
	router.PATCH("/todos/:id", todoHandler.Update)
	router.DELETE("/todos/:id", todoHandler.Delete)

	return router
}

func TestCreateTodo_Integration(t *testing.T) {
	db := SetUpTest(t)
	router := setUpRouter(db)

	reqBody := `{"taskName":"integration test","description":"test description"}`
	req := httptest.NewRequest("POST", "/todos", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var response models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.CreatedAt)
	assert.False(t, response.Completed)
	assert.Equal(t, "integration test", response.TaskName)
	assert.Equal(t, "test description", *response.Description)

	var todoInDB models.Todo
	query := "SELECT id, task_name, description, completed, created_at FROM todos WHERE id=$1"
	err = db.QueryRow(query, response.ID).Scan(
		&todoInDB.ID,
		&todoInDB.TaskName,
		&todoInDB.Description,
		&todoInDB.Completed,
		&todoInDB.CreatedAt,
	)

	assert.NoError(t, err)
	assert.Equal(t, response.ID, todoInDB.ID)
	assert.Equal(t, "integration test", todoInDB.TaskName)
	assert.NotNil(t, response.Description)
	assert.Equal(t, "test description", *todoInDB.Description)
}

func TestGetByID_Integration(t *testing.T) {
	db := SetUpTest(t)
	router := setUpRouter(db)

	description := "test description"
	todo := CreateTestTodo(db, "test task", &description)

	req := httptest.NewRequest("GET", "/todos/"+todo.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response *models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test task", response.TaskName)
	assert.Equal(t, "test description", *response.Description)
}

func TestGetAllTask_Integration(t *testing.T) {
	db := SetUpTest(t)
	router := setUpRouter(db)

	description1 := "test description1"
	_ = CreateTestTodo(db, "test task1", &description1)

	description2 := "test description2"
	_ = CreateTestTodo(db, "test task2", &description2)

	req := httptest.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response []*models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response[0].ID)
	assert.NotEmpty(t, response[0].CreatedAt)
	assert.Equal(t, "test task1", response[0].TaskName)
	assert.Equal(t, "test description1", *response[0].Description)

	assert.NotEmpty(t, response[1].ID)
	assert.NotEmpty(t, response[1].CreatedAt)
	assert.Equal(t, "test task2", response[1].TaskName)
	assert.Equal(t, "test description2", *response[1].Description)
}

func TestUpdateTodo_Integration(t *testing.T) {
	db := SetUpTest(t)
	router := setUpRouter(db)

	description := "test description"
	todo := CreateTestTodo(db, "test task", &description)

	reqBody := `{"taskName":"integration test new","description":"test description new", "completed":true}`
	req := httptest.NewRequest("PATCH", "/todos/"+todo.ID, bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response *models.Todo

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "integration test new", response.TaskName)
	assert.Equal(t, "test description new", *response.Description)
	assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.TaskName)

	var todoInDB models.Todo
	query := "SELECT id, task_name, description, completed, created_at FROM todos WHERE id=$1"
	err = db.QueryRow(query, response.ID).Scan(
		&todoInDB.ID,
		&todoInDB.TaskName,
		&todoInDB.Description,
		&todoInDB.Completed,
		&todoInDB.CreatedAt,
	)

	assert.NoError(t, err)
	assert.Equal(t, response.ID, todoInDB.ID)
	assert.Equal(t, "integration test new", todoInDB.TaskName)
	assert.NotNil(t, response.Description)
	assert.Equal(t, "test description new", *todoInDB.Description)

}

func TestDeleteTodo_Integration(t *testing.T) {
	db := SetUpTest(t)
	router := setUpRouter(db)

	description := "test description"
	todo := CreateTestTodo(db, "test task", &description)

	req := httptest.NewRequest("DELETE", "/todos/"+todo.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
	assert.Empty(t, w.Body.String())
}
