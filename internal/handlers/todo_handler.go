package handlers

import (
	"github.com/gin-gonic/gin"

	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/internal/services"
)

type TodoHandler struct {
	service services.TodoService
}

func NewTodoHandler(service services.TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var request models.CreateTodoRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "неверный JSON"})
		return
	}

	task, err := h.service.CreateTodo(&request)
	if err != nil {
		switch err {
		case repository.ErrEmptyID:
			c.JSON(400, gin.H{"error": err.Error()})
			return
		case repository.ErrEmptyData:
			c.JSON(400, gin.H{"error": err.Error()})
			return
		case repository.ErrAllreadyExist:
			c.JSON(409, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(500, gin.H{"error": "внутренняя ошибка сервера"})
			return
		}
	}

	c.JSON(201, task)
}

func (h *TodoHandler) GetById(c *gin.Context) {
	id := c.Param("id")

	task, err := h.service.GetById(id)

	if err != nil {
		switch err {
		case repository.ErrEmptyID:
			c.JSON(400, gin.H{"error": err.Error()})
			return
		case repository.ErrInvalidID:
			c.JSON(404, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(500, gin.H{"error": "внутренняя ошибка сервера"})
			return
		}
	}

	c.JSON(200, task)
}

func (h *TodoHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var updateData models.UpdateTodoRequest

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{"error": "неверный JSON"})
		return
	}

	task, err := h.service.UpdateTodo(id, &updateData)

	if err != nil {
		switch err {
		case repository.ErrEmptyID:
			c.JSON(400, gin.H{"error": err.Error()})
			return
		case repository.ErrInvalidID:
			c.JSON(404, gin.H{"error": err.Error()})
			return
		case repository.ErrEmptyData:
			c.JSON(400, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(500, gin.H{"error": "внутренняя ошибка сервера"})
			return
		}
	}

	c.JSON(200, task)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteTodo(id)

	if err != nil {
		switch err {
		case repository.ErrEmptyID:
			c.JSON(400, gin.H{"error": err.Error()})
			return
		case repository.ErrInvalidID:
			c.JSON(404, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(500, gin.H{"error": "внутренняя ошибка сервера"})
			return
		}
	}

	c.JSON(204, nil)
}

func (h *TodoHandler) GetAllTask(c *gin.Context) {
	tasks, err := h.service.GetAllTodos()

	if err != nil {
		c.JSON(500, gin.H{"error": "внутренняя ошибка сервера"})
		return
	}

	c.JSON(200, tasks)
}
