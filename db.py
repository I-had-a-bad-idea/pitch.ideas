import os
from datetime import datetime
from sqlalchemy import (
    create_engine,
    ForeignKey,
    Integer,
    String,
    Text,
    DateTime,
    func
)
from sqlalchemy.orm import (
    DeclarativeBase,
    Mapped,
    mapped_column,
    relationship,
    Session
)

# Check for Vercel Postgres environment variable
database_url = os.environ.get("DATABASE_URL")

if database_url:
    # Fix Vercel/SQLAlchemy compatibility issue
    if database_url.startswith("postgres://"):
        database_url = database_url.replace("postgres://", "postgresql://", 1)
else:
    database_url = "sqlite:///local_development.db" # Fallback to local SQLite for personal PC   (personal Personal Compputer :) )

engine = create_engine(database_url)

class Base(DeclarativeBase):
    pass

class User(Base):
    __tablename__ = "users"

    id: Mapped[int] = mapped_column(primary_key=True)
    username: Mapped[str] = mapped_column(String, unique=True)
    password_hash: Mapped[str] = mapped_column(String)

    ideas: Mapped[list["Idea"]] = relationship(back_populates="user")
    comments: Mapped[list["Comment"]] = relationship(back_populates="user")

user : User | None = None

class Idea(Base):
    __tablename__ = "ideas"

    id: Mapped[int] = mapped_column(primary_key=True)
    title: Mapped[str] = mapped_column(String)
    topic: Mapped[str] = mapped_column(String)
    description: Mapped[str] = mapped_column(Text)

    user_id: Mapped[int] = mapped_column(ForeignKey("users.id"))

    created_at: Mapped[datetime] = mapped_column(
        DateTime,
        server_default=func.now(),
    )

    votes: Mapped[int] = mapped_column(Integer, default=0)

    user: Mapped["User"] = relationship(back_populates="ideas")
    comments: Mapped[list["Comment"]] = relationship(
        back_populates="idea",
        cascade="all, delete-orphan",
    )

    def to_dict(self):
        return {
            "id": self.id,
            "title": self.title,
            "topic": self.topic,
            "description": self.description,
            "user_id": self.user_id,
            "user_name": self.user.username if self.user else None,
            "created_at": self.created_at.isoformat(),
            "created_at_pretty": self.created_at.strftime("%d %b %Y, %H:%M"),
            "votes": self.votes,
            "comment_count": len(self.comments),
        }

class Comment(Base):
    __tablename__ = "comments"

    id: Mapped[int] = mapped_column(primary_key=True)

    idea_id: Mapped[int] = mapped_column(ForeignKey("ideas.id"))
    user_id: Mapped[int] = mapped_column(ForeignKey("users.id"))

    created_at: Mapped[datetime] = mapped_column(
        DateTime,
        server_default=func.now(),
    )

    content: Mapped[str] = mapped_column(Text)
    votes: Mapped[int] = mapped_column(Integer, default=0)

    idea: Mapped["Idea"] = relationship(back_populates="comments")
    user: Mapped["User"] = relationship(back_populates="comments")

    def to_dict(self) -> dict:
        return {
            "id": self.id,
            "idea_id": self.idea_id,
            "user_id": self.user_id,
            "user_name": self.user.username if self.user else None,
            "created_at": self.created_at.isoformat(),
            "created_at_pretty": self.created_at.strftime("%d %b %Y, %H:%M"),
            "content": self.content,
            "votes": self.votes,
        }


def get_session():
    return Session(engine)

def init_db():
    global user
    Base.metadata.create_all(engine)
    
    _user = get_user("test")
    if _user is None:
        _user = create_user("test", "test") # TODO: Remove this if users get added (only for now, where everything gets attributed to user 0)

    user = _user
    print("Database and tables created!")

def idea_count() -> int:
    with get_session() as session:
        return session.query(Idea).count()

def create_user(username: str, password_hash: str) -> User:
    with get_session() as session:
        user = User(
            username=username,
            password_hash=password_hash
        )
        session.add(user)
        session.commit()
        session.refresh(user)
        return user

def get_user(username: str) -> User | None:
    with get_session() as session:
        return (
            session.query(User)
            .filter_by(username=username)
            .first()
        )



def create_idea(title: str, topic: str, description: str, user_id: int) -> int | None:
    """Creates a pitch, returns id"""
    with get_session() as session:
        idea = Idea(
            title=title,
            topic=topic,
            description=description,
            user_id=user_id,
        )
        session.add(idea)
        session.commit()
        session.refresh(idea)
        return idea.id

def get_all_ideas_as_dicts(limit: int = 20) -> list[dict]:
    with get_session() as session:
        ideas = (
            session.query(Idea)
            .order_by(Idea.created_at.desc())
            .limit(limit)
            .all()
        )
        return [idea.to_dict() for idea in ideas]



def get_idea_dict(idea_id: int) -> dict | None:
    with get_session() as session:
        idea = session.get(Idea, idea_id)
        if idea:
            return idea.to_dict()


def update_votes(idea_id: int, amount: int):
    """Adds amount to the votes of the idea"""
    with get_session() as session:
        idea = session.get(Idea, idea_id)
        if idea:
            idea.votes += amount
            session.commit()

def create_comment(idea_id: int, user_id: int, content: str):
    with get_session() as session:
        session.add(
            Comment(
                idea_id=idea_id,
                user_id=user_id,
                content=content
            )
        )
        session.commit()

def get_comment_count(idea_id: int) -> int:
    with get_session() as session:
        return (
            session.query(Comment)
            .filter_by(idea_id=idea_id)
            .count()
        )

def get_comments_dict(idea_id: int, limit: int = 50) -> list[dict]:
    with get_session() as session:
        comments =  (
            session.query(Comment)
            .filter_by(idea_id=idea_id)
            .order_by(Comment.created_at)
            .limit(limit)
            .all()
        )
        return [comment.to_dict() for comment in comments]

def update_comment_votes(comment_id: int, amount: int):
    """Amount gets added to current comment votes"""
    with get_session() as session:
        comment = session.get(Comment, comment_id)
        if comment:
            comment.votes += amount
            session.commit()

