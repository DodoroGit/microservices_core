from fastapi import APIRouter, Depends

from models.schemas import ChatRequest, ChatResponse, HistoryResponse, MessageResponse
from routers.dependencies import get_chat_service
from services.chat_service import ChatService

router = APIRouter(prefix="/chat", tags=["chat"])


@router.post("", response_model=ChatResponse)
async def chat(req: ChatRequest, svc: ChatService = Depends(get_chat_service)):
    return await svc.send_message(req.user_id, req.message)


@router.get("/history", response_model=HistoryResponse)
def get_history(user_id: str, svc: ChatService = Depends(get_chat_service)):
    return {"history": svc.fetch_history(user_id)}


@router.delete("/history", response_model=MessageResponse)
def clear_history(user_id: str, svc: ChatService = Depends(get_chat_service)):
    svc.delete_history(user_id)
    return {"message": "對話紀錄已清除"}
