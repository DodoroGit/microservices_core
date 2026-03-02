from fastapi import APIRouter

from models.schemas import ChatRequest, ChatResponse, HistoryResponse, MessageResponse
from services import chat_service

router = APIRouter(prefix="/chat", tags=["chat"])


@router.post("", response_model=ChatResponse)
async def chat(req: ChatRequest):
    return await chat_service.chat(req.user_id, req.message)


@router.get("/history", response_model=HistoryResponse)
def get_history(user_id: str):
    return {"history": chat_service.get_history(user_id)}


@router.delete("/history", response_model=MessageResponse)
def clear_history(user_id: str):
    chat_service.clear_history(user_id)
    return {"message": "對話紀錄已清除"}
