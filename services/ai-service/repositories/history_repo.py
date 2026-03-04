import json

import redis

from config import cfg


class HistoryRepository:
    def __init__(self, client: redis.Redis):
        self._client = client

    def _key(self, user_id: str) -> str:
        return f"chat:{user_id}"

    def get(self, user_id: str) -> list[dict]:
        raw = self._client.get(self._key(user_id))
        return json.loads(raw) if raw else []

    def save(self, user_id: str, history: list[dict]) -> None:
        self._client.setex(self._key(user_id), cfg.CHAT_HISTORY_TTL, json.dumps(history, ensure_ascii=False))

    def delete(self, user_id: str) -> None:
        self._client.delete(self._key(user_id))
