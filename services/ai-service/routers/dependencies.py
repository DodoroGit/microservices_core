from core.cache import redis_client
from repositories.history_repo import HistoryRepository
from services.chat_service import ChatService, ChatServiceProtocol


def get_chat_service() -> ChatServiceProtocol:
    repo = HistoryRepository(redis_client)
    return ChatService(repo)
