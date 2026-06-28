from flask import Flask, render_template, jsonify, request, g
from flask_cors import CORS
from pydantic import BaseModel, Field, ValidationError

import logging
from logging.handlers import TimedRotatingFileHandler
import os
import time

import db

IS_VERCEL = os.environ.get("VERCEL") == "1"

app = Flask(__name__)
CORS(app) # Enable CORS for all routes
# Init the DB
db.init_db()

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
    
    if not IS_VERCEL: # only log to files locally (Vercel is read-only)
        # File logging handler
        log_path = os.path.normpath(os.path.join("logs", 'server.log'))
        os.makedirs(os.path.dirname(f"{log_path}"), exist_ok=True) # ensure the logs directory exists
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

def validate_request(request, requires_content=True):
    """
    Validates the request.

    Returns:

    (False, (response, status_code)) if invalid
    (True, data) if valid
    
    request: Request var
    """
    if not requires_content:
        return True, {}
    # Validate request
    if not request.content_length:
        logger.warning(f"Request with no content length received from IP: {request.remote_addr}")
        logger.debug(f"No content length in request: {request.get_data(as_text=True)[:500]}") # Log the request object for debugging
        return bool(False), (jsonify({"message": "No content length"}), 400)
    # Parse JSON
    data = request.get_json(silent=True)
    if not data:
        logger.warning(f"Request with no content received from IP: {request.remote_addr}")
        logger.debug(f"No content in request: {request.get_data(as_text=True)[:500]}") # Log the request object for debugging
        return False, (jsonify({"message": "Missing content"}), 400)
    
    return True, data

# Handle pydantic validation errors for request data
@app.errorhandler(ValidationError)
def handle_validation_error(e):
    errors = e.errors()
    # True if every error is a max-length violation
    only_too_long = all(
        err.get("type") == "string_too_long"
        for err in errors
    )
    only_too_short = all(
        err.get("type") == "string_too_short"
        for err in errors
    )
    if only_too_long:
        return jsonify({"message": "One or more fields are too long", "details": e.errors()}), 413
    if only_too_short:
        return jsonify({"message": "One or more fields are too short", "details": e.errors()}), 422
    return jsonify({"message": "Invalid request data", "details": e.errors()}), 422

@app.route("/")
def home():
    return render_template("index.html")

# This for the web page
@app.route("/create-pitch", methods=["GET"])
def create_pitch_page():
    return render_template("create-pitch.html")


class CreatePitchRequest(BaseModel):
    title: str = Field(min_length=1, max_length=100)
    topic: str = Field(min_length=1, max_length=50)
    description: str = Field(min_length=1, max_length=5000)

# and this for the API functionality
@app.route("/create-pitch", methods=["PUT"])
def create_pitch():
    ok, result = validate_request(request)
    if not ok:
        response, status = result
        return response, status
    data = CreatePitchRequest.model_validate(result) # type: ignore

    db.create_idea(title=data.title, topic=data.topic, description=data.description, user_id=db.user.id) # type: ignore # set 1 for now, as no real users exist
    return {}, 200

@app.route("/pitches", methods=["GET"])
def get_pitches():
    pitches = db.get_all_ideas_as_dicts(limit=20)
    return jsonify({"pitches": pitches}), 200

@app.route("/pitches/<int:idea_id>", methods=["GET"])
def get_pitch(idea_id: int):
    idea = db.get_idea_dict(idea_id)
    if not idea:
        return jsonify({"message": "Pitch not found"}), 404
    return render_template("pitch.html", idea=idea, comments=db.get_comments_dict(idea_id=idea_id, limit=50))

class AddCommentRequest(BaseModel):
    content: str = Field(min_length=1, max_length=1000)

@app.route("/pitches/<int:idea_id>/comment", methods=["POST"])
def add_comment(idea_id: int):
    ok, result = validate_request(request)
    if not ok:
        response, status = result
        return response, status
    data = AddCommentRequest.model_validate(result) # type: ignore
    
    content = data.content
    db.create_comment(idea_id=idea_id, content=content, user_id=db.user.id) # type: ignore # set 1 for now, as no real users exist
    return {}, 200

def add_test_pitch(title: str, topic: str, description: str, vote_amount: int):
    id = db.create_idea(title=title, topic=topic, description=description, user_id=db.user.id) # type: ignore # set 1 for now, as no real users exist
    if id is None: 
        raise RuntimeError()
    
    db.update_votes(id, amount=vote_amount)

def add_test_pitches():
    add_test_pitch(
        title="AI Meeting Assistant",
        topic="AI",
        description="AI that joins meetings, creates summaries, and automatically generates tasks.",
        vote_amount=421
    )

    add_test_pitch(
        title="BudgetFlow",
        topic="FinTech",
        description="Modern financial planning platform built for freelancers and creators.",
        vote_amount=312
    )

    add_test_pitch(
        title="MedTrack",
        topic="HealthTech",
        description="Patient monitoring system that helps clinics reduce administrative work.",
        vote_amount=198
    )

if db.idea_count() == 0:
    # Add some test data
    add_test_pitches()

if __name__ == "__main__":
    app.run(host="localhost", port=4000)