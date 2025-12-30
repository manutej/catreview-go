"""Base class for testing inheritance extraction."""


class BaseService:
    """Base service with common functionality."""

    def __init__(self, name: str):
        """Initialize with name."""
        self.name = name

    def log(self, message: str) -> None:
        """Log a message."""
        print(f"[{self.name}] {message}")
