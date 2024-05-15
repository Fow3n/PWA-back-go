package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pwa/internal/models"
	"pwa/internal/repository" // Ensure this import is correct
)

type NotificationHandler struct {
	Repo *repository.WebPushRepository
}

func NewNotificationHandler(repo *repository.WebPushRepository) *NotificationHandler {
	return &NotificationHandler{
		Repo: repo,
	}
}

func (h *NotificationHandler) Subscribe(c *gin.Context) {
	var subscription models.WebPushSubscription
	if err := c.ShouldBindJSON(&subscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription data", "details": err.Error()})
		return
	}

	if _, err := h.Repo.CreateWebPushSubscription(c.Request.Context(), subscription); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully subscribed to notifications"})
}

func (h *NotificationHandler) Unsubscribe(c *gin.Context) {
	subscriptionID := c.Param("id")

	if _, err := h.Repo.DeleteWebPushSubscription(c.Request.Context(), subscriptionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unsubscribed from notifications"})
}
