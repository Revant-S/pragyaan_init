package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"main/config"
	"main/internal/models"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateUniqueID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("Error generating user ID:", err)
		return ""
	}
	return fmt.Sprintf("%x", b)
}

func CreateUser(user models.User) (mongo.InsertOneResult, error) {
	coll := config.GetCollection("users")

	result, err := coll.InsertOne(context.TODO(), user)

	if err != nil {
		return mongo.InsertOneResult{}, err
	}

	return *result, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	coll := config.GetCollection("users")
	err := coll.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, fmt.Errorf("no user found with email: %s", email)
		}
		return user, err
	}

	return user, nil
}

func ValidateSignupInput(user models.User) error {

	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)

	if user.Username == "" {
		return fmt.Errorf("username is required")
	}
	if len(user.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(user.Email) {
		return fmt.Errorf("invalid email format")
	}
	fmt.Println("HERE IS THE PASSWORD")
	fmt.Println(user)
	if len(user.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	return nil
}
