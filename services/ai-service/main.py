import logging

from fastapi import FastAPI

from routers import ai_router, health_router

logging.basicConfig(level=logging.INFO)

app = FastAPI(title="AI Service")
app.include_router(health_router.router)
app.include_router(ai_router.router)
