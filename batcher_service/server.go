package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {

}

type PredictionRequest struct {
	Instances []string `json:"instances"`
}

type PredictionResponse struct {
	Predictions []string `json:"predictions"`
}

func (s *Server) CreatePredictions(c echo.Context) error {
	req := new(PredictionRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	var predictions []string
	for _, instance := range req.Instances {
		predictions = append(predictions, instance+"_response")
	}
	res := PredictionResponse{predictions}
	return c.JSON(http.StatusCreated, res)
}
