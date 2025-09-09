package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func RunMigrations(db *sql.DB) error {
	file, err := os.ReadFile("migrations/001_create_todos.sql")
	if err != nil {
		return err
	}

	sqlQuery := string(file)

	_, err = db.Exec(sqlQuery)

	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			fmt.Println("Таблица todos уже существует")
			return nil
		}
		return fmt.Errorf("ошибка выполнения миграции: %w", err)
	}

	return nil
}
