package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"pwa/internal/handlers"
	"pwa/internal/repository"
	"pwa/pkg/mongodb"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI must be set in .env")
	}

	mongoClient, err := mongodb.NewClient(mongoURI)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	router := setupRouter()
	setupRoutes(router, mongoClient)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	trustedProxies := []string{"192.168.1.0/24"}
	if err := router.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	return router
}

func setupRoutes(router *gin.Engine, client *mongo.Client) {
	userRepo := &repository.UserRepository{Collection: client.Database("pwa").Collection("users")}
	userHandler := handlers.NewUserHandler(userRepo)

	router.POST("/users", userHandler.CreateUser)
	router.GET("/users/:id", userHandler.GetUser)
	router.GET("/users", userHandler.GetUsers)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)
	router.POST("/login", userHandler.LoginUser)

	channelRepo := &repository.ChannelRepository{Collection: client.Database("pwa").Collection("channels")}
	channelHandler := handlers.NewChannelHandler(channelRepo)

	router.POST("/channels", channelHandler.CreateChannel)
	router.GET("/channels/:id", channelHandler.GetChannel)
	router.GET("/users/:id/channels", channelHandler.GetChannelsByUserID)
	router.PUT("/channels/:id", channelHandler.UpdateChannel)
	router.DELETE("/channels/:id", channelHandler.DeleteChannel)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API is running",
		})
	})
}
