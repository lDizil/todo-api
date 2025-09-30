package database

import (
	"database/sql"
	"fmt"
	"time"
	"todo-api/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.DSN()
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("не удалось открыть соединение с БД: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	fmt.Printf("Успешное подключение к базе данных %s:%s\n",
		cfg.Host, cfg.Port)

	return db, nil
}

func ConnectTestDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.TestDSN()
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("не удалось открыть соединение с БД: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	fmt.Printf("Успешное подключение к базе данных %s:%s\n",
		cfg.Host, cfg.Port)

	return db, nil
}
