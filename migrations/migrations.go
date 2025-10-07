package migrations

import (
	"database/sql"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed 001_create_todos.sql
var createTodosSQL string

func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(createTodosSQL)

	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			fmt.Println("Таблица todos уже существует")
			return nil
		}
		return fmt.Errorf("ошибка выполнения миграции: %w", err)
	}

	return nil
}
