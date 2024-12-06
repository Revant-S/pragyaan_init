package routes

import (
	"main/internal/controllers"

	"github.com/labstack/echo/v4"
)

type UserRoutes struct {
	UserController controllers.UserController
}

func (ur *UserRoutes) SetupRoutes(e *echo.Echo) {
	usersGroup := e.Group("/users")
	usersGroup.POST("/signup", ur.UserController.Signup)
}
