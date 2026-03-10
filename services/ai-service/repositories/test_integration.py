"""
整合測試：需要真實的 note-service 與 Anthropic API Key

執行方式：
    pytest -v -m integration
    ANTHROPIC_API_KEY=xxx NOTE_SERVICE_URL=http://localhost:8083 pytest -v -m integration
"""
import os

import pytest

from repositories.claude_client import ClaudeClient
from repositories.note_client import NoteClient


def get_note_service_url() -> str:
    return os.getenv("NOTE_SERVICE_URL", "http://localhost:8083")


def get_api_key() -> str:
    return os.getenv("ANTHROPIC_API_KEY", "")


# ===================================================================
# NoteClient 整合測試（需要真實 note-service）
# ===================================================================

@pytest.mark.integration
class TestNoteClientIntegration:
    async def test_get_project_notes_returns_list(self):
        client = NoteClient(get_note_service_url())

        notes = await client.get_project_notes("test")

        # 不管有沒有筆記，應該回傳 list
        assert isinstance(notes, list)

    async def test_filter_by_project_title(self):
        client = NoteClient(get_note_service_url())

        notes = await client.get_project_notes("絕對不存在的專案名稱xyz123")

        assert notes == []


# ===================================================================
# ClaudeClient 整合測試（需要真實 Anthropic API Key）
# ===================================================================

@pytest.mark.integration
class TestClaudeClientIntegration:
    async def test_generate_project_story_returns_structured_output(self):
        api_key = get_api_key()
        if not api_key:
            pytest.skip("未設定 ANTHROPIC_API_KEY，跳過 Claude API 整合測試")

        client = ClaudeClient(api_key, "claude-haiku-4-5-20251001")  # 用 Haiku 降低測試成本
        notes = [
            {
                "title": "測試專案",
                "content": "使用 Go 建立 REST API，連接 PostgreSQL 資料庫",
                "category": "專案",
            }
        ]

        result = await client.generate_project_story("測試專案", notes)

        assert "resume_description" in result
        assert "tech_keywords" in result
        assert "interview_story" in result
        assert isinstance(result["tech_keywords"], list)
        assert len(result["resume_description"]) > 0
        assert len(result["interview_story"]) > 0
