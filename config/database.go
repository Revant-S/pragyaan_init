package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	mongoURI := Env.MongoURI
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(50). // not sure
		SetMinPoolSize(10)  // not sure

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
	dbName := Env.DatabaseName
	if dbName == "" {
		dbName = "Super_Hero_Game_DB"
	}

	// Initialize global DB config
	DB = &DatabaseConfig{
		Client:   client,
		Database: client.Database(dbName),
	}
	collections, err := DB.Database.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %v", err)
	}

	// Check if collections list is empty and create a default collection if necessary
	if len(collections) == 0 {
		defaultCollectionName := "default_collection"
		err = DB.Database.CreateCollection(ctx, defaultCollectionName)
		if err != nil {
			return fmt.Errorf("failed to create default collection: %v", err)
		}
		log.Printf("Created default collection: %s", defaultCollectionName)
	}

	log.Printf("Collections in database: %v", collections)

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
