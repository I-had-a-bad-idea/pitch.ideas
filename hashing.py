from argon2 import PasswordHasher
from argon2.exceptions import VerifyMismatchError, InvalidHashError



ph = PasswordHasher()

def verify_password(password: str, password_hash: str) -> bool:
    try:
        ph.verify(password_hash, password)
        return True
    except (VerifyMismatchError, InvalidHashError):
        return False
    
def hash_password(password) -> str:
    return ph.hash(password)

def is_valid_username(username: str) -> bool:
    return (username.isascii()
            and all(ch.isalnum() or ch in "_-" for ch in username))
