import json
from unittest.mock import MagicMock

import pytest
import redis

from repositories.history_repo import HistoryRepository


# ── Unit Test Fixtures ────────────────────────────────────────────────────────

@pytest.fixture
def mock_client() -> MagicMock:
    return MagicMock(spec=redis.Redis)

@pytest.fixture
def repo(mock_client: MagicMock) -> HistoryRepository:
    return HistoryRepository(mock_client)


# ── Integration Test Fixtures ─────────────────────────────────────────────────

@pytest.fixture
def real_repo():
    client = redis.Redis(host="localhost", port=6379, decode_responses=True)
    repo = HistoryRepository(client)
    yield repo
    client.delete("chat:test-user")
    client.close()


# ── Unit Tests: get ───────────────────────────────────────────────────────────

def test_get_returns_empty_list_when_no_history(repo: HistoryRepository, mock_client: MagicMock):
    mock_client.get.return_value = None

    result = repo.get("user-1")

    assert result == []

def test_get_returns_parsed_history(repo: HistoryRepository, mock_client: MagicMock):
    history = [{"role": "user", "content": "你好"}]
    mock_client.get.return_value = json.dumps(history)

    result = repo.get("user-1")

    assert result == history


# ── Unit Tests: save ──────────────────────────────────────────────────────────

def test_save_calls_setex(repo: HistoryRepository, mock_client: MagicMock):
    history = [{"role": "user", "content": "你好"}]

    repo.save("user-1", history)

    assert mock_client.setex.called
    args = mock_client.setex.call_args[0]
    assert args[0] == "chat:user-1"
    assert json.loads(args[2]) == history


# ── Unit Tests: delete ────────────────────────────────────────────────────────

def test_delete_calls_redis_delete(repo: HistoryRepository, mock_client: MagicMock):
    repo.delete("user-1")

    mock_client.delete.assert_called_once_with("chat:user-1")


# ── Integration Tests ─────────────────────────────────────────────────────────

@pytest.mark.integration
def test_integration_save_and_get(real_repo: HistoryRepository):
    history = [{"role": "user", "content": "你好"}, {"role": "assistant", "content": "哈囉"}]

    real_repo.save("test-user", history)
    result = real_repo.get("test-user")

    assert result == history

@pytest.mark.integration
def test_integration_get_returns_empty_before_save(real_repo: HistoryRepository):
    result = real_repo.get("test-user")

    assert result == []

@pytest.mark.integration
def test_integration_delete_clears_history(real_repo: HistoryRepository):
    real_repo.save("test-user", [{"role": "user", "content": "你好"}])

    real_repo.delete("test-user")
    result = real_repo.get("test-user")

    assert result == []
