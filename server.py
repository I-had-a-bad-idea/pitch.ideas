from flask import Flask, render_template, jsonify, request, g, make_response
from flask_cors import CORS
from pydantic import BaseModel, Field, StringConstraints, ValidationError
from typing import Annotated
import os
import time
from functools import wraps

import db
import hashing
from logger import logger

SESSION_COOKIE_NAME = "session_id"

IS_VERCEL = os.environ.get("VERCEL") == "1"

app = Flask(__name__)
CORS(app) # Enable CORS for all routes
# Init the DB
db.init_db()


def require_auth(f):
    @wraps(f)
    def wrapper(*args, **kwargs):
        session_id = request.cookies.get("session_id")

        if not session_id:
            return {"error": "unauthorized"}, 401

        user = db.get_user_by_session(session_id)

        if not user:
            return {"error": "expired"}, 401

        setattr(request, "user", user)

        return f(*args, **kwargs)

    return wrapper

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
        logger.warning('%s - "%s %s %s" - %s - | %dms | %s | %s',
                request.remote_addr, request.method, request.path,
                request.environ.get("SERVER_PROTOCOL"),response.status_code,
                duration, user_part, msg_part)
    else:
        logger.info('%s - "%s %s %s" - %s - | %dms | %s | %s',
                request.remote_addr, request.method, request.path,
                request.environ.get("SERVER_PROTOCOL"), response.status_code,
                duration, user_part, msg_part)

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

Username = Annotated[str, StringConstraints(min_length=1, max_length=100)]
Password = Annotated[str, StringConstraints(min_length=1, max_length=100)]

@app.route("/")
def home():
    return render_template("index.html")

# This for the web page
@app.route("/create-pitch", methods=["GET"])
@require_auth
def create_pitch_page():
    return render_template("create-pitch.html")


class CreatePitchRequest(BaseModel):
    title: str = Field(min_length=1, max_length=100)
    topic: str = Field(min_length=1, max_length=50)
    description: str = Field(min_length=1, max_length=5000)

# and this for the API functionality
@app.route("/create-pitch", methods=["PUT"])
@require_auth
def create_pitch():
    ok, result = validate_request(request)
    if not ok:
        response, status = result
        return response, status
    data = CreatePitchRequest.model_validate(result) # type: ignore
    user = getattr(request, "user")

    db.create_idea(title=data.title, topic=data.topic, description=data.description, user_id=user.id) # type: ignore 
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

@app.route("/pitches/<int:idea_id>/upvote", methods=["POST"])
@require_auth
def vote_pitch(idea_id: int):
    user = getattr(request, "user")
    db.vote_idea(idea_id=idea_id, user_id=user.id, value=1) # currently just upvote by 1
    return {}, 200

class AddCommentRequest(BaseModel):
    content: str = Field(min_length=1, max_length=1000)

@app.route("/pitches/<int:idea_id>/comment", methods=["POST"])
@require_auth
def add_comment(idea_id: int):
    ok, result = validate_request(request)
    if not ok:
        response, status = result
        return response, status
    data = AddCommentRequest.model_validate(result) # type: ignore
    
    user = getattr(request, "user")
    content = data.content
    db.create_comment(idea_id=idea_id, content=content, user_id=user.id) # type: ignore
    return {}, 200

class AuthRequest(BaseModel):
    password: Password
    username: Username

# This for the web page
@app.route("/auth/login", methods=["GET"])
def login_page():
    return render_template("login.html")

# and this for the API functionality
@app.route("/auth/login", methods=["POST"])
def login():
    ok, result = validate_request(request)
    if not ok:
        response, status = result
        return response, status
    data = AuthRequest.model_validate(result) # type:ignore

    user = db.get_user_by_username(data.username)
    if not user or not hashing.verify_password(data.password, user.password_hash):
        return jsonify({"message": "Invalid credentials"}), 401
    
    session_id = db.create_session(user_id=user.id, days=7)
    
    resp = make_response({"ok": True})
    resp.set_cookie(
        SESSION_COOKIE_NAME,
        session_id,
        httponly=True,
        secure=True,
        samesite="Lax",
        max_age=60 * 60 * 24 * 7
    )
    return resp
    
# This for the web page
@app.route("/auth/register", methods=["GET"])
def register_page():
    return render_template("register.html")

# and this for the API functionality
@app.route("/auth/register", methods=["POST"])
def register():
    ok, result = validate_request(request)
    if not ok:
        response, status = result
        return response, status
    data = AuthRequest.model_validate(result) # type:ignore

    user = db.create_user(username=data.username, password_hash=hashing.hash_password(data.password))
    if not user:
        return jsonify({"message": "Username already exists"}), 409
    
    session_id = db.create_session(user_id=user.id, days=7)
    
    resp = make_response({"ok": True})
    resp.set_cookie(
        SESSION_COOKIE_NAME,
        session_id,
        httponly=True,
        secure=True,
        samesite="Lax",
        max_age=60 * 60 * 24 * 7
    )
    return resp

# This for the web page
@app.route("/auth/logout", methods=["GET"])
def logout_page():
    return render_template("logout.html")

# and this for the API functionality
@app.route("/auth/logout", methods=["POST"])
@require_auth
def logout():
    session_id = request.cookies.get(SESSION_COOKIE_NAME)

    if session_id:
        db.delete_session(session_id=session_id)

    resp = jsonify({"message": "logged out"})
    resp.delete_cookie(SESSION_COOKIE_NAME)
    return resp

@app.route("/auth/status", methods=["GET"])
def auth_status():
    session_id = request.cookies.get(SESSION_COOKIE_NAME)
    if not session_id:
        return jsonify({"logged_in": False}), 200

    user = db.get_user_by_session(session_id=session_id)
    if not user:
        return jsonify({"logged_in": False}), 200
    

    return jsonify({
        "logged_in": True,
        "user": {
            "id": user.id,
            "username": user.username
        }
    })

def add_test_pitch(title: str, topic: str, description: str, vote_amount: int):
    id = db.create_idea(title=title, topic=topic, description=description, user_id=db.user.id) # type: ignore # Use test user for adding the stuff
    if id is None: 
        raise RuntimeError()
    
    db.vote_idea(idea_id=id, user_id=db.user.id, value=vote_amount) # type: ignore

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