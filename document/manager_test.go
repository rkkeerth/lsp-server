package document

import (
	"testing"

	"go.lsp.dev/protocol"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Fatal("NewManager returned nil")
	}
	if manager.documents == nil {
		t.Fatal("Manager documents map is nil")
	}
}

func TestManagerOpen(t *testing.T) {
	manager := NewManager()
	uri := protocol.DocumentURI("file:///test.go")
	content := "package main\n\nfunc main() {}"
	version := int32(1)

	manager.Open(uri, content, version)

	doc, exists := manager.Get(uri)
	if !exists {
		t.Fatal("Document not found after opening")
	}
	if doc.URI != uri {
		t.Errorf("Expected URI %s, got %s", uri, doc.URI)
	}
	if doc.Content != content {
		t.Errorf("Expected content %s, got %s", content, doc.Content)
	}
	if doc.Version != version {
		t.Errorf("Expected version %d, got %d", version, doc.Version)
	}
	if len(doc.Lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(doc.Lines))
	}
}

func TestManagerUpdate(t *testing.T) {
	manager := NewManager()
	uri := protocol.DocumentURI("file:///test.go")
	
	manager.Open(uri, "initial content", 1)
	
	newContent := "updated content"
	newVersion := int32(2)
	manager.Update(uri, newContent, newVersion)

	doc, exists := manager.Get(uri)
	if !exists {
		t.Fatal("Document not found after update")
	}
	if doc.Content != newContent {
		t.Errorf("Expected content %s, got %s", newContent, doc.Content)
	}
	if doc.Version != newVersion {
		t.Errorf("Expected version %d, got %d", newVersion, doc.Version)
	}
}

func TestManagerClose(t *testing.T) {
	manager := NewManager()
	uri := protocol.DocumentURI("file:///test.go")
	
	manager.Open(uri, "content", 1)
	manager.Close(uri)

	_, exists := manager.Get(uri)
	if exists {
		t.Error("Document still exists after closing")
	}
}

func TestManagerGetAll(t *testing.T) {
	manager := NewManager()
	
	manager.Open(protocol.DocumentURI("file:///test1.go"), "content1", 1)
	manager.Open(protocol.DocumentURI("file:///test2.go"), "content2", 1)
	manager.Open(protocol.DocumentURI("file:///test3.go"), "content3", 1)

	docs := manager.GetAll()
	if len(docs) != 3 {
		t.Errorf("Expected 3 documents, got %d", len(docs))
	}
}

func TestDocumentGetWordAt(t *testing.T) {
	doc := &Document{
		Content: "package main\n\nfunc hello() {}",
		Lines:   []string{"package main", "", "func hello() {}"},
	}

	tests := []struct {
		name     string
		pos      protocol.Position
		expected string
	}{
		{
			name:     "word at beginning",
			pos:      protocol.Position{Line: 0, Character: 0},
			expected: "package",
		},
		{
			name:     "word in middle",
			pos:      protocol.Position{Line: 0, Character: 8},
			expected: "main",
		},
		{
			name:     "function name",
			pos:      protocol.Position{Line: 2, Character: 5},
			expected: "hello",
		},
		{
			name:     "empty line",
			pos:      protocol.Position{Line: 1, Character: 0},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := doc.GetWordAt(tt.pos)
			if result != tt.expected {
				t.Errorf("Expected word %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDocumentGetLineAt(t *testing.T) {
	doc := &Document{
		Lines: []string{"line 0", "line 1", "line 2"},
	}

	tests := []struct {
		name     string
		pos      protocol.Position
		expected string
	}{
		{
			name:     "first line",
			pos:      protocol.Position{Line: 0},
			expected: "line 0",
		},
		{
			name:     "middle line",
			pos:      protocol.Position{Line: 1},
			expected: "line 1",
		},
		{
			name:     "last line",
			pos:      protocol.Position{Line: 2},
			expected: "line 2",
		},
		{
			name:     "out of bounds",
			pos:      protocol.Position{Line: 10},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := doc.GetLineAt(tt.pos)
			if result != tt.expected {
				t.Errorf("Expected line %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestIsWordChar(t *testing.T) {
	tests := []struct {
		char     rune
		expected bool
	}{
		{'a', true},
		{'Z', true},
		{'5', true},
		{'_', true},
		{' ', false},
		{'-', false},
		{'(', false},
		{'.', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			result := isWordChar(tt.char)
			if result != tt.expected {
				t.Errorf("isWordChar(%c) = %v, expected %v", tt.char, result, tt.expected)
			}
		})
	}
}
