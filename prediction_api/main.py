import os

from fastapi import FastAPI
from pydantic import BaseModel
from typing import List
import uvicorn
import time

from sympy.printing.pytorch import torch
from transformers import pipeline

# ---------------- Load model ----------------
DEVICE = torch.device("cuda" if torch.cuda.is_available() else "cpu")
MODEL_DIR = os.environ.get("MODEL_DIR")
if MODEL_DIR is not None and os.path.isdir(MODEL_DIR):
    # The model exists locally. This should always happen when running the image.
    print(f"Loading model from {MODEL_DIR}")
    SENTIMENT_CLASSIFIER = pipeline("text-classification", model=MODEL_DIR, tokenizer=MODEL_DIR, device=DEVICE)
else:
    # The model does not exist locally. This will happen if the model wasn't pre-downloaded.
    # Downloading from huggingface if it's not in the cache folder.
    print("The model is not found locally. downloading from huggingface if not in cache")
    MODEL_NAME = "distilbert/distilbert-base-uncased-finetuned-sst-2-english"
    SENTIMENT_CLASSIFIER = pipeline("text-classification", model=MODEL_NAME, device=DEVICE)

# --- FastAPI App Initialization ---
app = FastAPI(
    title="Prediction Service",
    description="A simple API to simulate ML model predictions.",
    version="2.0.0",
)

# --- Pydantic Models for Request and Response ---
# These models provide automatic data validation and documentation.
class PredictionRequest(BaseModel):
    instances: List[str]

class PredictionResponse(BaseModel):
    predictions: List[str]

# --- Run Model Function ---
# In a real scenario, this is where you would load your model
# and run inference on the input instances.
async def run_model(instances: List[str]) -> List[str]:
    """
    Run model inference.
    - Takes a list of strings (instances).
    - Returns a list of strings (predictions).
    """
    predictions = SENTIMENT_CLASSIFIER(instances)
    return [str(prediction) for prediction in predictions]

# --- API Endpoint ---
@app.post("/api/v1/predict", response_model=PredictionResponse)
async def predict(request: PredictionRequest) -> PredictionResponse:
    """
    API endpoint to get predictions.
    Accepts a JSON payload with "instances" and returns "predictions".
    """
    
    print(time.strftime("--- Received request at %Y-%m-%d %H:%M:%S", time.gmtime()), f"{int((time.time() % 1) * 1000):03d} UTC")
    
    predictions = await run_model(request.instances)
    
    response = PredictionResponse(predictions=predictions)
    
    return response

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8080)

