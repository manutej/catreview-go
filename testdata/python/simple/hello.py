"""Simple Python module for testing categorical extraction."""

from typing import Optional


def greet(name: str) -> str:
    """Return a greeting message."""
    message = format_greeting(name)
    return message


def format_greeting(name: str) -> str:
    """Format the greeting with the name."""
    return f"Hello, {name}!"


class Greeter:
    """A simple greeter class."""

    def __init__(self, default_name: str = "World"):
        """Initialize with default name."""
        self.default_name = default_name

    def greet(self, name: Optional[str] = None) -> str:
        """Greet with provided name or default."""
        actual_name = name if name else self.default_name
        return format_greeting(actual_name)


if __name__ == "__main__":
    greeter = Greeter()
    print(greeter.greet("Python"))
