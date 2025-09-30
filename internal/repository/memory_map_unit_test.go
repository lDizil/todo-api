package repository

import (
	"testing"
	"todo-api/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestStorageRepo_Create(t *testing.T) {
	repo := Constructor()
	todo := &models.Todo{ID: "1", TaskName: "test"}

	err := repo.Create(todo)
	assert.NoError(t, err)
}

func TestStorageRepo_Create_ErrEmptyTask(t *testing.T) {
	repo := Constructor()
	todo := &models.Todo{ID: "", TaskName: ""}
	err := repo.Create(todo)
	assert.ErrorIs(t, err, ErrEmptyTask)
}

func TestStorageRepo_Create_ErrEmptyID(t *testing.T) {
	repo := Constructor()
	todo := &models.Todo{ID: "", TaskName: "test"}
	err := repo.Create(todo)
	assert.ErrorIs(t, err, ErrEmptyID)
}

func TestStorageRepo_Create_ErrAlreadyExist(t *testing.T) {
	repo := Constructor()
	todo := &models.Todo{ID: "1", TaskName: "test"}
	_ = repo.Create(todo)
	err := repo.Create(todo)
	assert.ErrorIs(t, err, ErrAlreadyExist)
}

func TestStorageRepo_GetByID(t *testing.T) {
	repo := Constructor()
	todo := &models.Todo{ID: "1", TaskName: "test"}

	_ = repo.Create(todo)

	_, err := repo.GetById(todo.ID)

	assert.NoError(t, err)
}

func TestStorageRepo_GetByID_ErrEmptyID(t *testing.T) {
	repo := Constructor()
	_, err := repo.GetById("")
	assert.ErrorIs(t, err, ErrEmptyID)
}

func TestStorageRepo_GetByID_ErrInvalidID(t *testing.T) {
	repo := Constructor()
	_, err := repo.GetById("2")
	assert.ErrorIs(t, err, ErrInvalidID)
}

func TestStorageRepo_Update_Success(t *testing.T) {
	repo := Constructor()
	todo := &models.Todo{ID: "1", TaskName: "test"}
	_ = repo.Create(todo)

	newName := "updated"
	updateData := &models.UpdateTodoRequest{TaskName: &newName}

	err := repo.Update(todo.ID, updateData)
	assert.NoError(t, err)

	updated, err := repo.GetById(todo.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updated", updated.TaskName)
}

func TestStorageRepo_Update_Err_ErrEmptyID(t *testing.T) {
	repo := Constructor()

	taskName := "testNew"
	completed := true

	updateData := &models.UpdateTodoRequest{
		TaskName:  &taskName,
		Completed: &completed,
	}

	err := repo.Update("", updateData)

	assert.ErrorIs(t, err, ErrEmptyID)
}

func TestStorageRepo_Update_Err_ErrEmptyTask(t *testing.T) {
	repo := Constructor()

	todo := &models.Todo{ID: "1", TaskName: "test"}

	err := repo.Update(todo.ID, nil)

	assert.ErrorIs(t, err, ErrEmptyTask)
}

func TestStorageRepo_Update_Err_ErrInvalidID(t *testing.T) {
	repo := Constructor()

	taskName := "testNew"
	completed := true

	updateData := &models.UpdateTodoRequest{
		TaskName:  &taskName,
		Completed: &completed,
	}

	err := repo.Update("1", updateData)

	assert.ErrorIs(t, err, ErrInvalidID)
}

func TestStorageRepo_Delete(t *testing.T) {
	repo := Constructor()
	todo := &models.Todo{ID: "1", TaskName: "test"}
	_ = repo.Create(todo)

	err := repo.Delete(todo.ID)
	assert.NoError(t, err)

	_, err = repo.GetById(todo.ID)
	assert.ErrorIs(t, err, ErrInvalidID)
}

func TestStorageRepo_Delete_ErrEmptyID(t *testing.T) {
	repo := Constructor()

	err := repo.Delete("")

	assert.ErrorIs(t, err, ErrEmptyID)
}

func TestStorageRepo_Delete_ErrInvalidID(t *testing.T) {
	repo := Constructor()

	err := repo.Delete("2")

	assert.ErrorIs(t, err, ErrInvalidID)
}

func TestStorageRepo_GetAllTask(t *testing.T) {
	repo := Constructor()
	_ = repo.Create(&models.Todo{ID: "1", TaskName: "test1"})
	_ = repo.Create(&models.Todo{ID: "2", TaskName: "test2"})

	todos, err := repo.GetAllTask()
	assert.NoError(t, err)
	assert.Len(t, todos, 2)
}
