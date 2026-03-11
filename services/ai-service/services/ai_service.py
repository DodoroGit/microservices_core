from typing import Protocol

from models.schemas import (
    BackgroundResponse,
    DailyResponse,
    ProjectSummary,
    ProjectsOverviewResponse,
    SkillsResponse,
)
from repositories.claude_client import ClaudeClientProtocol
from repositories.note_client import NoteClientProtocol


class AIServiceProtocol(Protocol):
    async def generate_background(self) -> BackgroundResponse: ...
    async def generate_projects_overview(self) -> ProjectsOverviewResponse: ...
    async def generate_skills(self) -> SkillsResponse: ...
    async def generate_daily(self) -> DailyResponse: ...


class AIService:
    def __init__(
        self,
        note_client: NoteClientProtocol,
        claude_client: ClaudeClientProtocol,
    ):
        self._note_client = note_client
        self._claude_client = claude_client

    # ── 個人背景介紹 ────────────────────────────────────────────────────────
    async def generate_background(self) -> BackgroundResponse:
        notes = await self._note_client.get_notes_by_category("履歷")
        if not notes:
            raise ValueError("找不到「履歷」類別的筆記，請先在筆記頁面新增履歷相關內容")
        result = await self._claude_client.generate_background(notes)
        return BackgroundResponse(
            background_intro=result["background_intro"],
            key_roles=result["key_roles"],
            years_summary=result["years_summary"],
            source_notes_count=len(notes),
        )

    # ── 專案介紹 ────────────────────────────────────────────────────────────
    async def generate_projects_overview(self) -> ProjectsOverviewResponse:
        notes = await self._note_client.get_notes_by_category("專案")
        if not notes:
            raise ValueError("找不到「專案」類別的筆記，請先在筆記頁面新增專案相關內容")
        result = await self._claude_client.generate_projects_overview(notes)
        projects = [
            ProjectSummary(
                project_name=p["project_name"],
                description=p["description"],
                tech_keywords=p["tech_keywords"],
            )
            for p in result["projects"]
        ]
        return ProjectsOverviewResponse(
            projects=projects,
            source_notes_count=len(notes),
        )

    # ── 技術能力總覽 ────────────────────────────────────────────────────────
    async def generate_skills(self) -> SkillsResponse:
        notes = await self._note_client.get_notes_by_category("技術筆記")
        if not notes:
            raise ValueError("找不到「技術筆記」類別的筆記，請先在筆記頁面新增技術相關內容")
        result = await self._claude_client.generate_skills(notes)
        return SkillsResponse(
            languages=result["languages"],
            frameworks=result["frameworks"],
            databases=result["databases"],
            tools=result["tools"],
            summary=result["summary"],
            source_notes_count=len(notes),
        )

    # ── 生活雜談 ────────────────────────────────────────────────────────────
    async def generate_daily(self) -> DailyResponse:
        notes = await self._note_client.get_notes_by_category("日常隨筆")
        if not notes:
            raise ValueError("找不到「日常隨筆」類別的筆記，請先在筆記頁面新增生活相關內容")
        result = await self._claude_client.generate_daily(notes)
        return DailyResponse(
            introduction=result["introduction"],
            interests=result["interests"],
            source_notes_count=len(notes),
        )
