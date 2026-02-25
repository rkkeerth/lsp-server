// Package document provides document state management.
package document

import "sync"

// Document represents an open text document.
type Document struct {
	URI        string
	LanguageID string
	Version    int
	Content    string
}

// Manager manages open documents.
type Manager struct {
	mu        sync.RWMutex
	documents map[string]*Document
}

// NewManager creates a new document manager.
func NewManager() *Manager {
	return &Manager{
		documents: make(map[string]*Document),
	}
}

// Open opens a document and stores it in the manager.
func (m *Manager) Open(uri, languageID string, version int, content string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.documents[uri] = &Document{
		URI:        uri,
		LanguageID: languageID,
		Version:    version,
		Content:    content,
	}
}

// Change updates the content of an open document.
// For full sync mode, the entire content is replaced.
func (m *Manager) Change(uri string, version int, content string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	doc, ok := m.documents[uri]
	if !ok {
		return false
	}

	doc.Version = version
	doc.Content = content
	return true
}

// Close removes a document from the manager.
func (m *Manager) Close(uri string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.documents[uri]; !ok {
		return false
	}

	delete(m.documents, uri)
	return true
}

// Get retrieves a document by URI.
func (m *Manager) Get(uri string) (*Document, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	doc, ok := m.documents[uri]
	if !ok {
		return nil, false
	}

	// Return a copy to avoid race conditions
	return &Document{
		URI:        doc.URI,
		LanguageID: doc.LanguageID,
		Version:    doc.Version,
		Content:    doc.Content,
	}, true
}

// All returns all open documents.
func (m *Manager) All() []*Document {
	m.mu.RLock()
	defer m.mu.RUnlock()

	docs := make([]*Document, 0, len(m.documents))
	for _, doc := range m.documents {
		docs = append(docs, &Document{
			URI:        doc.URI,
			LanguageID: doc.LanguageID,
			Version:    doc.Version,
			Content:    doc.Content,
		})
	}
	return docs
}
