from fastapi import APIRouter, Depends, HTTPException

from models.schemas import (
    BackgroundResponse,
    DailyResponse,
    ProjectStoryRequest,
    ProjectStoryResponse,
    ProjectsOverviewResponse,
    SkillsResponse,
)
from routers.dependencies import get_ai_service
from services.ai_service import AIServiceProtocol

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


# ── 舊版（保留相容性）──────────────────────────────────────────────────────
@router.post("/project-story", response_model=ProjectStoryResponse)
async def project_story(
    req: ProjectStoryRequest,
    svc: AIServiceProtocol = Depends(get_ai_service),
):
    try:
        return await svc.generate_project_story(req.project_title)
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))
