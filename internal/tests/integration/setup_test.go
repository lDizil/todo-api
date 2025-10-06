package integration_tests

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"todo-api/internal/config"
	"todo-api/internal/models"
	"todo-api/migrations"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	cfg := config.Load()

	var err error
	testDB, err = sql.Open("postgres", cfg.Database.TestDSN())
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := testDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err := migrations.RunMigrations(testDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	code := m.Run()

	CleanUpDatabase()
	testDB.Close()

	os.Exit(code)
}

func SetUpTest(t *testing.T) *sql.DB {
	t.Helper()

	_, err := testDB.Exec("TRUNCATE TABLE todos RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Failed to truncate table: %v", err)
	}

	return testDB
}

func CleanUpDatabase() {
	testDB.Exec("DROP TABLE IF EXISTS todos CASCADE")
}

func CreateTestTodo(db *sql.DB, taskName string, description *string) *models.Todo {
	todo := &models.Todo{
		TaskName:    taskName,
		Description: description,
		Completed:   false,
	}

	query := "INSERT INTO todos (task_name, description, completed) VALUES ($1, $2, $3) RETURNING id, created_at"
	err := db.QueryRow(query, todo.TaskName, todo.Description, todo.Completed).Scan(&todo.ID, &todo.CreatedAt)

	if err != nil {
		panic(err)
	}

	return todo
}
