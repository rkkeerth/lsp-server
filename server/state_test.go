package server

import "testing"

func TestOpenDocument(t *testing.T) {
	state := NewState()

	uri := "file:///test.txt"
	content := "hello world"

	state.OpenDocument(uri, content)

	got, ok := state.GetDocument(uri)
	if !ok {
		t.Fatal("Document should exist after opening")
	}
	if got != content {
		t.Errorf("Content mismatch: got %q, want %q", got, content)
	}
}

func TestUpdateDocument(t *testing.T) {
	state := NewState()

	uri := "file:///test.txt"
	initialContent := "hello"
	updatedContent := "hello world"

	state.OpenDocument(uri, initialContent)
	state.UpdateDocument(uri, updatedContent)

	got, ok := state.GetDocument(uri)
	if !ok {
		t.Fatal("Document should exist after update")
	}
	if got != updatedContent {
		t.Errorf("Content mismatch: got %q, want %q", got, updatedContent)
	}
}

func TestGetDocumentNotFound(t *testing.T) {
	state := NewState()

	_, ok := state.GetDocument("file:///nonexistent.txt")
	if ok {
		t.Error("GetDocument should return false for missing document")
	}
}

func TestCloseDocument(t *testing.T) {
	state := NewState()

	uri := "file:///test.txt"
	content := "hello world"

	state.OpenDocument(uri, content)
	state.CloseDocument(uri)

	_, ok := state.GetDocument(uri)
	if ok {
		t.Error("Document should not exist after closing")
	}
}

func TestMultipleDocuments(t *testing.T) {
	state := NewState()

	uri1 := "file:///test1.txt"
	uri2 := "file:///test2.txt"
	content1 := "content one"
	content2 := "content two"

	state.OpenDocument(uri1, content1)
	state.OpenDocument(uri2, content2)

	got1, ok1 := state.GetDocument(uri1)
	got2, ok2 := state.GetDocument(uri2)

	if !ok1 || !ok2 {
		t.Fatal("Both documents should exist")
	}
	if got1 != content1 {
		t.Errorf("Document 1 content mismatch: got %q, want %q", got1, content1)
	}
	if got2 != content2 {
		t.Errorf("Document 2 content mismatch: got %q, want %q", got2, content2)
	}
}
