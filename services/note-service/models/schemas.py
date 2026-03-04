from datetime import datetime
from enum import Enum
from typing import Optional

from pydantic import BaseModel


class Category(str, Enum):
    project = "專案"
    tech = "技術筆記"
    daily = "日常隨筆"


class NoteCreate(BaseModel):
    title: str
    category: Category
    content: str


class NoteUpdate(BaseModel):
    title: Optional[str] = None
    category: Optional[Category] = None
    content: Optional[str] = None


class NoteResponse(BaseModel):
    id: str
    title: str
    category: str
    content: str
    created_at: datetime
    updated_at: datetime
