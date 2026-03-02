import logging

import httpx
from fastapi import HTTPException

from config import cfg
from repositories import history_repo

logger = logging.getLogger(__name__)


async def chat(user_id: str, message: str) -> dict:
    history = history_repo.get(user_id)
    history.append({"role": "user", "content": message})

    async with httpx.AsyncClient(timeout=120.0) as client:
        try:
            response = await client.post(
                f"{cfg.OLLAMA_URL}/api/chat",
                json={"model": cfg.MODEL_NAME, "messages": history, "stream": False},
            )
            response.raise_for_status()
        except httpx.TimeoutException:
            raise HTTPException(status_code=504, detail="模型回應超時，請再試一次")
        except Exception as e:
            logger.error(f"Ollama error: {e}")
            raise HTTPException(status_code=503, detail="AI 服務暫時無法使用")

    assistant_message = response.json()["message"]["content"]
    history.append({"role": "assistant", "content": assistant_message})

    history_repo.save(user_id, history)

    return {"response": assistant_message, "history": history}


def get_history(user_id: str) -> list[dict]:
    return history_repo.get(user_id)


def clear_history(user_id: str) -> None:
    history_repo.delete(user_id)
