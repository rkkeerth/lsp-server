package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

// handleTextDocumentDidOpen processes didOpen notifications
func (s *Server) handleTextDocumentDidOpen(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidOpenTextDocumentParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling didOpen params: %v", err)
		return reply(ctx, nil, err)
	}

	uri := params.TextDocument.URI
	content := params.TextDocument.Text
	version := params.TextDocument.Version

	s.documents.Set(uri, content, version)
	log.Printf("Document opened: %s (version %d)", uri, version)

	return reply(ctx, nil, nil)
}

// handleTextDocumentDidChange processes didChange notifications
func (s *Server) handleTextDocumentDidChange(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidChangeTextDocumentParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling didChange params: %v", err)
		return reply(ctx, nil, err)
	}

	uri := params.TextDocument.URI
	version := params.TextDocument.Version

	// Since we're using TextDocumentSyncKindFull, there should be one change with the full content
	if len(params.ContentChanges) > 0 {
		// Get the full text from the last change
		newContent := params.ContentChanges[len(params.ContentChanges)-1].Text
		s.documents.Set(uri, newContent, version)
		log.Printf("Document changed: %s (version %d)", uri, version)
	}

	return reply(ctx, nil, nil)
}

// handleTextDocumentDidClose processes didClose notifications
func (s *Server) handleTextDocumentDidClose(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidCloseTextDocumentParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling didClose params: %v", err)
		return reply(ctx, nil, err)
	}

	uri := params.TextDocument.URI
	s.documents.Delete(uri)
	log.Printf("Document closed: %s", uri)

	return reply(ctx, nil, nil)
}

// handleTextDocumentHover processes hover requests
func (s *Server) handleTextDocumentHover(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.HoverParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling hover params: %v", err)
		return reply(ctx, nil, err)
	}

	uri := params.TextDocument.URI
	position := params.Position

	doc, exists := s.documents.Get(uri)
	if !exists {
		log.Printf("Document not found: %s", uri)
		return reply(ctx, nil, nil)
	}

	// Get the word at the position
	word := getWordAtPosition(doc.Content, position)
	log.Printf("Hover request for word: %s at %s:%d:%d", word, uri, position.Line, position.Character)

	// Create hover response with basic information
	hoverContent := fmt.Sprintf("**%s**\n\nHover information for '%s'\n\nLine: %d, Character: %d",
		word, word, position.Line, position.Character)

	result := &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: hoverContent,
		},
	}

	return reply(ctx, result, nil)
}

// handleTextDocumentCompletion processes completion requests
func (s *Server) handleTextDocumentCompletion(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.CompletionParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling completion params: %v", err)
		return reply(ctx, nil, err)
	}

	uri := params.TextDocument.URI
	position := params.Position

	log.Printf("Completion request at %s:%d:%d", uri, position.Line, position.Character)

	// Return some example completion items
	items := []protocol.CompletionItem{
		{
			Label:  "example",
			Kind:   protocol.CompletionItemKindText,
			Detail: "Example completion item",
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "This is an example completion item provided by the LSP server.",
			},
			InsertText: "example",
		},
		{
			Label:  "function",
			Kind:   protocol.CompletionItemKindFunction,
			Detail: "Example function",
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "An example function completion.",
			},
			InsertText: "function()",
		},
		{
			Label:  "variable",
			Kind:   protocol.CompletionItemKindVariable,
			Detail: "Example variable",
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "An example variable completion.",
			},
			InsertText: "variable",
		},
		{
			Label:  "class",
			Kind:   protocol.CompletionItemKindClass,
			Detail: "Example class",
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "An example class completion.",
			},
			InsertText: "class",
		},
	}

	result := protocol.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}

	return reply(ctx, result, nil)
}

// getWordAtPosition extracts the word at the given position in the document
func getWordAtPosition(content string, position protocol.Position) string {
	lines := strings.Split(content, "\n")
	if int(position.Line) >= len(lines) {
		return ""
	}

	line := lines[position.Line]
	if int(position.Character) >= len(line) {
		return ""
	}

	// Find word boundaries
	start := int(position.Character)
	end := int(position.Character)

	// Move start backward to find the beginning of the word
	for start > 0 && isWordChar(rune(line[start-1])) {
		start--
	}

	// Move end forward to find the end of the word
	for end < len(line) && isWordChar(rune(line[end])) {
		end++
	}

	if start >= end {
		return ""
	}

	return line[start:end]
}

// isWordChar determines if a character is part of a word
func isWordChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}
