package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"main/config"
	"main/internal/models"
	"main/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)


func TestMain(m *testing.M) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("DATABASE_NAME", "test_database")

	// Initialize database
	if err := config.InitializeDatabase(ctx); err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	// Clean up test database before running tests
	cleanupTestDatabase(ctx)

	// Run tests
	code := m.Run()

	// Cleanup
	if err := config.CloseConnection(ctx); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	os.Exit(code)
}

// cleanupTestDatabase removes all documents from user collection
func cleanupTestDatabase(ctx context.Context) {
	userCollection := config.GetCollection("users")
	_, err := userCollection.DeleteMany(ctx, bson.D{})
	if err != nil {
		log.Fatalf("Failed to clean up test database: %v", err)
	}
}

func setupTestEnvironment() (*echo.Echo, *httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users/signup", nil)
	req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
	return e, recorder, e.NewContext(req, recorder)
}

func TestSignup_Successful(t *testing.T) {
	// Clean up before test
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cleanupTestDatabase(ctx)

	// Setup test environment
	e, recorder, c := setupTestEnvironment()

	// Prepare user data
	userData := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	payload, err := json.Marshal(userData)
	assert.NoError(t, err)

	// Update request with payload
	req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader(payload))
	req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
	c = e.NewContext(req, recorder)

	// Create controller and call Signup
	uc := UserController{}
	err = uc.Signup(c)

	// Assertions
	assert.NoError(t, err, "Signup should not return an error")
	assert.Equal(t, http.StatusCreated, recorder.Code, "Expected status code 201")

	// Verify user was created in database
	createdUser, err := utils.GetUserByEmail(userData.Email)
	assert.NoError(t, err, "Should be able to retrieve created user")
	assert.NotEmpty(t, createdUser.ID, "Created user should have an ID")
	assert.Equal(t, userData.Username, createdUser.Username, "Username should match")
	assert.Equal(t, userData.Email, createdUser.Email, "Email should match")

	// Verify password was hashed
	err = bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(userData.Password))
	assert.NoError(t, err, "Password should be correctly hashed")
}

func TestSignup_ExistingUser(t *testing.T) {
	// Clean up before test
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cleanupTestDatabase(ctx)

	// Create an existing user
	existingUser := models.User{
		Username: "existinguser",
		Email:    "existing@example.com",
		Password: "password123",
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(existingUser.Password), bcrypt.DefaultCost)
	existingUser.Password = string(hashedPassword)
	_, err := utils.CreateUser(existingUser)
	assert.NoError(t, err, "Should be able to create existing user")

	// Setup test environment
	e, recorder, c := setupTestEnvironment()

	// Prepare duplicate user data
	payload, err := json.Marshal(existingUser)
	assert.NoError(t, err)

	// Update request with payload
	req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader(payload))
	req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
	c = e.NewContext(req, recorder)

	// Create controller and call Signup
	uc := UserController{}
	err = uc.Signup(c)

	// Assertions
	assert.Error(t, err, "Signup should return an error for existing user")
	assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
}

func TestSignup_InvalidInput(t *testing.T) {
	testCases := []struct {
		name          string
		userData      models.User
		expectedError string
	}{
		{
			name: "Empty Username",
			userData: models.User{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedError: "username is required",
		},
		{
			name: "Invalid Email",
			userData: models.User{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
			},
			expectedError: "invalid email format",
		},
		{
			name: "Short Password",
			userData: models.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "short",
			},
			expectedError: "password must be at least 8 characters long",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test environment
			e, recorder, c := setupTestEnvironment()

			// Prepare user data
			payload, err := json.Marshal(tc.userData)
			assert.NoError(t, err)

			// Update request with payload
			req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader(payload))
			req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
			c = e.NewContext(req, recorder)

			// Create controller and call Signup
			uc := UserController{}
			err = uc.Signup(c)

			// Assertions
			assert.Error(t, err, "Signup should return an error for invalid input")
			assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
		})
	}
}