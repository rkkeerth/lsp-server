package document

import (
	"strings"
	"sync"

	"go.lsp.dev/protocol"
)

// Document represents a text document with its content and version
type Document struct {
	URI     protocol.DocumentURI
	Content string
	Version int32
	Lines   []string
}

// Manager manages all open documents
type Manager struct {
	documents map[protocol.DocumentURI]*Document
	mu        sync.RWMutex
}

// NewManager creates a new document manager
func NewManager() *Manager {
	return &Manager{
		documents: make(map[protocol.DocumentURI]*Document),
	}
}

// Open adds a new document to the manager
func (m *Manager) Open(uri protocol.DocumentURI, content string, version int32) {
	m.mu.Lock()
	defer m.mu.Unlock()

	doc := &Document{
		URI:     uri,
		Content: content,
		Version: version,
		Lines:   strings.Split(content, "\n"),
	}
	m.documents[uri] = doc
}

// Update updates an existing document's content
func (m *Manager) Update(uri protocol.DocumentURI, content string, version int32) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if doc, exists := m.documents[uri]; exists {
		doc.Content = content
		doc.Version = version
		doc.Lines = strings.Split(content, "\n")
	}
}

// Close removes a document from the manager
func (m *Manager) Close(uri protocol.DocumentURI) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.documents, uri)
}

// Get retrieves a document by URI
func (m *Manager) Get(uri protocol.DocumentURI) (*Document, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	doc, exists := m.documents[uri]
	return doc, exists
}

// GetAll returns all managed documents
func (m *Manager) GetAll() []*Document {
	m.mu.RLock()
	defer m.mu.RUnlock()

	docs := make([]*Document, 0, len(m.documents))
	for _, doc := range m.documents {
		docs = append(docs, doc)
	}
	return docs
}

// GetWordAt retrieves the word at a specific position in a document
func (d *Document) GetWordAt(pos protocol.Position) string {
	if int(pos.Line) >= len(d.Lines) {
		return ""
	}

	line := d.Lines[pos.Line]
	if int(pos.Character) >= len(line) {
		return ""
	}

	// Find word boundaries
	start := int(pos.Character)
	end := int(pos.Character)

	// Move start backward to beginning of word
	for start > 0 && isWordChar(rune(line[start-1])) {
		start--
	}

	// Move end forward to end of word
	for end < len(line) && isWordChar(rune(line[end])) {
		end++
	}

	return line[start:end]
}

// GetLineAt returns the line at the specified position
func (d *Document) GetLineAt(pos protocol.Position) string {
	if int(pos.Line) >= len(d.Lines) {
		return ""
	}
	return d.Lines[pos.Line]
}

// isWordChar determines if a character is part of a word
func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '_'
}
