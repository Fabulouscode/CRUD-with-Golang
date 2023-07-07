package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github/fabulousCode/services/controllers"
	"github/fabulousCode/services/services"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := connectToMongoDB(ctx)
	defer client.Disconnect(ctx)

	userCollection := getMongoDBCollection(client, "got_tv", "users")

	userService := createUserServiceImpl(userCollection, ctx)

	userController := createUserController(userService)

	router := setupRouter(userController)

	startServer(router)
}

func connectToMongoDB(ctx context.Context) *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return client
}

func getMongoDBCollection(client *mongo.Client, dbName string, collectionName string) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}

func createUserServiceImpl(userCollection *mongo.Collection, ctx context.Context) services.UserService {
	return services.NewUserServiceImpl(userCollection, ctx)
}

func createUserController(userService services.UserService) *controllers.UserController {
	return controllers.NewUserController(userService)
}

func setupRouter(userController *controllers.UserController) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	userController.SetupRoutes(api)

	return router
}

func startServer(router *gin.Engine) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
