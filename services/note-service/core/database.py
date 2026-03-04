from motor.motor_asyncio import AsyncIOMotorClient

from config import cfg

_client = AsyncIOMotorClient(cfg.MONGODB_URL)
db = _client[cfg.MONGODB_DB]
