package server

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/rkkeerth/lsp-server/protocol"
	"github.com/sourcegraph/jsonrpc2"
)

func TestServerInitialize(t *testing.T) {
	logger := log.New(os.Stderr, "[TEST] ", log.Lshortfile)
	srv := NewServer(logger)

	params := protocol.InitializeParams{
		RootURI: "file:///tmp",
		Capabilities: protocol.ClientCapabilities{
			TextDocument: protocol.TextDocumentClientCapabilities{
				Hover: protocol.HoverCapabilities{
					ContentFormat: []string{"markdown"},
				},
			},
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	rawParams := json.RawMessage(paramsJSON)
	req := &jsonrpc2.Request{
		Method: "initialize",
		Params: &rawParams,
	}

	result, err := srv.handleInitialize(context.Background(), req)
	if err != nil {
		t.Fatalf("handleInitialize failed: %v", err)
	}

	initResult, ok := result.(protocol.InitializeResult)
	if !ok {
		t.Fatalf("Expected InitializeResult, got %T", result)
	}

	if !initResult.Capabilities.HoverProvider {
		t.Error("Expected HoverProvider to be true")
	}

	if initResult.Capabilities.TextDocumentSync.OpenClose != true {
		t.Error("Expected TextDocumentSync.OpenClose to be true")
	}

	if initResult.ServerInfo == nil || initResult.ServerInfo.Name != "basic-lsp-server" {
		t.Error("Expected server info to be set correctly")
	}
}

func TestDocumentOperations(t *testing.T) {
	logger := log.New(os.Stderr, "[TEST] ", log.Lshortfile)
	srv := NewServer(logger)

	// Test didOpen
	openParams := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///tmp/test.txt",
			LanguageID: "text",
			Version:    1,
			Text:       "Hello world\nTest line",
		},
	}

	paramsJSON, _ := json.Marshal(openParams)
	rawParams := json.RawMessage(paramsJSON)
	req := &jsonrpc2.Request{
		Method: "textDocument/didOpen",
		Params: &rawParams,
	}

	_, err := srv.handleTextDocumentDidOpen(context.Background(), req)
	if err != nil {
		t.Fatalf("handleTextDocumentDidOpen failed: %v", err)
	}

	// Verify document was stored
	doc, exists := srv.documents.Get("file:///tmp/test.txt")
	if !exists {
		t.Fatal("Document was not stored")
	}

	if doc.Content != "Hello world\nTest line" {
		t.Errorf("Expected content 'Hello world\\nTest line', got '%s'", doc.Content)
	}

	if len(doc.Lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(doc.Lines))
	}

	// Test didChange
	changeParams := protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			URI:     "file:///tmp/test.txt",
			Version: 2,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{
				Text: "Updated content\nNew line",
			},
		},
	}

	paramsJSON, _ = json.Marshal(changeParams)
	rawParams = json.RawMessage(paramsJSON)
	req = &jsonrpc2.Request{
		Method: "textDocument/didChange",
		Params: &rawParams,
	}

	_, err = srv.handleTextDocumentDidChange(context.Background(), req)
	if err != nil {
		t.Fatalf("handleTextDocumentDidChange failed: %v", err)
	}

	// Verify document was updated
	doc, _ = srv.documents.Get("file:///tmp/test.txt")
	if doc.Content != "Updated content\nNew line" {
		t.Errorf("Expected updated content, got '%s'", doc.Content)
	}

	if doc.Version != 2 {
		t.Errorf("Expected version 2, got %d", doc.Version)
	}

	// Test didClose
	closeParams := protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: "file:///tmp/test.txt",
		},
	}

	paramsJSON, _ = json.Marshal(closeParams)
	rawParams = json.RawMessage(paramsJSON)
	req = &jsonrpc2.Request{
		Method: "textDocument/didClose",
		Params: &rawParams,
	}

	_, err = srv.handleTextDocumentDidClose(context.Background(), req)
	if err != nil {
		t.Fatalf("handleTextDocumentDidClose failed: %v", err)
	}

	// Verify document was removed
	_, exists = srv.documents.Get("file:///tmp/test.txt")
	if exists {
		t.Error("Document should have been removed")
	}
}

func TestHover(t *testing.T) {
	logger := log.New(os.Stderr, "[TEST] ", log.Lshortfile)
	srv := NewServer(logger)

	// Open a document first
	srv.documents.Open("file:///tmp/test.txt", "go", 1, "package main\n\nfunc Hello() {}")

	hoverParams := protocol.HoverParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: "file:///tmp/test.txt",
		},
		Position: protocol.Position{
			Line:      2,
			Character: 5,
		},
	}

	paramsJSON, _ := json.Marshal(hoverParams)
	rawParams := json.RawMessage(paramsJSON)
	req := &jsonrpc2.Request{
		Method: "textDocument/hover",
		Params: &rawParams,
	}

	result, err := srv.handleTextDocumentHover(context.Background(), req)
	if err != nil {
		t.Fatalf("handleTextDocumentHover failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected hover result, got nil")
	}

	hover, ok := result.(protocol.Hover)
	if !ok {
		t.Fatalf("Expected Hover result, got %T", result)
	}

	if hover.Contents.Kind != protocol.Markdown {
		t.Errorf("Expected markdown content, got %s", hover.Contents.Kind)
	}

	if hover.Contents.Value == "" {
		t.Error("Expected non-empty hover content")
	}
}

func TestGetWordAtPosition(t *testing.T) {
	doc := &Document{
		URI:        "file:///test.txt",
		LanguageID: "text",
		Version:    1,
		Content:    "Hello world_test\nAnother line",
		Lines:      []string{"Hello world_test", "Another line"},
	}

	tests := []struct {
		line      int
		character int
		expected  string
	}{
		{0, 0, "Hello"},
		{0, 2, "Hello"},
		{0, 4, "Hello"},
		{0, 6, "world_test"},
		{0, 11, "world_test"},
		{1, 0, "Another"},
		{1, 8, "line"},
		{0, 100, ""}, // Out of bounds
		{10, 0, ""},  // Invalid line
	}

	for _, tt := range tests {
		result := doc.GetWordAtPosition(tt.line, tt.character)
		if result != tt.expected {
			t.Errorf("GetWordAtPosition(%d, %d) = %q, expected %q",
				tt.line, tt.character, result, tt.expected)
		}
	}
}
