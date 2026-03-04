from unittest.mock import AsyncMock, MagicMock, patch

import pytest
import redis
from fastapi import HTTPException

from repositories.history_repo import HistoryRepository
from services.chat_service import ChatService


# ── Fake Repository（Unit Test 用）────────────────────────────────────────────

class FakeHistoryRepository:
    def __init__(self):
        self._store: dict[str, list[dict]] = {}

    def get(self, user_id: str) -> list[dict]:
        return self._store.get(user_id, [])

    def save(self, user_id: str, history: list[dict]) -> None:
        self._store[user_id] = history

    def delete(self, user_id: str) -> None:
        self._store.pop(user_id, None)


# ── Unit Test Fixtures ────────────────────────────────────────────────────────

@pytest.fixture
def repo() -> FakeHistoryRepository:
    return FakeHistoryRepository()

@pytest.fixture
def svc(repo: FakeHistoryRepository) -> ChatService:
    return ChatService(repo)


# ── Integration Test Fixtures ─────────────────────────────────────────────────

@pytest.fixture
def real_svc():
    client = redis.Redis(host="localhost", port=6379, decode_responses=True)
    repo = HistoryRepository(client)
    yield ChatService(repo)
    client.delete("chat:test-user")
    client.close()


# ── Helpers ───────────────────────────────────────────────────────────────────

def mock_ollama(response_text: str):
    """mock httpx.AsyncClient，讓 Ollama 回傳固定文字"""
    mock_response = MagicMock()
    mock_response.json.return_value = {"message": {"content": response_text}}
    mock_response.raise_for_status = MagicMock()

    mock_instance = AsyncMock()
    mock_instance.post = AsyncMock(return_value=mock_response)
    mock_instance.__aenter__ = AsyncMock(return_value=mock_instance)
    mock_instance.__aexit__ = AsyncMock(return_value=None)

    return patch("services.chat_service.httpx.AsyncClient", return_value=mock_instance)


# ── Unit Tests: send_message ──────────────────────────────────────────────────

@pytest.mark.asyncio
async def test_send_message_returns_response(svc: ChatService):
    with mock_ollama("你好！"):
        result = await svc.send_message("user-1", "你好")

    assert result["response"] == "你好！"

@pytest.mark.asyncio
async def test_send_message_appends_to_history(svc: ChatService, repo: FakeHistoryRepository):
    with mock_ollama("回覆"):
        await svc.send_message("user-1", "訊息")

    history = repo.get("user-1")
    assert len(history) == 2
    assert history[0] == {"role": "user", "content": "訊息"}
    assert history[1] == {"role": "assistant", "content": "回覆"}

@pytest.mark.asyncio
async def test_send_message_raises_504_on_timeout(svc: ChatService):
    import httpx
    mock_instance = AsyncMock()
    mock_instance.post = AsyncMock(side_effect=httpx.TimeoutException("timeout"))
    mock_instance.__aenter__ = AsyncMock(return_value=mock_instance)
    mock_instance.__aexit__ = AsyncMock(return_value=None)

    with patch("services.chat_service.httpx.AsyncClient", return_value=mock_instance):
        with pytest.raises(HTTPException) as exc:
            await svc.send_message("user-1", "你好")

    assert exc.value.status_code == 504

@pytest.mark.asyncio
async def test_send_message_raises_503_on_error(svc: ChatService):
    mock_instance = AsyncMock()
    mock_instance.post = AsyncMock(side_effect=Exception("connection error"))
    mock_instance.__aenter__ = AsyncMock(return_value=mock_instance)
    mock_instance.__aexit__ = AsyncMock(return_value=None)

    with patch("services.chat_service.httpx.AsyncClient", return_value=mock_instance):
        with pytest.raises(HTTPException) as exc:
            await svc.send_message("user-1", "你好")

    assert exc.value.status_code == 503


# ── Unit Tests: fetch_history / delete_history ────────────────────────────────

def test_fetch_history_returns_empty_initially(svc: ChatService):
    result = svc.fetch_history("user-1")

    assert result == []

@pytest.mark.asyncio
async def test_fetch_history_returns_saved_history(svc: ChatService):
    with mock_ollama("回覆"):
        await svc.send_message("user-1", "訊息")

    result = svc.fetch_history("user-1")

    assert len(result) == 2

def test_delete_history_clears_history(svc: ChatService, repo: FakeHistoryRepository):
    repo.save("user-1", [{"role": "user", "content": "訊息"}])

    svc.delete_history("user-1")

    assert svc.fetch_history("user-1") == []


# ── Integration Tests（real Redis + mock Ollama）──────────────────────────────

@pytest.mark.integration
@pytest.mark.asyncio
async def test_integration_send_message_saves_to_redis(real_svc: ChatService):
    with mock_ollama("整合測試回覆"):
        result = await real_svc.send_message("test-user", "整合測試訊息")

    history = real_svc.fetch_history("test-user")

    assert result["response"] == "整合測試回覆"
    assert len(history) == 2

@pytest.mark.integration
def test_integration_delete_history_clears_redis(real_svc: ChatService):
    real_svc._repo.save("test-user", [{"role": "user", "content": "訊息"}])

    real_svc.delete_history("test-user")

    assert real_svc.fetch_history("test-user") == []
