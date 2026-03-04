from unittest.mock import AsyncMock, MagicMock

import pytest
from motor.motor_asyncio import AsyncIOMotorClient

from repositories.note_repository import NoteRepository


# ── Unit Test Fixtures ────────────────────────────────────────────────────────

@pytest.fixture
def mock_col() -> MagicMock:
    return MagicMock()

@pytest.fixture
def repo(mock_col: MagicMock) -> NoteRepository:
    return NoteRepository(mock_col)


# ── Integration Test Fixtures ─────────────────────────────────────────────────

@pytest.fixture
async def real_repo():  # type: ignore[return]
    client = AsyncIOMotorClient("mongodb://admin:admin123@localhost:27017")
    col = client["test_notedb"]["notes"]
    yield NoteRepository(col)
    await col.drop()
    client.close()


# ── Unit Tests: find_all ──────────────────────────────────────────────────────

@pytest.mark.asyncio
async def test_find_all_no_filter(repo: NoteRepository, mock_col: MagicMock):
    mock_cursor = MagicMock()
    mock_cursor.to_list = AsyncMock(return_value=[{"_id": "1", "title": "A"}])
    mock_col.find.return_value.sort.return_value = mock_cursor

    result = await repo.find_all()

    mock_col.find.assert_called_once_with({})
    assert len(result) == 1

@pytest.mark.asyncio
async def test_find_all_with_category_filter(repo: NoteRepository, mock_col: MagicMock):
    mock_cursor = MagicMock()
    mock_cursor.to_list = AsyncMock(return_value=[])
    mock_col.find.return_value.sort.return_value = mock_cursor

    await repo.find_all(category="專案")

    mock_col.find.assert_called_once_with({"category": "專案"})


# ── Unit Tests: find_by_id ────────────────────────────────────────────────────

@pytest.mark.asyncio
async def test_find_by_id_returns_note(repo: NoteRepository, mock_col: MagicMock):
    mock_col.find_one = AsyncMock(return_value={"_id": "abc", "title": "標題"})

    result = await repo.find_by_id("abc")

    assert result["title"] == "標題"

@pytest.mark.asyncio
async def test_find_by_id_returns_none_on_invalid_id(repo: NoteRepository):
    result = await repo.find_by_id("not-a-valid-objectid")

    assert result is None


# ── Unit Tests: create ────────────────────────────────────────────────────────

@pytest.mark.asyncio
async def test_create_inserts_and_returns_note(repo: NoteRepository, mock_col: MagicMock):
    inserted_id = "new-id-123"
    mock_col.insert_one = AsyncMock(return_value=MagicMock(inserted_id=inserted_id))
    mock_col.find_one = AsyncMock(return_value={"_id": inserted_id, "title": "新筆記"})

    result = await repo.create({"title": "新筆記", "category": "專案", "content": ""})

    assert mock_col.insert_one.called
    assert result["title"] == "新筆記"


# ── Unit Tests: update ────────────────────────────────────────────────────────

@pytest.mark.asyncio
async def test_update_calls_update_one(repo: NoteRepository, mock_col: MagicMock):
    mock_col.update_one = AsyncMock()
    mock_col.find_one = AsyncMock(return_value={"_id": "abc", "title": "新標題"})

    result = await repo.update("abc", {"title": "新標題"})

    assert mock_col.update_one.called
    assert result["title"] == "新標題"


# ── Unit Tests: delete ────────────────────────────────────────────────────────

@pytest.mark.asyncio
async def test_delete_returns_true_on_success(repo: NoteRepository, mock_col: MagicMock):
    mock_col.delete_one = AsyncMock(return_value=MagicMock(deleted_count=1))

    result = await repo.delete("abc")

    assert result is True

@pytest.mark.asyncio
async def test_delete_returns_false_when_not_found(repo: NoteRepository, mock_col: MagicMock):
    mock_col.delete_one = AsyncMock(return_value=MagicMock(deleted_count=0))

    result = await repo.delete("nonexistent")

    assert result is False


# ── Integration Tests ─────────────────────────────────────────────────────────

@pytest.mark.integration
@pytest.mark.asyncio
async def test_integration_create_and_find_by_id(real_repo: NoteRepository):
    created = await real_repo.create({"title": "標題", "category": "專案", "content": "內容"})

    fetched = await real_repo.find_by_id(str(created["_id"]))

    assert fetched["title"] == "標題"
    assert fetched["category"] == "專案"

@pytest.mark.integration
@pytest.mark.asyncio
async def test_integration_find_all_returns_all(real_repo: NoteRepository):
    await real_repo.create({"title": "A", "category": "專案", "content": ""})
    await real_repo.create({"title": "B", "category": "技術筆記", "content": ""})

    result = await real_repo.find_all()

    assert len(result) == 2

@pytest.mark.integration
@pytest.mark.asyncio
async def test_integration_find_all_filters_by_category(real_repo: NoteRepository):
    await real_repo.create({"title": "A", "category": "專案", "content": ""})
    await real_repo.create({"title": "B", "category": "技術筆記", "content": ""})

    result = await real_repo.find_all(category="專案")

    assert len(result) == 1
    assert result[0]["title"] == "A"

@pytest.mark.integration
@pytest.mark.asyncio
async def test_integration_update_modifies_document(real_repo: NoteRepository):
    created = await real_repo.create({"title": "舊標題", "category": "專案", "content": "內容"})

    updated = await real_repo.update(str(created["_id"]), {"title": "新標題"})

    assert updated["title"] == "新標題"
    assert updated["content"] == "內容"

@pytest.mark.integration
@pytest.mark.asyncio
async def test_integration_delete_removes_document(real_repo: NoteRepository):
    created = await real_repo.create({"title": "標題", "category": "專案", "content": "內容"})

    result = await real_repo.delete(str(created["_id"]))
    fetched = await real_repo.find_by_id(str(created["_id"]))

    assert result is True
    assert fetched is None

@pytest.mark.integration
@pytest.mark.asyncio
async def test_integration_find_by_id_returns_none_when_not_found(real_repo: NoteRepository):
    result = await real_repo.find_by_id("000000000000000000000000")

    assert result is None
