package main

import (
	"context"
	"fmt"
	"log"
	"main/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default configurations")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := config.InitializeDatabase(ctx); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

			router := http.NewServeMux()

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Server starting on port %s", port)
	testCRUD(ctx)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

}

func testCRUD(ctx context.Context) {
	collection := config.GetCollection("users")
	fmt.Println("Connected to MongoDB!")
	// Create
	newUser := bson.M{
		"name":  "John Doe",
		"email": "john.doe@example.com",
		"age":   30,
	}
	insertResult, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}
	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)

	// Read
	var result bson.M
	err = collection.FindOne(ctx, bson.M{"name": "John Doe"}).Decode(&result)
	if err != nil {
		log.Fatalf("Failed to find document: %v", err)
	}
	fmt.Printf("Found document: %v\n", result)

	// Update
	update := bson.M{
		"$set": bson.M{
			"age": 35,
		},
	}
			updateResult, err := collection.UpdateOne(ctx, bson.M{"name": "John Doe"}, update)
	if err != nil {
		log.Fatalf("Failed to update document: %v", err)
	}
	fmt.Printf("Matched %d document(s) and updated %d document(s)\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}
