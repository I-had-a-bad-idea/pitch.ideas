from flask import Flask, render_template, jsonify, request, g
from flask_cors import CORS
import logging
from logging.handlers import TimedRotatingFileHandler
import os
import time

import db

app = Flask(__name__)
CORS(app) # Enable CORS for all routes

def get_logging_level(config_logging_level: str="INFO"):
    """
    Converts a logging level from the configs to a logging level that can be used by the logging module
    
    All supported inputs: "NOTSET", "DEBUG", "INFO", "WARN", "ERROR", "CRITICAL"
    """
    # Implementation of those I found with a quick google search
    if config_logging_level == "NOTSET":
        return logging.NOTSET
    elif config_logging_level == "DEBUG":
        return logging.DEBUG
    elif config_logging_level == "INFO":
        return logging.INFO
    elif config_logging_level == "WARN":
        return logging.WARN
    elif config_logging_level == "ERROR":
        return logging.ERROR
    elif config_logging_level == "CRITICAL":
        return logging.CRITICAL
    
    else:
        raise ValueError(f"Invalid config-logging-level: {config_logging_level}")

class ColorFormatter(logging.Formatter):

    COLORS = {
        logging.INFO: "\033[37m",     # white
        logging.WARNING: "\033[33m",  # orange/yellow
        logging.ERROR: "\033[31m",    # red
    }
    RESET = "\033[0m"

    def format(self, record):
        color = self.COLORS.get(record.levelno, "")
        message = super().format(record)
        return f"{color}{message}{self.RESET}"

def setup_logging():
    """
    Sets up logging for the server.
    """
    # Load logging configuration
    GENERAL_LOGGING_LEVEL = logging.DEBUG       # general logging level
    CONSOLE_LOGGING_LEVEL = logging.INFO        # stuff that gets logged in the console
    FILE_LOGGING_LEVEL = logging.DEBUG          # stuff that gets logged in the file

    logger = logging.getLogger()
    logger.setLevel(GENERAL_LOGGING_LEVEL)

    formatter = logging.Formatter('[%(asctime)s %(name)s/%(levelname)s]: %(message)s')
    color_formatter = ColorFormatter('[%(asctime)s %(name)s/%(levelname)s]: %(message)s')

    # Console logging handler
    console_handler = logging.StreamHandler()
    console_handler.setLevel(CONSOLE_LOGGING_LEVEL)
    console_handler.setFormatter(color_formatter)
    logger.addHandler(console_handler)
    
    # File logging handler
    log_path = os.path.normpath(os.path.join("logs", 'server.log'))
    os.makedirs(os.path.dirname(f"../{log_path}"), exist_ok=True) # ensure the logs directory exists
    file_handler = TimedRotatingFileHandler(f"../{log_path}", when='midnight', interval=1) # create a new file daily at midnight
    file_handler.setLevel(FILE_LOGGING_LEVEL)
    file_handler.setFormatter(formatter)
    logger.addHandler(file_handler)
    
    # Flask/werkzeug logs only to console, not to file
    werkzeug_logger = logging.getLogger('werkzeug')
    werkzeug_logger.disabled = True
    
    logger.info("Logging initialized")
    return logger

# logger
logger = setup_logging() # setup logger

@app.before_request
def before_request():
    g.start_time = time.time()

@app.after_request
def log_request(response):
    message = None
    if response.status_code != 200:
        try:
            data = response.get_json()
            message = data.get("message") if isinstance(data, dict) else None # get message from response data dict
        except Exception:
            message = None

    duration = (time.time() - g.start_time) * 1000
    user_id = getattr(request, "user_id", None)

    msg_part = f"message: {message}" if message else ""
    user_part = f"user_id: {user_id}" if user_id else ""

    if response.status_code != 200:
        logger.warning(
                '%s - "%s %s %s" - %s - | %dms | %s | %s',
                request.remote_addr,
                request.method,
                request.path,
                request.environ.get("SERVER_PROTOCOL"),
                response.status_code,
                duration,
                user_part,
                msg_part
            )
    else:
        logger.info(
                '%s - "%s %s %s" - %s - | %dms | %s | %s',
                request.remote_addr,
                request.method,
                request.path,
                request.environ.get("SERVER_PROTOCOL"),
                response.status_code,
                duration,
                user_part,
                msg_part
            )

    return response

@app.route("/")
def home():
    return render_template("index.html")



if __name__ == "__main__":
    db.init_db()
    app.run(host="localhost", port=4000)