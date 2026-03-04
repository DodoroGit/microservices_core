import pytest
from fastapi import HTTPException
from motor.motor_asyncio import AsyncIOMotorClient

from models.schemas import NoteCreate, NoteUpdate
from repositories.note_repository import NoteRepository
from services.note_service import NoteService


# ── Fake Repository（Unit Test 用）────────────────────────────────────────────

class FakeNoteRepository:
    def __init__(self):
        self._store: dict[str, dict] = {}
        self._counter = 0

    def _next_id(self) -> str:
        self._counter += 1
        return str(self._counter)

    async def find_all(self, category: str | None = None) -> list[dict]:
        notes = list(self._store.values())
        if category:
            notes = [n for n in notes if n["category"] == category]
        return notes

    async def find_by_id(self, note_id: str) -> dict | None:
        return self._store.get(note_id)

    async def create(self, data: dict) -> dict:
        note_id = self._next_id()
        note = {"_id": note_id, **data}
        self._store[note_id] = note
        return note

    async def update(self, note_id: str, data: dict) -> dict | None:
        if note_id not in self._store:
            return None
        self._store[note_id].update(data)
        return self._store[note_id]

    async def delete(self, note_id: str) -> bool:
        if note_id not in self._store:
            return False
        del self._store[note_id]
        return True


# ── Unit Test Fixtures ────────────────────────────────────────────────────────

@pytest.fixture
def repo() -> FakeNoteRepository:
    return FakeNoteRepository()

@pytest.fixture
def svc(repo: FakeNoteRepository) -> NoteService:
    return NoteService(repo)


# ── Integration Test Fixtures ─────────────────────────────────────────────────

@pytest.fixture
async def real_svc():  # type: ignore[return]
    client = AsyncIOMotorClient("mongodb://admin:admin123@localhost:27017")
    col = client["test_notedb"]["notes"]
    yield NoteService(NoteRepository(col))
    await col.drop()
    client.close()


# ── Unit Tests: create ────────────────────────────────────────────────────────

@pytest.mark.asyncio
async def test_create_returns_note(svc: NoteService):
    note = await svc.create(NoteCreate(title="標題", category="專案", content="內容"))

    assert note["title"] == "標題"
    assert note["category"] == "專案"
    assert note["content"] == "內容"


# ── Unit Tests: get_all ───────────────────────────────────────────────────────

@pytest.mark.asyncio
async def test_get_all_returns_all_notes(svc: NoteService):
    await svc.create(NoteCreate(title="A", category="專案", content=""))
    await svc.create(NoteCreate(title="B", category="技術筆記", content=""))

    notes = await svc.get_all()
    assert len(notes) == 2

@pytest.mark.asyncio
async def test_get_all_filters_by_category(svc: NoteService):
    await svc.create(NoteCreate(title="A", category="專案", content=""))
    await svc.create(NoteCreate(title="B", category="技術筆記", content=""))

    notes = await svc.get_all(category="專案")
    assert len(notes) == 1
    assert notes[0]["title"] == "A"


# ── Unit Tests: get_by_id ─────────────────────────────────────────────────────

@pytest.mark.asyncio
async def test_get_by_id_returns_note(svc: NoteService):
    created = await svc.create(NoteCreate(title="標題", category="專案", content="內容"))
    note = await svc.get_by_id(created["_id"])

    assert note["title"] == "標題"

@pytest.mark.asyncio
async def test_get_by_id_raises_404_when_not_found(svc: NoteService):
    with pytest.raises(HTTPException) as exc:
        await svc.get_by_id("nonexistent-id")

    assert exc.value.status_code == 404


# ── Unit Tests: update ────────────────────────────────────────────────────────

@pytest.mark.asyncio
async def test_update_modifies_note(svc: NoteService):
    created = await svc.create(NoteCreate(title="舊標題", category="專案", content="舊內容"))
    updated = await svc.update(created["_id"], NoteUpdate(title="新標題"))

    assert updated["title"] == "新標題"
    assert updated["content"] == "舊內容"

@pytest.mark.asyncio
async def test_update_raises_404_when_not_found(svc: NoteService):
    with pytest.raises(HTTPException) as exc:
        await svc.update("nonexistent-id", NoteUpdate(title="新標題"))

    assert exc.value.status_code == 404


# ── Unit Tests: delete ────────────────────────────────────────────────────────

@pytest.mark.asyncio
async def test_delete_removes_note(svc: NoteService, repo: FakeNoteRepository):
    created = await svc.create(NoteCreate(title="標題", category="專案", content="內容"))
    await svc.delete(created["_id"])

    assert await repo.find_by_id(created["_id"]) is None

@pytest.mark.asyncio
async def test_delete_raises_404_when_not_found(svc: NoteService):
    with pytest.raises(HTTPException) as exc:
        await svc.delete("nonexistent-id")

    assert exc.value.status_code == 404


# ── Integration Tests ─────────────────────────────────────────────────────────

@pytest.mark.integration
@pytest.mark.asyncio
async def test_integration_create_and_get_by_id(real_svc: NoteService):
    created = await real_svc.create(NoteCreate(title="標題", category="專案", content="內容"))

    fetched = await real_svc.get_by_id(created["_id"])

    assert fetched["title"] == "標題"

@pytest.mark.integration
@pytest.mark.asyncio
async def test_integration_get_all_and_filter(real_svc: NoteService):
    await real_svc.create(NoteCreate(title="A", category="專案", content=""))
    await real_svc.create(NoteCreate(title="B", category="技術筆記", content=""))

    all_notes = await real_svc.get_all()
    filtered = await real_svc.get_all(category="專案")

    assert len(all_notes) == 2
    assert len(filtered) == 1
    assert filtered[0]["title"] == "A"

@pytest.mark.integration
@pytest.mark.asyncio
async def test_integration_update_note(real_svc: NoteService):
    created = await real_svc.create(NoteCreate(title="舊標題", category="專案", content="內容"))

    updated = await real_svc.update(created["_id"], NoteUpdate(title="新標題"))

    assert updated["title"] == "新標題"
    assert updated["content"] == "內容"

@pytest.mark.integration
@pytest.mark.asyncio
async def test_integration_delete_note(real_svc: NoteService):
    created = await real_svc.create(NoteCreate(title="標題", category="專案", content="內容"))

    await real_svc.delete(created["_id"])

    with pytest.raises(HTTPException) as exc:
        await real_svc.get_by_id(created["_id"])
    assert exc.value.status_code == 404

@pytest.mark.integration
@pytest.mark.asyncio
async def test_integration_get_by_id_raises_404(real_svc: NoteService):
    with pytest.raises(HTTPException) as exc:
        await real_svc.get_by_id("000000000000000000000000")

    assert exc.value.status_code == 404
