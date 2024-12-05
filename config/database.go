package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConfig struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var (
	DB *DatabaseConfig
)

func InitializeDatabase(ctx context.Context) error {
	// Retrieve MongoDB connection URI
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(50).
		SetMinPoolSize(10)

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	// Get database name from environment or use default
	dbName := os.Getenv("MONGO_DATABASE")
	if dbName == "" {
		dbName = "default_database"
	}

	// Initialize global DB config
	DB = &DatabaseConfig{
		Client:   client,
		Database: client.Database(dbName),
	}

	log.Printf("Successfully connected to MongoDB database: %s", dbName)
	return nil
}

func GetCollection(collectionName string) *mongo.Collection {
	if DB == nil || DB.Database == nil {
		log.Fatal("Database not initialized")
	}
	return DB.Database.Collection(collectionName)
}

func CloseConnection(ctx context.Context) error {
	if DB == nil || DB.Client == nil {
		return nil
	}

	return DB.Client.Disconnect(ctx)
}
