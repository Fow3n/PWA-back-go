package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"pwa/internal/models"
	"pwa/internal/repository"
	"pwa/pkg/jwt"
	"time"
	_ "time"
)

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

type UserHandler struct {
	Repo *repository.UserRepository
}

// CreateUser godoc
// @Summary Create a new user
// @Description Adds a new user to the system with the provided information.
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} map[string]interface{} "Successful creation with new user ID"
// @Failure 400 {object} map[string]interface{} "Bad request when the JSON data is invalid"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid user data")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to process password")
		return
	}
	newUser.Password = string(hashedPassword)
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	result, err := h.Repo.CreateUser(c, newUser)
	if mongo.IsDuplicateKeyError(err) {
		respondWithError(c, http.StatusBadRequest, "Username or email already exists")
		return
	} else if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to register user")
		return
	}

	objID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		respondWithError(c, http.StatusInternalServerError, "Failed to generate user ID")
		return
	}

	token, err := jwt.GenerateToken(objID.Hex())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "token": token})
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieves a list of all users in the system.
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 404 {object} ErrorResponse "Users not found"
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.Repo.FindUsers(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieves a user by their unique identifier.
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} ErrorResponse "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.Repo.FindUserByIdentifier(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Updates a user's information with the provided data.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body User true "User update data"
// @Success 200 {object}
// @Failure 400 {object} ErrorResponse "Bad request when the JSON data is invalid"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	updateDoc := bson.M{"$set": bson.M{}}
	for key, value := range updateData {
		updateDoc["$set"].(bson.M)[key] = value
		if key == "password" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(value.(string)), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
				return
			}
			updateDoc["$set"].(bson.M)[key] = string(hashedPassword)
		}
	}
	updateDoc["$set"].(bson.M)["updatedAt"] = time.Now()

	objID, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"_id": objID}
	result, err := h.Repo.UpdateUser(c, filter, updateDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result.ModifiedCount})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Removes a user from the system.
// @Tags users
// @Param id path string true "User ID"
// @Success 200 {object}
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	result, err := h.Repo.DeleteUser(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted_count": result.DeletedCount})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var loginDetails models.LoginRequest
	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}

	user, err := h.Repo.FindUserByIdentifier(c, loginDetails.Identifier)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDetails.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := jwt.GenerateToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}
