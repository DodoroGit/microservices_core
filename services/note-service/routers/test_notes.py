from datetime import datetime, timezone

import pytest
from fastapi.testclient import TestClient

from main import app
from models.schemas import NoteCreate, NoteUpdate
from routers.dependencies import get_note_service

# ── Fake Service ──────────────────────────────────────────────────────────────
# 實作 NoteServiceProtocol，回傳固定資料，不碰任何 DB
# 等同於 Go 裡 dependency_overrides 用的 mock handler

NOW = datetime(2024, 1, 1, tzinfo=timezone.utc)

FAKE_NOTE = {
    "_id": "note-1",
    "title": "測試標題",
    "category": "專案",
    "content": "測試內容",
    "created_at": NOW,
    "updated_at": NOW,
}


class FakeNoteService:
    async def get_all(self, category: str | None = None) -> list[dict]:
        if category and category != FAKE_NOTE["category"]:
            return []
        return [FAKE_NOTE]

    async def get_by_id(self, note_id: str) -> dict:
        from fastapi import HTTPException
        if note_id != FAKE_NOTE["_id"]:
            raise HTTPException(status_code=404, detail="筆記不存在")
        return FAKE_NOTE

    async def create(self, data: NoteCreate) -> dict:
        return {**FAKE_NOTE, "title": data.title, "category": data.category, "content": data.content}

    async def update(self, note_id: str, data: NoteUpdate) -> dict:
        from fastapi import HTTPException
        if note_id != FAKE_NOTE["_id"]:
            raise HTTPException(status_code=404, detail="筆記不存在")
        return {**FAKE_NOTE, **(data.model_dump(exclude_none=True))}

    async def delete(self, note_id: str) -> None:
        from fastapi import HTTPException
        if note_id != FAKE_NOTE["_id"]:
            raise HTTPException(status_code=404, detail="筆記不存在")


# ── Fixtures ──────────────────────────────────────────────────────────────────

@pytest.fixture
def client():
    app.dependency_overrides[get_note_service] = lambda: FakeNoteService()
    yield TestClient(app)
    app.dependency_overrides.clear()


# ── Tests: GET /notes ─────────────────────────────────────────────────────────

def test_list_notes_returns_200(client: TestClient):
    res = client.get("/notes")

    assert res.status_code == 200
    assert len(res.json()) == 1

def test_list_notes_filter_by_category(client: TestClient):
    res = client.get("/notes", params={"category": "技術筆記"})

    assert res.status_code == 200
    assert len(res.json()) == 0


# ── Tests: GET /notes/{id} ────────────────────────────────────────────────────

def test_get_note_returns_200(client: TestClient):
    res = client.get("/notes/note-1")

    assert res.status_code == 200
    assert res.json()["title"] == "測試標題"

def test_get_note_returns_404_when_not_found(client: TestClient):
    res = client.get("/notes/nonexistent")

    assert res.status_code == 404


# ── Tests: POST /notes ────────────────────────────────────────────────────────

def test_create_note_returns_201(client: TestClient):
    payload = {"title": "新標題", "category": "專案", "content": "新內容"}
    res = client.post("/notes", json=payload)

    assert res.status_code == 201
    assert res.json()["title"] == "新標題"

def test_create_note_returns_422_on_invalid_payload(client: TestClient):
    res = client.post("/notes", json={"title": "缺少 category 和 content"})

    assert res.status_code == 422


# ── Tests: PUT /notes/{id} ────────────────────────────────────────────────────

def test_update_note_returns_200(client: TestClient):
    payload = {"title": "更新後標題"}
    res = client.put("/notes/note-1", json=payload)

    assert res.status_code == 200
    assert res.json()["title"] == "更新後標題"

def test_update_note_returns_404_when_not_found(client: TestClient):
    res = client.put("/notes/nonexistent", json={"title": "X"})

    assert res.status_code == 404


# ── Tests: DELETE /notes/{id} ─────────────────────────────────────────────────

def test_delete_note_returns_204(client: TestClient):
    res = client.delete("/notes/note-1")

    assert res.status_code == 204

def test_delete_note_returns_404_when_not_found(client: TestClient):
    res = client.delete("/notes/nonexistent")

    assert res.status_code == 404
