## Python Prediction Service

A tiny FastAPI service that simulates ML predictions for a batch of input strings.

### Requirements
- Python 3.12+
- `uv` package manager

### Quickstart
1. Create a venv using UV (use python 3.12)
2. Install dependencies using UV (see pyproject.toml file)
3. Run the server. (Either run main.py or run: `uvicorn main:app --host 0.0.0.0 --port 8080`)

The API will be available at http://localhost:8080 and OpenAPI docs at:
- Swagger UI: http://localhost:8080/docs
- ReDoc: http://localhost:8080/redoc

### API
- POST /api/v1/predict
  - Request:
    ```json
    {
      "instances": ["Hello World", "test"]
    }
    ```
  - Response:
    ```json
    {
      "predictions": [
        "predicted_label_for_hello_world",
        "predicted_label_for_test"
      ]
    }
    ```

Example:
```bash
curl -X POST http://localhost:8080/api/v1/predict \
  -H "Content-Type: application/json" \
  -d '{"instances": ["Hello World", "Another sample"]}'
```

### Notes
- Project name: `python_detection_service` (see `pyproject.toml`).
- Main entrypoint: `main.py`, FastAPI app instance is `app`.