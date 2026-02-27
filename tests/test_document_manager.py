"""Tests for the document manager."""
from lsp_server.document.manager import Document, Manager


class TestManagerOpen:
    def test_open_then_get_returns_correct_document(self):
        mgr = Manager()
        mgr.open("file:///test.py", "python", 1, "print('hello')\n")
        doc = mgr.get("file:///test.py")
        assert doc is not None
        assert doc.uri == "file:///test.py"
        assert doc.language_id == "python"
        assert doc.version == 1
        assert doc.content == "print('hello')\n"


class TestManagerGet:
    def test_get_nonexistent_uri_returns_none(self):
        mgr = Manager()
        assert mgr.get("file:///does_not_exist.py") is None

    def test_get_returns_copy_not_reference(self):
        mgr = Manager()
        mgr.open("file:///test.py", "python", 1, "original")
        doc = mgr.get("file:///test.py")
        assert doc is not None
        # Mutate the returned copy
        doc.content = "mutated"
        doc.version = 999
        # The stored document should be unaffected
        stored = mgr.get("file:///test.py")
        assert stored is not None
        assert stored.content == "original"
        assert stored.version == 1


class TestManagerChange:
    def test_change_updates_content_and_version(self):
        mgr = Manager()
        mgr.open("file:///test.py", "python", 1, "old content")
        result = mgr.change("file:///test.py", 2, "new content")
        assert result is True
        doc = mgr.get("file:///test.py")
        assert doc is not None
        assert doc.version == 2
        assert doc.content == "new content"

    def test_change_nonexistent_uri_returns_false(self):
        mgr = Manager()
        result = mgr.change("file:///nope.py", 1, "content")
        assert result is False


class TestManagerClose:
    def test_close_removes_document(self):
        mgr = Manager()
        mgr.open("file:///test.py", "python", 1, "content")
        result = mgr.close("file:///test.py")
        assert result is True
        assert mgr.get("file:///test.py") is None

    def test_close_nonexistent_uri_returns_false(self):
        mgr = Manager()
        result = mgr.close("file:///nope.py")
        assert result is False


class TestManagerAll:
    def test_all_returns_copies_of_all_documents(self):
        mgr = Manager()
        mgr.open("file:///a.py", "python", 1, "aaa")
        mgr.open("file:///b.py", "python", 1, "bbb")
        docs = mgr.all()
        assert len(docs) == 2
        uris = {d.uri for d in docs}
        assert uris == {"file:///a.py", "file:///b.py"}
        # Mutating returned docs should not affect stored ones
        for d in docs:
            d.content = "mutated"
        stored = mgr.get("file:///a.py")
        assert stored is not None
        assert stored.content == "aaa"

    def test_all_on_empty_manager_returns_empty_list(self):
        mgr = Manager()
        docs = mgr.all()
        assert docs == []
