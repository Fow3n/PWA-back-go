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

// CreateChannel godoc
// @Summary Create a new channel
// @Description Adds a new channel to the system with the provided information.
// @Tags channels
// @Accept json
// @Produce json
// @Success 201 {object} object "Successful creation with new channel ID" {id string}
// @Failure 400 {object} ErrorResponse "Bad request when the JSON data is invalid"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /channels [post]
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

// GetChannel godoc
// @Summary Get a channel by ID
// @Description Retrieves a channel by its unique identifier.
// @Tags channels
// @Accept json
// @Produce json
// @Param id path string true "Channel ID"
// @Success 200 {object} Channel
// @Failure 404 {object} ErrorResponse "Channel not found"
// @Router /channels/{id} [get]
func (h *ChannelHandler) GetChannel(c *gin.Context) {
	id := c.Param("id")
	channel, err := h.Repo.FindChannelByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	c.JSON(http.StatusOK, channel)
}

// GetChannelsByUserID godoc
// @Summary Get all channels for a user
// @Description Retrieves a list of all channels that a user is a member of.
// @Tags channels
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {array} Channel
// @Failure 404 {object} ErrorResponse "Channels not found"
// @Router /users/{id}/channels [get]
func (h *ChannelHandler) GetChannelsByUserID(c *gin.Context) {
	userID := c.Param("id")
	channels, err := h.Repo.FindChannelsByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channels not found"})
		return
	}

	c.JSON(http.StatusOK, channels)
}

// UpdateChannel godoc
// @Summary Update a channel
// @Description Updates a channel with the provided information.
// @Tags channels
// @Accept json
// @Produce json
// @Param id path string true "Channel ID"
// @Param channel body Channel true "Channel update data"
// @Success 200 {object} Channel
// @Failure 400 {object} ErrorResponse "Bad request when the JSON data is invalid"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /channels/{id} [put]
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

// DeleteChannel godoc
// @Summary Delete a channel
// @Description Removes a channel from the system.
// @Tags channels
// @Accept json
// @Produce json
// @Param id path string true "Channel ID"
// @Success 200 {object} DeleteResponse
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /channels/{id} [delete]
func (h *ChannelHandler) DeleteChannel(c *gin.Context) {
	id := c.Param("id")
	result, err := h.Repo.DeleteChannel(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
