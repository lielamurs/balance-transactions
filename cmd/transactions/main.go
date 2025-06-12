package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lielamurs/balance-transactions/internal/database"
	"github.com/lielamurs/balance-transactions/internal/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	setupLogger()

	database.Init()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userHandler := handler.NewUserHandler()

	e.GET("/user/:userId/balance", userHandler.GetBalance)
	e.POST("/user/:userId/transaction", userHandler.ProcessTransaction)

	logrus.Info("Starting server on :8080")
	log.Fatal(e.Start(":8080"))
}

func setupLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)

	logrus.Info("Logger configured successfully")
}
