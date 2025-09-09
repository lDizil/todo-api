package main

import (
	"log"
	_ "todo-api/docs"
	"todo-api/internal/config"
	"todo-api/internal/database"
	"todo-api/internal/handlers"
	"todo-api/internal/repository"
	"todo-api/internal/services"
	"todo-api/migrations"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title TODO API
// @version 1.0
// @description API для управления задачами
// @host localhost:8080
// @BasePath /
func main() {

	cfg := config.Load()

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal("ошибка при подключении к базе данных: ", err)
	}

	defer db.Close()

	err = migrations.RunMigrations(db)
	if err != nil {
		log.Fatal("ошибка миграции: ", err)
	}

	repo := repository.NewPostgresRepository(db)

	service := services.NewTodoService(repo)

	handlers := handlers.NewTodoHandler(service)

	gin.SetMode(cfg.Server.Mode)
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	todosGroup := router.Group("/todos")
	{
		todosGroup.POST("", handlers.CreateTodo)
		todosGroup.GET("", handlers.GetAllTask)
		todosGroup.GET("/:id", handlers.GetById)
		todosGroup.PATCH("/:id", handlers.Update)
		todosGroup.DELETE("/:id", handlers.Delete)
	}

	router.Run(":" + cfg.Server.Port)
}
