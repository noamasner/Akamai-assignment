package main

import (
	"batcher_service/batcher"
	"batcher_service/payload"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	http.Server
	batcher *batcher.Batcher
}

func NewServer(batcher *batcher.Batcher) *Server {
	return &Server{
		Server:  http.Server{},
		batcher: batcher,
	}
}

func (s *Server) Start() {

}

func (s *Server) CreatePredictions(c echo.Context) error {
	req := new(payload.PredictionRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	predictions, err := s.collectPredications(c, req.Instances)
	if err != nil {
		c.Error(err)
	}
	res := payload.PredictionResponse{
		Predictions: predictions,
	}
	return c.JSON(http.StatusCreated, &res)
}

func (s *Server) collectPredications(c echo.Context, instances []string) ([]string, error) {
	resCh := make(chan batcher.PredictResult, len(instances))
	for i, instance := range instances {
		message := batcher.PredictMessage{
			Instance: instance,
			Seq:      i,
			Size:     len(instances),
			ResCh:    resCh,
		}
		err := s.batcher.PushMessage(c.Request().Context(), message)
		if err != nil {
			return nil, err
		}
	}

	var predictions []string
	for r := range resCh {
		if r.Err != nil {
			return nil, r.Err
		}
		predictions = append(predictions, r.Pred)
	}
	return predictions, nil
}
