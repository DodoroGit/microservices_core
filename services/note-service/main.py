import logging

from fastapi import FastAPI

from routers import notes

logging.basicConfig(level=logging.INFO)

app = FastAPI(title="Note Service")
app.include_router(notes.router)


@app.get("/health")
def health():
    return {"status": "healthy", "service": "note-service"}
