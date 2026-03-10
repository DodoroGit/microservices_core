from config import settings
from repositories.note_client import NoteClient
from repositories.claude_client import ClaudeClient
from services.ai_service import AIService, AIServiceProtocol


def get_ai_service() -> AIServiceProtocol:
    note_client = NoteClient(settings.note_service_url)
    claude_client = ClaudeClient(settings.anthropic_api_key, settings.claude_model)
    return AIService(note_client, claude_client)
