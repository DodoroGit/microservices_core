from pydantic import BaseModel


# ── 個人背景介紹（履歷類別）──────────────────────────────────────────────────

class BackgroundResponse(BaseModel):
    background_intro: str    # 2～3 段的專業背景介紹
    key_roles: list[str]     # 曾擔任的角色 / 職位
    years_summary: str       # 年資或職涯階段的一句話摘要
    source_notes_count: int


# ── 專案介紹（專案類別）──────────────────────────────────────────────────────

class ProjectSummary(BaseModel):
    project_name: str
    description: str         # 2～3 句 STAR 格式的專案描述
    tech_keywords: list[str]


class ProjectsOverviewResponse(BaseModel):
    projects: list[ProjectSummary]
    source_notes_count: int


# ── 技術能力總覽（技術筆記類別）──────────────────────────────────────────────

class SkillsResponse(BaseModel):
    languages: list[str]
    frameworks: list[str]
    databases: list[str]
    tools: list[str]
    summary: str             # 一段文字描述的技術能力概述
    source_notes_count: int


# ── 生活雜談（日常隨筆類別）──────────────────────────────────────────────────

class DailyResponse(BaseModel):
    introduction: str        # 輕鬆的個人生活介紹段落
    interests: list[str]     # 從筆記提取的興趣 / 關鍵詞
    source_notes_count: int
