import logging

from fastapi import FastAPI

from routers import chat

logging.basicConfig(level=logging.INFO)

app = FastAPI(title="AI Service")
app.include_router(chat.router)


@app.get("/health")
def health():
    from config import cfg
    return {"status": "healthy", "service": "ai-service", "model": cfg.MODEL_NAME}
