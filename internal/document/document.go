package document

import "sync"

// TextDocument represents an open text document
type TextDocument struct {
	URI        string
	LanguageID string
	Version    int
	Content    string
}

// Store manages open text documents
type Store struct {
	mu        sync.RWMutex
	documents map[string]*TextDocument
}

// NewStore creates a new document store
func NewStore() *Store {
	return &Store{
		documents: make(map[string]*TextDocument),
	}
}

// Open adds a document to the store
func (s *Store) Open(uri, languageID string, version int, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.documents[uri] = &TextDocument{
		URI:        uri,
		LanguageID: languageID,
		Version:    version,
		Content:    content,
	}
}

// Close removes a document from the store
func (s *Store) Close(uri string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.documents, uri)
}

// Get retrieves a document by URI
func (s *Store) Get(uri string) (*TextDocument, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, ok := s.documents[uri]
	return doc, ok
}

// Update applies changes to a document (full sync)
func (s *Store) Update(uri string, version int, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if doc, ok := s.documents[uri]; ok {
		doc.Version = version
		doc.Content = content
	}
}
