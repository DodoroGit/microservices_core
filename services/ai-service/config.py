import os


class Config:
    OLLAMA_URL: str = os.getenv("OLLAMA_URL", "http://ollama:11434")
    MODEL_NAME: str = os.getenv("MODEL_NAME", "llama3.2:3b")
    REDIS_HOST: str = os.getenv("REDIS_HOST", "redis")
    REDIS_PORT: int = int(os.getenv("REDIS_PORT", "6379"))
    CHAT_HISTORY_TTL: int = 60 * 60 * 24  # 24 小時


cfg = Config()
