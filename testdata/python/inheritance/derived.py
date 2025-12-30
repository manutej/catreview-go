"""Derived class for testing inheritance extraction."""

from typing import List
from .base import BaseService


class UserService(BaseService):
    """User service extending BaseService."""

    def __init__(self, name: str):
        """Initialize user service."""
        super().__init__(name)
        self.users: List[str] = []

    def add_user(self, user: str) -> None:
        """Add a user."""
        self.users.append(user)
        self.log(f"Added user: {user}")

    def get_users(self) -> List[str]:
        """Get all users."""
        return self.users
