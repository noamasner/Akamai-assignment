from fastapi import FastAPI
from pydantic import BaseModel
from typing import List
import uvicorn
import time

# --- FastAPI App Initialization ---
app = FastAPI(
    title="Prediction Service",
    description="A simple API to simulate ML model predictions.",
    version="1.0.0",
)

# --- Pydantic Models for Request and Response ---
# These models provide automatic data validation and documentation.
class PredictionRequest(BaseModel):
    instances: List[str]

class PredictionResponse(BaseModel):
    predictions: List[str]

# --- Mock Model Function ---
# In a real scenario, this is where you would load your model
# and run inference on the input instances.
async def run_model(instances: List[str]) -> List[str]:
    """
    Simulates model inference.
    - Takes a list of strings (instances).
    - Returns a list of strings (predictions).
    """

    # Simulate a long-running model inference that takes time for a batch of instances
    time.sleep(0.5)
    print(f"--- Running model on batch of {len(instances)} instances ---")
    predictions = []
    for instance in instances:
        prediction = f"predicted_label_for_{instance.lower().replace(' ', '_')}"
        predictions.append(prediction)
    print(f"--- Finished model run ---")
    return predictions

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

