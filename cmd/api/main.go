package main

import (
	"log"
	"net/http"
	"os"
	"pwa/internal/middleware"
	"pwa/pkg/mongodb"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"pwa/internal/handlers"
	"pwa/internal/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoClient := setupMongoClient()

	router := setupRouter()
	setupRoutes(router, mongoClient)
	configureCORS(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupMongoClient() *mongo.Client {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI must be set in .env")
	}

	mongoClient, err := mongodb.NewClient(mongoURI)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}
	return mongoClient
}

func configureCORS(router *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))
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
	channelRepo := &repository.ChannelRepository{Collection: client.Database("pwa").Collection("channels")}
	channelHandler := handlers.NewChannelHandler(channelRepo)
	todoListRepo := &repository.TodoListRepository{Collection: client.Database("pwa").Collection("todoLists")}
	todoListHandler := handlers.NewTodoListHandler(todoListRepo)
	notificationRepo := &repository.WebPushRepository{Collection: client.Database("pwa").Collection("webPushSubscriptions")}
	notificationHandler := handlers.NewNotificationHandler(notificationRepo)

	router.POST("/login", userHandler.LoginUser)
	router.POST("/users", userHandler.CreateUser)

	userRoutes := router.Group("/users")
	userRoutes.Use(middleware.JWTAuthMiddleware())
	{
		userRoutes.GET("/", userHandler.GetUsers)
		userRoutes.GET("/:id", userHandler.GetUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	channelRoutes := router.Group("/channels")
	channelRoutes.Use(middleware.JWTAuthMiddleware())
	{
		channelRoutes.POST("/", channelHandler.CreateChannel)
		channelRoutes.GET("/:id", channelHandler.GetChannel)
		channelRoutes.GET("/users/:id", channelHandler.GetChannelsByUserID)
		channelRoutes.PUT("/:id", channelHandler.UpdateChannel)
		channelRoutes.DELETE("/:id", channelHandler.DeleteChannel)
		channelRoutes.POST("/:id/join", channelHandler.JoinChannel)
		channelRoutes.POST("/:id/leave", channelHandler.LeaveChannel)
	}

	todoListRoutes := router.Group("/todoLists")
	todoListRoutes.Use(middleware.JWTAuthMiddleware())
	{
		todoListRoutes.POST("/:id/tasks", todoListHandler.AddTask)
		todoListRoutes.PUT("/:todoListId/tasks/:taskId", todoListHandler.UpdateTask)
		todoListRoutes.DELETE("/:todoListId/tasks/:taskId", todoListHandler.DeleteTask)
		todoListRoutes.GET("/channels/:id", todoListHandler.GetTodoListByChannelID)
		todoListRoutes.GET("/:id", todoListHandler.GetTodoList)
		todoListRoutes.PUT("/:id", todoListHandler.UpdateTodoList)
		todoListRoutes.DELETE("/:id", todoListHandler.DeleteTodoList)
		todoListRoutes.POST("/", todoListHandler.CreateTodoList)
	}

	router.POST("/subscribe", notificationHandler.Subscribe)
	router.POST("/unsubscribe/:id", notificationHandler.Unsubscribe)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API is running",
		})
	})
}
