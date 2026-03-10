import pytest
from fastapi.testclient import TestClient
from fastapi import FastAPI

from models.schemas import ProjectStoryResponse
from routers.ai import router
from routers.dependencies import get_ai_service
from services.ai_service import AIServiceProtocol


# -------------------------------------------------------------------
# FakeAIService：模擬 AIService，不發真實請求
# -------------------------------------------------------------------

class FakeAIService:
    async def generate_project_story(self, project_title: str) -> ProjectStoryResponse:
        return ProjectStoryResponse(
            project_title=project_title,
            resume_description="建立了微服務架構系統，採用 Go 與 Python 雙語言實作。",
            tech_keywords=["Go", "Python", "MongoDB", "Docker"],
            interview_story="面對多語言微服務的整合挑戰，設計 API Gateway 統一對外入口。",
            source_notes_count=3,
        )


class FakeAIServiceNotFound:
    async def generate_project_story(self, project_title: str) -> ProjectStoryResponse:
        raise ValueError(f"找不到與「{project_title}」相關的專案筆記")


def make_test_client(svc: AIServiceProtocol) -> TestClient:
    app = FastAPI()
    app.include_router(router)
    app.dependency_overrides[get_ai_service] = lambda: svc
    return TestClient(app)


# ===================================================================
# POST /ai/project-story 測試
# ===================================================================

class TestProjectStoryEndpoint:
    def test_success(self):
        client = make_test_client(FakeAIService())

        resp = client.post("/ai/project-story", json={"project_title": "microservices_core"})

        assert resp.status_code == 200
        data = resp.json()
        assert data["project_title"] == "microservices_core"
        assert data["resume_description"] != ""
        assert isinstance(data["tech_keywords"], list)
        assert len(data["tech_keywords"]) > 0
        assert data["interview_story"] != ""
        assert data["source_notes_count"] == 3

    def test_not_found_returns_404(self):
        client = make_test_client(FakeAIServiceNotFound())

        resp = client.post("/ai/project-story", json={"project_title": "不存在的專案"})

        assert resp.status_code == 404
        assert "找不到" in resp.json()["detail"]

    def test_missing_project_title_returns_422(self):
        # project_title 是必填欄位，缺少時 FastAPI 自動回 422
        client = make_test_client(FakeAIService())

        resp = client.post("/ai/project-story", json={})

        assert resp.status_code == 422

    def test_empty_project_title(self):
        # 空字串通過 Pydantic 驗證，但 service 層會因找不到筆記而拋 ValueError
        client = make_test_client(FakeAIServiceNotFound())

        resp = client.post("/ai/project-story", json={"project_title": ""})

        assert resp.status_code == 404
