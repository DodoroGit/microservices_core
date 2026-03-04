import redis

from config import cfg

redis_client = redis.Redis(host=cfg.REDIS_HOST, port=cfg.REDIS_PORT, decode_responses=True)
