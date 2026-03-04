from fastapi import HTTPException

from models.schemas import NoteCreate, NoteUpdate
from repositories.note_repository import NoteRepository


class NoteService:
    def __init__(self, repo: NoteRepository):
        self._repo = repo

    async def get_all(self, category: str | None = None) -> list[dict]:
        return await self._repo.find_all(category)

    async def get_by_id(self, note_id: str) -> dict:
        note = await self._repo.find_by_id(note_id)
        if not note:
            raise HTTPException(status_code=404, detail="筆記不存在")
        return note

    async def create(self, data: NoteCreate) -> dict:
        return await self._repo.create(data.model_dump())

    async def update(self, note_id: str, data: NoteUpdate) -> dict:
        await self.get_by_id(note_id)
        update_data = {k: v for k, v in data.model_dump().items() if v is not None}
        return await self._repo.update(note_id, update_data)

    async def delete(self, note_id: str) -> None:
        await self.get_by_id(note_id)
        await self._repo.delete(note_id)
