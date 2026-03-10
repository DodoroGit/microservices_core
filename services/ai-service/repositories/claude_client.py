from typing import Protocol

import anthropic


class ClaudeClientProtocol(Protocol):
    async def generate_project_story(
        self, project_title: str, notes: list[dict]
    ) -> dict: ...

    async def generate_background(self, notes: list[dict]) -> dict: ...
    async def generate_projects_overview(self, notes: list[dict]) -> dict: ...
    async def generate_skills(self, notes: list[dict]) -> dict: ...
    async def generate_daily(self, notes: list[dict]) -> dict: ...


class ClaudeClient:
    """透過 Anthropic SDK 呼叫 Claude API，使用 tool use 取得結構化輸出"""

    # ── 舊版工具（保留相容性）──────────────────────────────────────────────
    _TOOL_PROJECT_STORY = {
        "name": "generate_project_story",
        "description": "根據開發筆記產生專業的履歷與面試素材",
        "input_schema": {
            "type": "object",
            "properties": {
                "resume_description": {
                    "type": "string",
                    "description": (
                        "2～3 句話的 STAR 格式專案描述（繁體中文），"
                        "適合直接放入履歷。"
                        "涵蓋：面對的問題(S/T)、採取的行動(A)、達成的成果(R)。"
                    ),
                },
                "tech_keywords": {
                    "type": "array",
                    "items": {"type": "string"},
                    "description": "從筆記中提取的技術關鍵字清單（例如：Go, MongoDB, Docker）",
                },
                "interview_story": {
                    "type": "string",
                    "description": (
                        "3～4 句話的面試故事（繁體中文），"
                        "說明這個專案的挑戰、你的決策思路與學到的東西。"
                    ),
                },
            },
            "required": ["resume_description", "tech_keywords", "interview_story"],
        },
    }

    # ── 個人背景介紹工具 ────────────────────────────────────────────────────
    _TOOL_BACKGROUND = {
        "name": "generate_background",
        "description": "根據工作經歷筆記產生專業的個人背景介紹",
        "input_schema": {
            "type": "object",
            "properties": {
                "background_intro": {
                    "type": "string",
                    "description": (
                        "2～3 段的專業背景介紹（繁體中文）。"
                        "涵蓋：職涯歷程、主要專長領域、工作風格或價值觀。"
                        "語氣專業但不過於制式，適合放在履歷自我介紹或 LinkedIn 摘要。"
                    ),
                },
                "key_roles": {
                    "type": "array",
                    "items": {"type": "string"},
                    "description": "曾擔任或正在擔任的職位、角色（例如：後端工程師、Tech Lead）",
                },
                "years_summary": {
                    "type": "string",
                    "description": "一句話描述年資或職涯階段（例如：5 年後端開發經驗，專注於分散式系統）",
                },
            },
            "required": ["background_intro", "key_roles", "years_summary"],
        },
    }

    # ── 專案介紹工具 ────────────────────────────────────────────────────────
    _TOOL_PROJECTS_OVERVIEW = {
        "name": "generate_projects_overview",
        "description": "根據多筆專案筆記，整理成每個專案各自的技術介紹段落",
        "input_schema": {
            "type": "object",
            "properties": {
                "projects": {
                    "type": "array",
                    "description": "每個專案的獨立介紹，相關筆記應歸併為同一個專案",
                    "items": {
                        "type": "object",
                        "properties": {
                            "project_name": {
                                "type": "string",
                                "description": "專案名稱",
                            },
                            "description": {
                                "type": "string",
                                "description": (
                                    "2～3 句話的 STAR 格式專案描述（繁體中文）。"
                                    "說明專案背景、技術選型與達成的成果。"
                                ),
                            },
                            "tech_keywords": {
                                "type": "array",
                                "items": {"type": "string"},
                                "description": "此專案用到的技術關鍵字",
                            },
                        },
                        "required": ["project_name", "description", "tech_keywords"],
                    },
                },
            },
            "required": ["projects"],
        },
    }

    # ── 技術能力總覽工具 ────────────────────────────────────────────────────
    _TOOL_SKILLS = {
        "name": "generate_skills",
        "description": "根據技術筆記整理出結構化的技術能力總覽",
        "input_schema": {
            "type": "object",
            "properties": {
                "languages": {
                    "type": "array",
                    "items": {"type": "string"},
                    "description": "程式語言（例如：Go, Python, TypeScript）",
                },
                "frameworks": {
                    "type": "array",
                    "items": {"type": "string"},
                    "description": "框架與函式庫（例如：Gin, FastAPI, React）",
                },
                "databases": {
                    "type": "array",
                    "items": {"type": "string"},
                    "description": "資料庫與資料儲存（例如：PostgreSQL, MongoDB, Redis）",
                },
                "tools": {
                    "type": "array",
                    "items": {"type": "string"},
                    "description": "工具與平台（例如：Docker, GitHub Actions, AWS）",
                },
                "summary": {
                    "type": "string",
                    "description": (
                        "一段 2～3 句話的技術能力概述（繁體中文）。"
                        "說明技術廣度、深度與主要專注領域。"
                    ),
                },
            },
            "required": ["languages", "frameworks", "databases", "tools", "summary"],
        },
    }

    # ── 生活雜談工具 ────────────────────────────────────────────────────────
    _TOOL_DAILY = {
        "name": "generate_daily",
        "description": "根據日常隨筆整理出輕鬆的個人生活介紹",
        "input_schema": {
            "type": "object",
            "properties": {
                "introduction": {
                    "type": "string",
                    "description": (
                        "一段輕鬆自然的個人生活介紹（繁體中文）。"
                        "語氣友善，適合放在個人網站或社群媒體自我介紹。"
                        "涵蓋生活興趣、日常習慣或個人特質。"
                    ),
                },
                "interests": {
                    "type": "array",
                    "items": {"type": "string"},
                    "description": "從筆記提取的興趣、嗜好或關鍵詞（例如：攝影, 健身, 閱讀）",
                },
            },
            "required": ["introduction", "interests"],
        },
    }

    def __init__(self, api_key: str, model: str):
        self._client = anthropic.AsyncAnthropic(api_key=api_key)
        self._model = model

    # ── 舊版（保留相容性）──────────────────────────────────────────────────
    async def generate_project_story(
        self, project_title: str, notes: list[dict]
    ) -> dict:
        notes_text = self._format_notes(notes)
        response = await self._client.messages.create(
            model=self._model,
            max_tokens=1024,
            tools=[self._TOOL_PROJECT_STORY],
            tool_choice={"type": "tool", "name": "generate_project_story"},
            messages=[{
                "role": "user",
                "content": (
                    f"請根據以下關於「{project_title}」的開發筆記，"
                    f"產生適合求職使用的履歷描述與面試素材。\n\n"
                    f"開發筆記：\n{notes_text}"
                ),
            }],
        )
        return self._extract_tool_output(response)

    # ── 個人背景介紹 ────────────────────────────────────────────────────────
    async def generate_background(self, notes: list[dict]) -> dict:
        notes_text = self._format_notes(notes)
        response = await self._client.messages.create(
            model=self._model,
            max_tokens=1024,
            tools=[self._TOOL_BACKGROUND],
            tool_choice={"type": "tool", "name": "generate_background"},
            messages=[{
                "role": "user",
                "content": (
                    "以下是我的工作經歷與履歷相關筆記，"
                    "請根據這些內容幫我生成專業的個人背景介紹，"
                    "適合放在履歷或 LinkedIn 自我介紹欄位。\n\n"
                    f"筆記內容：\n{notes_text}"
                ),
            }],
        )
        return self._extract_tool_output(response)

    # ── 專案介紹 ────────────────────────────────────────────────────────────
    async def generate_projects_overview(self, notes: list[dict]) -> dict:
        notes_text = self._format_notes(notes)
        response = await self._client.messages.create(
            model=self._model,
            max_tokens=2048,
            tools=[self._TOOL_PROJECTS_OVERVIEW],
            tool_choice={"type": "tool", "name": "generate_projects_overview"},
            messages=[{
                "role": "user",
                "content": (
                    "以下是我記錄的各個專案開發筆記，"
                    "請將相關筆記歸併為同一個專案，"
                    "每個專案各自生成一段技術介紹。\n\n"
                    f"專案筆記：\n{notes_text}"
                ),
            }],
        )
        return self._extract_tool_output(response)

    # ── 技術能力總覽 ────────────────────────────────────────────────────────
    async def generate_skills(self, notes: list[dict]) -> dict:
        notes_text = self._format_notes(notes)
        response = await self._client.messages.create(
            model=self._model,
            max_tokens=1024,
            tools=[self._TOOL_SKILLS],
            tool_choice={"type": "tool", "name": "generate_skills"},
            messages=[{
                "role": "user",
                "content": (
                    "以下是我記錄的技術學習筆記與技術心得，"
                    "請整理出我的技術能力總覽，"
                    "依程式語言、框架、資料庫、工具分類。\n\n"
                    f"技術筆記：\n{notes_text}"
                ),
            }],
        )
        return self._extract_tool_output(response)

    # ── 生活雜談 ────────────────────────────────────────────────────────────
    async def generate_daily(self, notes: list[dict]) -> dict:
        notes_text = self._format_notes(notes)
        response = await self._client.messages.create(
            model=self._model,
            max_tokens=768,
            tools=[self._TOOL_DAILY],
            tool_choice={"type": "tool", "name": "generate_daily"},
            messages=[{
                "role": "user",
                "content": (
                    "以下是我的日常隨筆與生活記錄，"
                    "請根據這些內容幫我生成一段輕鬆自然的個人生活介紹。\n\n"
                    f"日常隨筆：\n{notes_text}"
                ),
            }],
        )
        return self._extract_tool_output(response)

    # ── 工具函式 ────────────────────────────────────────────────────────────
    def _extract_tool_output(self, response) -> dict:
        for block in response.content:
            if block.type == "tool_use":
                return block.input
        raise ValueError("Claude 未回傳預期的結構化輸出")

    @staticmethod
    def _format_notes(notes: list[dict]) -> str:
        FIELD_LABELS = [
            ("title",        "標題"),
            ("content",      "內容"),
            ("tags",         "技術標籤"),
            ("project_name", "專案名稱"),
            ("description",  "專案介紹"),
            ("tech_stack",   "技術棧"),
            ("github_url",   "GitHub"),
            ("year",         "建置年份"),
            ("company",      "任職公司"),
            ("position",     "任職職位"),
            ("period",       "任職日期"),
            ("projects",     "負責專案"),
        ]
        parts = []
        for i, note in enumerate(notes, 1):
            lines = [f"[筆記 {i}]"]
            for field, label in FIELD_LABELS:
                val = note.get(field)
                if val:
                    if isinstance(val, list):
                        lines.append(f"{label}：{', '.join(str(v) for v in val)}")
                    else:
                        lines.append(f"{label}：{val}")
            parts.append("\n".join(lines))
        return "\n\n".join(parts)
