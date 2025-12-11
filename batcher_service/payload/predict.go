package payload

type PredictionRequest struct {
	Instances []string `json:"instances"`
}

type PredictionResponse struct {
	Predictions []string `json:"predictions"`
}
