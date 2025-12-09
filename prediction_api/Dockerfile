# Base image with uv
FROM ghcr.io/astral-sh/uv:0.9-python3.12-bookworm-slim

WORKDIR /app

# Copy project metadata
COPY pyproject.toml uv.lock ./

# Directory where the model will live
ENV MODEL_DIR=/opt/models/distilbert-sst2

# Install dependencies and download the model at build time using huggingface-cli
RUN mkdir -p "$MODEL_DIR" && \
    uv sync --frozen && \
    uv run hf download \
        distilbert/distilbert-base-uncased-finetuned-sst-2-english \
        --local-dir "$MODEL_DIR"

# Copy application code
COPY main.py .

# Expose the service port
EXPOSE 8080

# Run service
CMD ["uv", "run", "uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
