package batcher

import (
	"batcher_service/payload"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PredictMessage struct {
	Instance string
	Seq      int
	Size     int
	ResCh    chan PredictResult
}

type PredictResult struct {
	Pred string
	Err  error
}

type Batcher struct {
	in         chan PredictMessage
	predURL    string
	maxBatch   int
	maxWait    time.Duration
	httpClient *http.Client
}

func NewBatcher(predURL string, maxBatch int, maxWait time.Duration) *Batcher {
	b := &Batcher{
		in:         make(chan PredictMessage, 1024), // buffered so callers don't block much
		maxBatch:   maxBatch,
		maxWait:    maxWait,
		predURL:    predURL,
		httpClient: &http.Client{Timeout: 2 * time.Second},
	}
	return b
}

func (b *Batcher) Start() {
	go b.loop()
}

func (b *Batcher) PushMessage(ctx context.Context, msg PredictMessage) error {
	select {
	case b.in <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *Batcher) loop() {
	for {
		// wait for the first req
		firstReq, ok := <-b.in
		if !ok {
			return // channel closed, exit
		}
		batch := []PredictMessage{firstReq}
		if ok := b.collect(batch); !ok {
			// channel is closed. exit the loop
			return
		}
	}
}

// Collects a batch until capacity is reached and then flushes
// Returns true if the channel is still open
func (b *Batcher) collect(batch []PredictMessage) bool {
	for {
		select {
		case req, ok := <-b.in:
			if !ok {
				// channel closed, we want to exit the loop
				b.flush(batch)
				return false
			}
			batch = append(batch, req)
			if len(batch) == b.maxBatch {
				b.flush(batch)
				return true
			}
		case <-time.After(b.maxWait):
			b.flush(batch)
			return true
		}
	}
}

func (b *Batcher) flush(batch []PredictMessage) {
	if len(batch) == 0 {
		return
	}

	// Prepare request payload
	instances := make([]string, len(batch))
	for i, r := range batch {
		instances[i] = r.Instance
	}
	reqBody := payload.PredictionRequest{
		Instances: instances,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		b.fail(batch, fmt.Errorf("marshal request: %w", err))
		return
	}

	resp, err := b.httpClient.Post(b.predURL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		b.fail(batch, fmt.Errorf("http call: %w", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b.fail(batch, fmt.Errorf("prediction_api status: %s", resp.Status))
		return
	}

	var respBody payload.PredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		b.fail(batch, fmt.Errorf("decode response: %w", err))
		return
	}

	if len(respBody.Predictions) != len(batch) {
		b.fail(batch, fmt.Errorf("mismatched lengths: got %d predictions for %d inputs",
			len(respBody.Predictions), len(batch)))
		return
	}

	for i := range respBody.Predictions {
		req := batch[i]
		req.ResCh <- PredictResult{
			Pred: respBody.Predictions[i],
		}
		if req.Seq == req.Size-1 {
			close(req.ResCh)
		}
	}
}

func (b *Batcher) fail(batch []PredictMessage, err error) {
	for i, r := range batch {
		r.ResCh <- PredictResult{
			Pred: "",
			Err:  err,
		}
		if i == len(batch)-1 {
			close(r.ResCh)
		}
	}
}
