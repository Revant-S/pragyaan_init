package utils

import (
	"github.com/labstack/echo/v4"
)

func JSONResponse(c echo.Context, statusCode int, message interface{}) error {
	return c.JSON(statusCode, map[string]interface{}{
		"message": message,
	})
}
