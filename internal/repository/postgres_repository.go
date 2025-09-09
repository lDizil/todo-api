package repository

import (
	"database/sql"
	"errors"
	"fmt"
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

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) TodoRepository {
	return &PostgresRepository{
		db: db,
	}
}

var ErrEmptyID = errors.New("передан пустой айди")
var ErrInvalidID = errors.New("задача с таким айди не найдена")
var ErrEmptyTask = errors.New("передана пустая задача")
var ErrEmptyData = errors.New("переданы пустые данные")
var ErrAlreadyExist = errors.New("задача с таким айди уже существует")
var ErrEmptyName = errors.New("необходимо передать наименование задачи")

func (r *PostgresRepository) Create(task *models.Todo) error {
	query := "INSERT INTO todos (task_name, description, completed) VALUES ($1, $2, $3) RETURNING id, created_at"

	err := r.db.QueryRow(query, task.TaskName, task.Description, task.Completed).Scan(&task.ID, &task.CreatedAt)
	return err
}

func (r *PostgresRepository) Update(id string, updateData *models.UpdateTodoRequest) error {
	args := []interface{}{}
	setParts := []string{}
	argIndex := 1

	if updateData.TaskName != nil {
		setParts = append(setParts, fmt.Sprintf("task_name = $%d", argIndex))
		args = append(args, *updateData.TaskName)
		argIndex++
	}

	if updateData.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *updateData.Description)
		argIndex++
	}

	if updateData.Completed != nil {
		setParts = append(setParts, fmt.Sprintf("completed = $%d", argIndex))
		args = append(args, *updateData.Completed)
		argIndex++
	}

	if len(setParts) == 0 {
		return ErrEmptyData
	}

	query := fmt.Sprintf("UPDATE todos SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *PostgresRepository) GetById(id string) (*models.Todo, error) {
	query := "SELECT id, task_name, description, completed, created_at FROM todos WHERE id = $1"
	row := r.db.QueryRow(query, id)

	var todo models.Todo
	err := row.Scan(&todo.ID, &todo.TaskName, &todo.Description, &todo.Completed, &todo.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrInvalidID
	}
	if err != nil {
		return nil, err
	}

	return &todo, nil

}

func (r *PostgresRepository) GetAllTask() ([]*models.Todo, error) {
	query := "SELECT id, task_name, description, completed, created_at FROM todos"
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Todo

	for rows.Next() {
		todo := &models.Todo{}
		err := rows.Scan(&todo.ID, &todo.TaskName, &todo.Description, &todo.Completed, &todo.CreatedAt)

		if err != nil {
			return nil, err
		}

		result = append(result, todo)
	}

	return result, nil
}

func (r *PostgresRepository) Delete(id string) error {
	query := "DELETE FROM todos WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}
