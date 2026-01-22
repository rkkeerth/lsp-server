package main

import (
	"sync"

	"go.lsp.dev/protocol"
)

// Document represents a text document with its content and version
type Document struct {
	URI     protocol.DocumentURI
	Content string
	Version int32
}

// DocumentStore manages the state of open documents
type DocumentStore struct {
	mu        sync.RWMutex
	documents map[protocol.DocumentURI]*Document
}

// NewDocumentStore creates a new document store
func NewDocumentStore() *DocumentStore {
	return &DocumentStore{
		documents: make(map[protocol.DocumentURI]*Document),
	}
}

// Set stores or updates a document in the store
func (ds *DocumentStore) Set(uri protocol.DocumentURI, content string, version int32) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.documents[uri] = &Document{
		URI:     uri,
		Content: content,
		Version: version,
	}
}

// Get retrieves a document from the store
func (ds *DocumentStore) Get(uri protocol.DocumentURI) (*Document, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	doc, exists := ds.documents[uri]
	return doc, exists
}

// Delete removes a document from the store
func (ds *DocumentStore) Delete(uri protocol.DocumentURI) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	delete(ds.documents, uri)
}

// List returns all URIs of documents in the store
func (ds *DocumentStore) List() []protocol.DocumentURI {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	uris := make([]protocol.DocumentURI, 0, len(ds.documents))
	for uri := range ds.documents {
		uris = append(uris, uri)
	}
	return uris
}

// Count returns the number of documents in the store
func (ds *DocumentStore) Count() int {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	return len(ds.documents)
}
