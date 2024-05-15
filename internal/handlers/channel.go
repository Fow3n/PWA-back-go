package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pwa/internal/models"
	"pwa/internal/repository"
	"time"
)

func NewChannelHandler(repo *repository.ChannelRepository) *ChannelHandler {
	return &ChannelHandler{Repo: repo}
}

type ChannelHandler struct {
	Repo *repository.ChannelRepository
}

func (h *ChannelHandler) CreateChannel(c *gin.Context) {
	var newChannel models.Channel
	if err := c.ShouldBindJSON(&newChannel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newChannel.CreatedAt = time.Now()
	newChannel.UpdatedAt = time.Now()

	result, err := h.Repo.CreateChannel(c, newChannel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

func (h *ChannelHandler) GetChannel(c *gin.Context) {
	id := c.Param("id")
	channel, err := h.Repo.FindChannelByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	c.JSON(http.StatusOK, channel)
}

func (h *ChannelHandler) GetChannelsByUserID(c *gin.Context) {
	userID := c.Param("id")
	channels, err := h.Repo.FindChannelsByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channels not found"})
		return
	}

	c.JSON(http.StatusOK, channels)
}

func (h *ChannelHandler) UpdateChannel(c *gin.Context) {
	id := c.Param("id")
	var channel models.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channel.UpdatedAt = time.Now()
	result, err := h.Repo.UpdateChannel(c, id, channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ChannelHandler) DeleteChannel(c *gin.Context) {
	id := c.Param("id")
	result, err := h.Repo.DeleteChannel(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ChannelHandler) JoinChannel(c *gin.Context) {
	var request struct {
		Password string `json:"password"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	userID := c.GetString("userID")
	channelID := c.Param("id")

	ok, err := h.Repo.CheckChannelPassword(c, channelID, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify channel password", "details": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid channel password"})
		return
	}

	if err := h.Repo.JoinChannel(c, channelID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join channel", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined channel"})
}

func (h *ChannelHandler) LeaveChannel(c *gin.Context) {
	userID := c.GetString("userID")
	channelID := c.Param("id")

	if err := h.Repo.LeaveChannel(c, channelID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave channel", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully left channel"})
}
