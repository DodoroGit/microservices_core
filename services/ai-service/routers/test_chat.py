import pytest
from fastapi.testclient import TestClient

from main import app
from routers.dependencies import get_chat_service


# ── Fake Service（Unit Test 用）───────────────────────────────────────────────

class FakeChatService:
    def __init__(self):
        self._history: list[dict] = []

    async def send_message(self, user_id: str, message: str) -> dict:
        self._history.append({"role": "user", "content": message})
        self._history.append({"role": "assistant", "content": "假回覆"})
        return {"response": "假回覆", "history": self._history}

    def fetch_history(self, user_id: str) -> list[dict]:
        return self._history

    def delete_history(self, user_id: str) -> None:
        self._history.clear()


# ── Fixtures ──────────────────────────────────────────────────────────────────

@pytest.fixture
def client():
    app.dependency_overrides[get_chat_service] = lambda: FakeChatService()
    yield TestClient(app)
    app.dependency_overrides.clear()


# ── Tests: POST /chat ─────────────────────────────────────────────────────────

def test_chat_returns_200(client: TestClient):
    res = client.post("/chat", json={"user_id": "user-1", "message": "你好"})

    assert res.status_code == 200
    assert res.json()["response"] == "假回覆"

def test_chat_returns_422_on_missing_fields(client: TestClient):
    res = client.post("/chat", json={"message": "缺少 user_id"})

    assert res.status_code == 422


# ── Tests: GET /chat/history ──────────────────────────────────────────────────

def test_get_history_returns_200(client: TestClient):
    res = client.get("/chat/history", params={"user_id": "user-1"})

    assert res.status_code == 200
    assert "history" in res.json()

def test_get_history_returns_422_on_missing_user_id(client: TestClient):
    res = client.get("/chat/history")

    assert res.status_code == 422


# ── Tests: DELETE /chat/history ───────────────────────────────────────────────

def test_delete_history_returns_200(client: TestClient):
    res = client.delete("/chat/history", params={"user_id": "user-1"})

    assert res.status_code == 200
    assert res.json()["message"] == "對話紀錄已清除"

def test_delete_history_returns_422_on_missing_user_id(client: TestClient):
    res = client.delete("/chat/history")

    assert res.status_code == 422
