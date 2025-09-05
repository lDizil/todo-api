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

// @Summary Создать задачу
// @Description Создание новой задачи
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body models.CreateTodoRequest true "Данные задачи"
// @Success 201 {object} models.Todo
// @Failure 400 {object} map[string]string "Неверный формат ID"
// @Failure 409 {object} map[string]string "Задача с таким айди уже существует"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /todos [post]
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
		case repository.ErrАlreadyExist:
			c.JSON(409, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(500, gin.H{"error": "внутренняя ошибка сервера"})
			return
		}
	}

	c.JSON(201, task)
}

// @Summary Получить задачу
// @Description Получение задачи по её ID
// @Tags todos
// @Produce json
// @Param id path string true "ID задачи"
// @Success 200 {object} models.Todo
// @Failure 400 {object} map[string]string "Неверный формат ID"
// @Failure 404 {object} map[string]string "Задача не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /todos/{id} [get]
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

// @Summary Обновить задачу
// @Description Обновление задачи по ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Param todo body models.UpdateTodoRequest true "Данные для обновления"
// @Success 200 {object} models.Todo
// @Failure 400 {object} map[string]string "Неверный формат ID или данных для обновления"
// @Failure 404 {object} map[string]string "Задача не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /todos/{id} [patch]
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

// @Summary Удалить задачу
// @Description Удаление задачи по айди
// @Tags todos
// @Produce json
// @Param id path string true "ID задачи"
// @Success 204 "Задача успешно удалена"
// @Failure 400 {object} map[string]string "Неверный формат ID"
// @Failure 404 {object} map[string]string "Задача не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /todos/{id} [delete]
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

	c.Status(204)
}

// @Summary Получить все задачи
// @Description Получения списка всех задач
// @Tags todos
// @Produce json
// @Success 200 {array} models.Todo
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /todos [get]
func (h *TodoHandler) GetAllTask(c *gin.Context) {
	tasks, err := h.service.GetAllTodos()

	if err != nil {
		c.JSON(500, gin.H{"error": "внутренняя ошибка сервера"})
		return
	}

	c.JSON(200, tasks)
}
