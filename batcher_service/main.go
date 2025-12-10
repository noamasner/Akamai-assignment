package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	server := NewServer()
	server.Start()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/v1/predict", server.CreatePredictions)
	e.Logger.Fatal(e.Start(":30001"))
}
