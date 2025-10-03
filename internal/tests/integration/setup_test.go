package integration_tests

import (
	"log"
	"testing"
	"todo-api/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	cfg := config.Load()

	testDB, err := gorm.Open(postgres.Open(cfg.Database.TestDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent)
	})

	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := testDB.AutoMigrate(&models.Todo{}); err != nil {
        log.Fatalf("Failed to migrate test database: %v", err)
    }

	code := m.Run()
	
	cleanupDatabase()

	os.Exit(code)
}
