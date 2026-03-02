from pydantic import BaseModel


class ChatRequest(BaseModel):
    message: str
    user_id: str


class Message(BaseModel):
    role: str
    content: str


class ChatResponse(BaseModel):
    response: str
    history: list[Message]


class HistoryResponse(BaseModel):
    history: list[Message]


class MessageResponse(BaseModel):
    message: str
