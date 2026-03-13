from fastapi import APIRouter, Depends, HTTPException

from clients.claude_client import ClaudeClient
from clients.note_client import NoteClient
from config import settings
from models.schemas import (
    BackgroundResponse,
    DailyResponse,
    ProjectsOverviewResponse,
    SkillsResponse,
)
from services.ai_service import AIService, AIServiceProtocol


def get_ai_service() -> AIServiceProtocol:
    note_client = NoteClient(settings.note_service_url)
    claude_client = ClaudeClient(settings.anthropic_api_key, settings.claude_model)
    return AIService(note_client, claude_client)

router = APIRouter(prefix="/ai", tags=["ai"])


@router.post("/background", response_model=BackgroundResponse)
async def background(svc: AIServiceProtocol = Depends(get_ai_service)):
    try:
        return await svc.generate_background()
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))


@router.post("/projects", response_model=ProjectsOverviewResponse)
async def projects_overview(svc: AIServiceProtocol = Depends(get_ai_service)):
    try:
        return await svc.generate_projects_overview()
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))


@router.post("/skills", response_model=SkillsResponse)
async def skills(svc: AIServiceProtocol = Depends(get_ai_service)):
    try:
        return await svc.generate_skills()
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))


@router.post("/daily", response_model=DailyResponse)
async def daily(svc: AIServiceProtocol = Depends(get_ai_service)):
    try:
        return await svc.generate_daily()
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))

