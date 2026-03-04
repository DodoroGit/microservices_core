from core.database import db
from repositories.note_repository import NoteRepository
from services.note_service import NoteService


def get_note_service() -> NoteService:
    repo = NoteRepository(db["notes"])
    return NoteService(repo)
