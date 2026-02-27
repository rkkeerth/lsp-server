"""Document manager for tracking open text documents."""
from __future__ import annotations

import copy
import threading
from dataclasses import dataclass
from typing import Optional


@dataclass
class Document:
    uri: str
    language_id: str
    version: int
    content: str


class Manager:
    """Thread-safe manager for open text documents."""

    def __init__(self) -> None:
        self._lock = threading.Lock()
        self._documents: dict[str, Document] = {}

    def open(self, uri: str, language_id: str, version: int, content: str) -> None:
        """Store a newly opened document."""
        with self._lock:
            self._documents[uri] = Document(
                uri=uri,
                language_id=language_id,
                version=version,
                content=content,
            )

    def change(self, uri: str, version: int, content: str) -> bool:
        """Update an existing document. Returns True if the document existed."""
        with self._lock:
            if uri not in self._documents:
                return False
            doc = self._documents[uri]
            doc.version = version
            doc.content = content
            return True

    def close(self, uri: str) -> bool:
        """Remove a document. Returns True if the document existed."""
        with self._lock:
            if uri not in self._documents:
                return False
            del self._documents[uri]
            return True

    def get(self, uri: str) -> Optional[Document]:
        """Get a copy of a document by URI, or None if not found."""
        with self._lock:
            doc = self._documents.get(uri)
            if doc is None:
                return None
            return copy.copy(doc)

    def all(self) -> list[Document]:
        """Get copies of all open documents."""
        with self._lock:
            return [copy.copy(d) for d in self._documents.values()]
