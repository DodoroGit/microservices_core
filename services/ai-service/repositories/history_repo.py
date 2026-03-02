import json

import redis

from config import cfg

_client = redis.Redis(host=cfg.REDIS_HOST, port=cfg.REDIS_PORT, decode_responses=True)


def _key(user_id: str) -> str:
    return f"chat:{user_id}"


def get(user_id: str) -> list[dict]:
    raw = _client.get(_key(user_id))
    return json.loads(raw) if raw else []


def save(user_id: str, history: list[dict]) -> None:
    _client.setex(_key(user_id), cfg.CHAT_HISTORY_TTL, json.dumps(history, ensure_ascii=False))


def delete(user_id: str) -> None:
    _client.delete(_key(user_id))
