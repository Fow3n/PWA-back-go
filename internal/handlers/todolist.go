package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pwa/internal/models"
	"pwa/internal/repository"
	"time"
)

type TodoListHandler struct {
	Repo *repository.TodoListRepository
}

func (h *TodoListHandler) CreateTodoList(c *gin.Context) {
	var newTodoList models.TodoList
	if err := c.ShouldBindJSON(&newTodoList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTodoList.CreatedAt = time.Now()
	newTodoList.UpdatedAt = time.Now()

	result, err := h.Repo.CreateTodoList(c, newTodoList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

func (h *TodoListHandler) GetTodoList(c *gin.Context) {
	id := c.Param("id")
	todoList, err := h.Repo.FindTodoListByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TodoList not found"})
		return
	}

	c.JSON(http.StatusOK, todoList)
}

func (h *TodoListHandler) UpdateTodoList(c *gin.Context) {
	id := c.Param("id")
	var todoList models.TodoList
	if err := c.ShouldBindJSON(&todoList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todoList.UpdatedAt = time.Now()
	result, err := h.Repo.UpdateTodoList(c, id, todoList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *TodoListHandler) DeleteTodoList(c *gin.Context) {
	id := c.Param("id")
	result, err := h.Repo.DeleteTodoList(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
