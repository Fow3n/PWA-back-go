package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pwa/internal/models"
	"pwa/internal/repository"
	"pwa/internal/service"
	"time"
)

func NewTodoListHandler(repo *repository.TodoListRepository) *TodoListHandler {
	return &TodoListHandler{Repo: repo}
}

type TodoListHandler struct {
	Repo           *repository.TodoListRepository
	WebPushService *service.WebPushService
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

func (h *TodoListHandler) GetTodoListsByUserID(c *gin.Context) {
	userID := c.Param("id")
	todoLists, err := h.Repo.FindTodoListsByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TodoLists not found"})
		return
	}

	c.JSON(http.StatusOK, todoLists)
}

func (h *TodoListHandler) GetTodoListByChannelID(c *gin.Context) {
	channelID := c.Param("id")
	todoLists, err := h.Repo.FindTodoListsByChannelID(c, channelID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TodoLists not found"})
		return
	}

	c.JSON(http.StatusOK, todoLists)
}

func (h *TodoListHandler) AddTask(c *gin.Context) {
	todoListID := c.Param("id")
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	if err := h.Repo.AddTaskToList(c, todoListID, task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Task added"})
}

func (h *TodoListHandler) UpdateTask(c *gin.Context) {
	todoListID := c.Param("todoListId")
	taskID := c.Param("taskId")
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.UpdatedAt = time.Now()

	oldTask, err := h.Repo.GetTaskByID(c, todoListID, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	if err := h.Repo.UpdateTask(c, todoListID, taskID, task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if oldTask.Completed != task.Completed {
		message := fmt.Sprintf("Task '%s' has been marked as %v.", task.Title, task.Completed)
		err := h.WebPushService.NotifyChannelMembers(c, todoListID, message)
		if err != nil {
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func (h *TodoListHandler) DeleteTask(c *gin.Context) {
	todoListID := c.Param("todoListId")
	taskID := c.Param("taskId")

	task, err := h.Repo.GetTaskByID(c, todoListID, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	if err := h.Repo.DeleteTask(c, todoListID, taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := fmt.Sprintf("Task '%s' has been deleted.", task.Title)
	err = h.WebPushService.NotifyChannelMembers(c, todoListID, message)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
