package main

import (
	"sync"
)

// Document represents a text document with its content and version
type Document struct {
	URI     string
	Content string
	Version int32
}

// DocumentManager manages the state of all open documents
type DocumentManager struct {
	mu        sync.RWMutex
	documents map[string]*Document
}

// NewDocumentManager creates a new document manager
func NewDocumentManager() *DocumentManager {
	return &DocumentManager{
		documents: make(map[string]*Document),
	}
}

// Open adds a new document to the manager
func (dm *DocumentManager) Open(uri, content string, version int32) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.documents[uri] = &Document{
		URI:     uri,
		Content: content,
		Version: version,
	}
}

// Update updates an existing document's content
func (dm *DocumentManager) Update(uri, content string, version int32) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if doc, exists := dm.documents[uri]; exists {
		doc.Content = content
		doc.Version = version
	}
}

// Close removes a document from the manager
func (dm *DocumentManager) Close(uri string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	delete(dm.documents, uri)
}

// Get retrieves a document from the manager
func (dm *DocumentManager) Get(uri string) (*Document, bool) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	doc, exists := dm.documents[uri]
	return doc, exists
}

// GetAll returns all documents
func (dm *DocumentManager) GetAll() []*Document {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	docs := make([]*Document, 0, len(dm.documents))
	for _, doc := range dm.documents {
		docs = append(docs, doc)
	}
	return docs
}

// Count returns the number of open documents
func (dm *DocumentManager) Count() int {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	return len(dm.documents)
}
