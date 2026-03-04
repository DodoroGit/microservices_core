from core.cache import redis_client
from repositories.history_repo import HistoryRepository
from services.chat_service import ChatService


def get_chat_service() -> ChatService:
    repo = HistoryRepository(redis_client)
    return ChatService(repo)
