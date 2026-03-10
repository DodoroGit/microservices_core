from typing import Protocol

import httpx


class NoteClientProtocol(Protocol):
    async def get_notes_by_category(self, category: str) -> list[dict]: ...
    async def get_project_notes(self, project_title: str) -> list[dict]: ...


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

    async def get_project_notes(self, project_title: str) -> list[dict]:
        """
        取得所有「專案」類別的筆記，並依 project_title 過濾。
        過濾邏輯：note 的 title 或 content 中含有 project_title（不分大小寫）。
        """
        notes = await self.get_notes_by_category("專案")
        keyword = project_title.lower()
        return [
            n for n in notes
            if keyword in n.get("title", "").lower()
            or keyword in n.get("content", "").lower()
        ]
