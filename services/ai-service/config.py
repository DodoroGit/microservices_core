import os


def get_env(key: str, default: str = "") -> str:
    return os.getenv(key, default)


class Settings:
    port: int = int(get_env("PORT", "8082"))
    note_service_url: str = get_env("NOTE_SERVICE_URL", "http://localhost:8083")
    anthropic_api_key: str = get_env("ANTHROPIC_API_KEY", "")
    claude_model: str = get_env("CLAUDE_MODEL", "claude-sonnet-4-6")


settings = Settings()
