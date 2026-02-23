package server

// State holds the server's document state
type State struct {
	// Documents maps document URIs to their content
	Documents map[string]string
}

// NewState creates a new State instance
func NewState() *State {
	return &State{
		Documents: make(map[string]string),
	}
}

// OpenDocument adds a document to the store
func (s *State) OpenDocument(uri, content string) {
	s.Documents[uri] = content
}

// UpdateDocument updates a document's content (full sync)
func (s *State) UpdateDocument(uri, content string) {
	s.Documents[uri] = content
}

// GetDocument returns a document's content and whether it exists
func (s *State) GetDocument(uri string) (string, bool) {
	content, ok := s.Documents[uri]
	return content, ok
}

// CloseDocument removes a document from the store
func (s *State) CloseDocument(uri string) {
	delete(s.Documents, uri)
}
