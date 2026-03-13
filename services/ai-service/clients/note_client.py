from typing import Protocol

import httpx


class NoteClientProtocol(Protocol):
    async def get_notes_by_category(self, category: str) -> list[dict]: ...


class NoteClient:
    """透過 HTTP 向 note-service 取得筆記資料"""

    def __init__(self, base_url: str):
        self._base_url = base_url

    async def get_notes_by_category(self, category: str) -> list[dict]:
        """取得指定類別的所有筆記"""
        async with httpx.AsyncClient(timeout=10.0) as client:
            resp = await client.get(
                f"{self._base_url}/notes",
                params={"category": category},
            )
            resp.raise_for_status()
            return resp.json()

