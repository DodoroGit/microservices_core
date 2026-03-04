from typing import Optional

from fastapi import APIRouter, Depends, Query

from dependencies import get_note_service
from models.schemas import NoteCreate, NoteResponse, NoteUpdate
from services.note_service import NoteService

router = APIRouter(prefix="/notes", tags=["notes"])


def _format(doc: dict) -> dict:
    return {
        "id": str(doc["_id"]),
        "title": doc["title"],
        "category": doc["category"],
        "content": doc["content"],
        "created_at": doc["created_at"],
        "updated_at": doc["updated_at"],
    }


@router.get("", response_model=list[NoteResponse])
async def list_notes(
    category: Optional[str] = Query(None),
    svc: NoteService = Depends(get_note_service),
):
    notes = await svc.get_all(category)
    return [_format(n) for n in notes]


@router.get("/{note_id}", response_model=NoteResponse)
async def get_note(note_id: str, svc: NoteService = Depends(get_note_service)):
    note = await svc.get_by_id(note_id)
    return _format(note)


@router.post("", response_model=NoteResponse, status_code=201)
async def create_note(body: NoteCreate, svc: NoteService = Depends(get_note_service)):
    note = await svc.create(body)
    return _format(note)


@router.put("/{note_id}", response_model=NoteResponse)
async def update_note(
    note_id: str, body: NoteUpdate, svc: NoteService = Depends(get_note_service)
):
    note = await svc.update(note_id, body)
    return _format(note)


@router.delete("/{note_id}", status_code=204)
async def delete_note(note_id: str, svc: NoteService = Depends(get_note_service)):
    await svc.delete(note_id)
