package controllers

import (
	"fmt"
	"main/internal/models"
	"main/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
}

// @title Golang MVC Backend API
// @version 1.0
// @description This is a sample API for user management.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support Team
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8081
// @BasePath /api/v1
// @accept json
// @produce json

// @tag.name Users
// @tag.description Operations related to user management.
// @tag.docs.url https://swagger.io
// @tag.docs.description Swagger Documentation for User Operations.

// @Summary Sign up a new user
// @Description This endpoint allows a new user to register with the system.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.User true "User signup data"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/signup [post]
func (uc *UserController) Signup(c echo.Context) error {

	var newUser models.User
	if err := c.Bind(&newUser); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "Invalid request body")
	}
	if err := utils.ValidateSignupInput(newUser); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, err.Error())
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "Error processing password")
	}

	existingUser, _ := utils.GetUserByEmail(newUser.Email)

	if existingUser != (models.User{}) {
		return utils.JSONResponse(c, http.StatusBadRequest, "User already exists")
	}

	newUser.Password = string(hashedPassword)
	savedUser, err := utils.CreateUser(newUser)
	if err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "Failed to create user")
	}

	message := fmt.Sprintf("User registered successfully with id %s", savedUser.InsertedID)
	return utils.JSONResponse(c, http.StatusCreated, message)
}
