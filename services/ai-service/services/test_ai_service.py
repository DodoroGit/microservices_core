import pytest

from models.schemas import ProjectStoryResponse
from services.ai_service import AIService


# -------------------------------------------------------------------
# Fake 實作：模擬 NoteClient 與 ClaudeClient，不發真實 HTTP/API 請求
# -------------------------------------------------------------------

class FakeNoteClient:
    def __init__(self, notes: list[dict]):
        self._notes = notes

    async def get_project_notes(self, project_title: str) -> list[dict]:
        return self._notes


class FakeClaudeClient:
    async def generate_project_story(
        self, project_title: str, notes: list[dict]
    ) -> dict:
        return {
            "resume_description": f"為 {project_title} 建立了微服務系統，採用 Go 與 Python 實作，解決了服務間通訊問題。",
            "tech_keywords": ["Go", "Python", "FastAPI", "MongoDB", "Docker"],
            "interview_story": f"在開發 {project_title} 時，面對多語言微服務的整合挑戰，設計了統一的 API Gateway 解決方案。",
        }


SAMPLE_NOTES = [
    {
        "title": "microservices_core 架構設計",
        "content": "使用 Go 實作 user-service，Python 實作 note-service 與 ai-service",
        "category": "專案",
    },
    {
        "title": "microservices_core API Gateway",
        "content": "使用 Gin 實作 reverse proxy，支援 JWT 認證",
        "category": "專案",
    },
]


# ===================================================================
# generate_project_story 測試
# ===================================================================

class TestGenerateProjectStory:
    async def test_success(self):
        svc = AIService(FakeNoteClient(SAMPLE_NOTES), FakeClaudeClient())

        result = await svc.generate_project_story("microservices_core")

        assert isinstance(result, ProjectStoryResponse)
        assert result.project_title == "microservices_core"
        assert result.resume_description != ""
        assert len(result.tech_keywords) > 0
        assert result.interview_story != ""
        # 使用的筆記數量要與 FakeNoteClient 回傳的一致
        assert result.source_notes_count == len(SAMPLE_NOTES)

    async def test_no_notes_found_raises_value_error(self):
        # FakeNoteClient 回傳空清單，應該拋出 ValueError
        svc = AIService(FakeNoteClient([]), FakeClaudeClient())

        with pytest.raises(ValueError, match="找不到與"):
            await svc.generate_project_story("不存在的專案")

    async def test_response_contains_project_title(self):
        svc = AIService(FakeNoteClient(SAMPLE_NOTES), FakeClaudeClient())

        result = await svc.generate_project_story("microservices_core")

        assert result.project_title == "microservices_core"

    async def test_source_notes_count_reflects_input(self):
        single_note = [SAMPLE_NOTES[0]]
        svc = AIService(FakeNoteClient(single_note), FakeClaudeClient())

        result = await svc.generate_project_story("microservices_core")

        assert result.source_notes_count == 1
