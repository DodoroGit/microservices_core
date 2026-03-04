import os


class Config:
    MONGODB_URL: str = os.getenv("MONGODB_URL", "mongodb://admin:admin123@mongodb:27017")
    MONGODB_DB: str = os.getenv("MONGODB_DB", "notedb")


cfg = Config()
