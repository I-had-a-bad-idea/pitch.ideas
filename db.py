import sqlite3
from dataclasses import dataclass
from datetime import datetime

@dataclass
class User:
    id: int
    username: str
    password_hash: str

@dataclass
class Idea:
    id: int
    title: str
    topic: str
    description: str
    user_id: int
    created_at: datetime
    votes: int

@dataclass
class Comment:
    id: int
    idea_id: int
    user_id: int
    created_at: datetime
    content: str
    votes: int

DB_NAME = "pitch-ideas.db"

def get_connection():
    conn = sqlite3.connect(DB_NAME)
    conn.row_factory = sqlite3.Row
    return conn

def init_db():
    conn = get_connection()
    cursor = conn.cursor()

    cursor.execute("""
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL
    )
    """)

    cursor.execute("""
    CREATE TABLE IF NOT EXISTS ideas (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        topic TEXT NOT NULL,
        description TEXT NOT NULL,
        user_id INTEGER,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        votes INTEGER DEFAULT 0,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )
    """)
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        idea_id INTEGER,
        user_id INTEGER,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        content TEXT NOT NULL,
        votes INTEGER DEFAULT 0,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (idea_Id) REFERENCES ideas(id)
        
    )
    """)

    conn.commit()
    conn.close()

    print("Database and tables created!")

def create_user(username: str, password_hash: str):
    conn = get_connection()
    cursor = conn.cursor()

    cursor.execute(
        "INSERT INTO users (username, password_hash) VALUES (?, ?)",
        (username, password_hash)
    )

    conn.commit()
    conn.close()


def get_user(username: str) -> User:
    conn = get_connection()
    cursor = conn.cursor()

    cursor.execute(
        "SELECT * FROM users WHERE username = ?",
        (username,)
    )

    user = cursor.fetchone()
    conn.close()
    return User(id=user["id"], username=user["username"], password_hash=user["password_hash"])



def create_idea(title: str, topic: str, description: str, user_id: int):
    conn = get_connection()
    cursor = conn.cursor()

    cursor.execute("""
        INSERT INTO ideas (title, topic, description, user_id)
        VALUES (?, ?, ?, ?)
    """, (title, topic, description, user_id))

    conn.commit()
    conn.close()

def row_to_idea(row: sqlite3.Row) -> Idea:
    return Idea(
        id=row["id"],
        title=row["title"],
        topic=row["topic"],
        description=row["description"],
        user_id=row["user_id"],
        created_at=datetime.fromisoformat(row["created_at"]),
        votes=row["votes"],
    )

def get_all_ideas(limit: int = 20) -> list[Idea]:
    """Currently no limit"""
    conn = get_connection()
    cursor = conn.cursor()

    cursor.execute("""
        SELECT * FROM ideas
        ORDER BY created_at DESC
        LIMIT ?
    """, (limit,))

    rows = cursor.fetchall()
    conn.close()
    return [row_to_idea(row) for row in rows]


def get_idea(idea_id: int) -> Idea | None:
    conn = get_connection()
    cursor = conn.cursor()

    cursor.execute(
        "SELECT * FROM ideas WHERE id = ?",
        (idea_id,)
    )

    row = cursor.fetchone()
    conn.close()
    if row is None:
        return None

    return row_to_idea(row)


def update_votes(idea_id: int, amount: int):
    """Adds amount to the votes of the idea"""
    conn = get_connection()
    cursor = conn.cursor()

    cursor.execute("""
        UPDATE ideas
        SET votes = votes + ?
        WHERE id = ?
    """, (amount, idea_id))

    conn.commit()
    conn.close()

def create_comment(idea_id: int, user_id: int, content: str):
    conn = get_connection()
    cursor = conn.cursor()

    cursor.execute("""
        INSERT INTO comments (idea_id, user_id, content)
        VALUES (?, ?, ?)
    """, (idea_id, user_id, content))

    conn.commit()
    conn.close()

def row_to_comment(row: sqlite3.Row) -> Comment:
    return Comment(
        id=row["id"],
        idea_id=row["idea_id"],
        user_id=row["user_id"],
        created_at=datetime.fromisoformat(row["created_at"]),
        content=row["content"],
        votes=row["votes"],
    )

def get_comments(idea_id: int, limit: int = 50) -> list[Comment]:
    conn = get_connection()
    cursor = conn.cursor()

    cursor.execute("""
        SELECT * FROM comments
        WHERE idea_id = ?
        ORDER BY created_at ASC
        LIMIT ?
    """, (idea_id, limit))

    rows = cursor.fetchall()
    conn.close()
    return [row_to_comment(row) for row in rows]


def update_comment_votes(comment_id: int, amount: int):
    """Amount gets added to current comment votes"""
    conn = get_connection()
    cursor = conn.cursor()

    cursor.execute("""
        UPDATE comments
        SET votes = votes + ?
        WHERE id = ?
    """, (amount, comment_id))

    conn.commit()
    conn.close()

