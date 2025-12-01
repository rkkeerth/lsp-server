package server

import (
	"strings"
	"sync"
)

// Document represents an open text document
type Document struct {
	URI        string
	LanguageID string
	Version    int
	Content    string
	Lines      []string
}

// DocumentStore manages open documents
type DocumentStore struct {
	mu        sync.RWMutex
	documents map[string]*Document
}

// NewDocumentStore creates a new document store
func NewDocumentStore() *DocumentStore {
	return &DocumentStore{
		documents: make(map[string]*Document),
	}
}

// Open adds a document to the store
func (ds *DocumentStore) Open(uri, languageID string, version int, text string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	doc := &Document{
		URI:        uri,
		LanguageID: languageID,
		Version:    version,
		Content:    text,
		Lines:      strings.Split(text, "\n"),
	}
	ds.documents[uri] = doc
}

// Update updates a document's content
func (ds *DocumentStore) Update(uri string, version int, text string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if doc, exists := ds.documents[uri]; exists {
		doc.Version = version
		doc.Content = text
		doc.Lines = strings.Split(text, "\n")
	}
}

// Close removes a document from the store
func (ds *DocumentStore) Close(uri string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	delete(ds.documents, uri)
}

// Get retrieves a document from the store
func (ds *DocumentStore) Get(uri string) (*Document, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	doc, exists := ds.documents[uri]
	return doc, exists
}

// GetWordAtPosition retrieves the word at the specified position
func (d *Document) GetWordAtPosition(line, character int) string {
	if line < 0 || line >= len(d.Lines) {
		return ""
	}

	lineText := d.Lines[line]
	if character < 0 || character >= len(lineText) {
		return ""
	}

	// Find word boundaries
	start := character
	for start > 0 && isWordChar(rune(lineText[start-1])) {
		start--
	}

	end := character
	for end < len(lineText) && isWordChar(rune(lineText[end])) {
		end++
	}

	return lineText[start:end]
}

func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}
