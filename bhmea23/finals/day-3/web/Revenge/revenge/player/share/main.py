from flask import Flask
from website import createApp
from website.config import APP_SECRET
app = createApp()
app.secret_key = APP_SECRET


if __name__ == '__main__':
    app.run("0.0.0.0", port=8080, debug=False)
