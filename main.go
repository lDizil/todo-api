package main

import (
	_ "todo-api/docs"
	"todo-api/internal/handlers"
	"todo-api/internal/repository"
	"todo-api/internal/services"

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
	repo := repository.Constructor()

	service := services.NewTodoService(repo)

	handlers := handlers.NewTodoHandler(service)

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

	router.Run(":8080")
}
