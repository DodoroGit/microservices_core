import logging

from fastapi import FastAPI

from routers import ai

logging.basicConfig(level=logging.INFO)

app = FastAPI(title="AI Service")
app.include_router(ai.router)


@app.get("/health")
def health():
    from config import settings
    return {
        "status": "healthy",
        "service": "ai-service",
        "model": settings.claude_model,
    }
