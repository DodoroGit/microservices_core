from motor.motor_asyncio import AsyncIOMotorClient

from config import cfg
from repositories.note_repository import NoteRepository
from services.note_service import NoteService

_client = AsyncIOMotorClient(cfg.MONGODB_URL)
_db = _client[cfg.MONGODB_DB]


def get_note_service() -> NoteService:
    repo = NoteRepository(_db["notes"])
    return NoteService(repo)
