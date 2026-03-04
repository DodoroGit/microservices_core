from datetime import datetime, timezone
from typing import Protocol

from bson import ObjectId
from motor.motor_asyncio import AsyncIOMotorCollection


class NoteRepositoryProtocol(Protocol):
    async def find_all(self, category: str | None = None) -> list[dict]: ...
    async def find_by_id(self, note_id: str) -> dict | None: ...
    async def create(self, data: dict) -> dict: ...
    async def update(self, note_id: str, data: dict) -> dict | None: ...
    async def delete(self, note_id: str) -> bool: ...


class NoteRepository:
    def __init__(self, collection: AsyncIOMotorCollection):
        self._col = collection

    async def find_all(self, category: str | None = None) -> list[dict]:
        query = {}
        if category:
            query["category"] = category
        cursor = self._col.find(query).sort("created_at", -1)
        return await cursor.to_list(length=None)

    async def find_by_id(self, note_id: str) -> dict | None:
        try:
            return await self._col.find_one({"_id": ObjectId(note_id)})
        except Exception:
            return None

    async def create(self, data: dict) -> dict:
        now = datetime.now(timezone.utc)
        data["created_at"] = now
        data["updated_at"] = now
        result = await self._col.insert_one(data)
        return await self.find_by_id(str(result.inserted_id))

    async def update(self, note_id: str, data: dict) -> dict | None:
        data["updated_at"] = datetime.now(timezone.utc)
        await self._col.update_one(
            {"_id": ObjectId(note_id)},
            {"$set": data}
        )
        return await self.find_by_id(note_id)

    async def delete(self, note_id: str) -> bool:
        result = await self._col.delete_one({"_id": ObjectId(note_id)})
        return result.deleted_count == 1
