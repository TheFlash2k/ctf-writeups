import os

# app secret session
APP_SECRET = os.urandom(12).hex()
# DB config
DB_NAME = os.environ.get("DB_NAME") or "db"
DB_USER = os.environ.get("DB_USER") or "flaskuser"
DB_PWD = os.environ.get("PASSWORD") or "REDACTED"
DB_HOST = os.environ.get("DB_HOST") or "localhost"
DB_PORT = os.environ.get("DB_PORT") or "5432"

# SALT
SALT = os.urandom(10).hex()

# BOT THREAD, higher value results in more memory resource consumption and less time to wait for bot to visit a link.
# LOWER valur result in less resource consumption but player need to wait for bot to visit his link
MAX_CONCURRENT_THREADS = 12

# USER agent, to identify the bot
SECRET_AGENT = os.environ.get("SECRET_AGENT") or "REDACTED"

# FLAG
FLAG = os.environ.get("FLAG") or "REDACTED"
