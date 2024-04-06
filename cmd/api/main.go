package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
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

	router := gin.Default()
	trustedProxies := []string{"192.168.1.0/24"}
	err = router.SetTrustedProxies(trustedProxies)
	if err != nil {
		return
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is running",
		})
	})
	repo := &repository.UserRepository{Collection: mongoClient.Database("pwa").Collection("users")}
	userHandler := &handlers.UserHandler{Repo: repo}

	router.POST("/users", userHandler.CreateUser)
	router.GET("/users/:id", userHandler.GetUser)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)

	err = router.Run(":8080")
	if err != nil {
		return
	}
}
