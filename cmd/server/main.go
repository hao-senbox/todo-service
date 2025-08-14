package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-service/config"
	"todo-service/internal/todo"
	"todo-service/internal/user"
	"todo-service/pkg/consul"
	"todo-service/pkg/zap"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// Load env
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		} else {
			log.Println("Successfully loaded .env file")
		}
	} else {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.LoadConfig()

	logger, err := zap.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	mongoClient, err := connectToMongoDB(cfg.MongoURI)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	consulConn := consul.NewConsulConn(logger, cfg)
	consulClient := consulConn.Connect()
	userService := user.NewUserService(consulClient)
	todoCollection := mongoClient.Database(cfg.MongoDB).Collection("todo")
	todoRepository := todo.NewTodoRepository(todoCollection)
	todoService := todo.NewTodoService(todoRepository, userService)
	todoHandler := todo.NewTodoHandler(todoService)

	r := gin.Default()

	todo.RegisterRoutes(r, todoHandler)
	// Handle OS signal để deregister
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server... De-registering from Consul...")
		consulConn.Deregister()
		os.Exit(0)
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8011"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Failed to connect to MongoDB")
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Failed to ping MongoDB")
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}
