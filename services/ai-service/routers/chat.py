from fastapi import APIRouter

from models.schemas import ChatRequest, ChatResponse, HistoryResponse, MessageResponse
from services.chat_service import send_message, fetch_history, delete_history

router = APIRouter(prefix="/chat", tags=["chat"])


@router.post("", response_model=ChatResponse)
async def chat(req: ChatRequest):
    return await send_message(req.user_id, req.message)


@router.get("/history", response_model=HistoryResponse)
def get_history(user_id: str):
    return {"history": fetch_history(user_id)}


@router.delete("/history", response_model=MessageResponse)
def clear_history(user_id: str):
    delete_history(user_id)
    return {"message": "對話紀錄已清除"}
