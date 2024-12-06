package main

import (
	"context"
	"fmt"
	"log"
	"main/config"
	"main/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"main/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Super Hero Game MVC Backend
// @version 1.0
// @description This is a Game server for Super Hero Game.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default configurations")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := config.InitializeDatabase(ctx); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	fmt.Println("Connected to MongoDB!")

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/swagger/*"
		},
	}))

	e.Use(middleware.Recover())
	docs.SwaggerInfo.BasePath = "/v1"
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	userRoutes := routes.UserRoutes{}
	userRoutes.SetupRoutes(e)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
	}

	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Printf("Server started on port %s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Failed to shut down server: %v", err)
	}
	log.Println("Server exited gracefully")
}
