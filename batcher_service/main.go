package main

import (
	"batcher_service/batcher"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	b := batcher.NewBatcher("http://localhost:8080/api/v1/predict", 5, 100*time.Millisecond)
	b.Start()
	server := NewServer(b)
	server.Start()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/v1/predict", server.CreatePredictions)
	e.Logger.Fatal(e.Start(":30001"))
}
