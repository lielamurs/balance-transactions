package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lielamurs/balance-transactions/internal/database"
	"github.com/lielamurs/balance-transactions/internal/handler"
)

func main() {
	// Initialize database
	database.Init()

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Create handlers
	userHandler := handler.NewUserHandler()

	// Routes
	e.GET("/user/:userId/balance", userHandler.GetBalance)

	// Start server
	log.Println("Starting server on :8080")
	log.Fatal(e.Start(":8080"))
}
