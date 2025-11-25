package lsp

import (
	"encoding/json"
	"log"
)

// handleTextDocumentDidOpen handles the textDocument/didOpen notification
func (s *Server) handleTextDocumentDidOpen(request Request) {
	if s.GetState() != StateInitialized {
		log.Println("Received didOpen but server not initialized")
		return
	}

	// Parse params
	var params DidOpenTextDocumentParams
	paramsData, err := json.Marshal(request.Params)
	if err != nil {
		log.Printf("Failed to marshal params: %v", err)
		return
	}

	if err := json.Unmarshal(paramsData, &params); err != nil {
		log.Printf("Failed to parse didOpen params: %v", err)
		return
	}

	// Store the document
	doc := &Document{
		URI:        params.TextDocument.URI,
		LanguageID: params.TextDocument.LanguageID,
		Version:    params.TextDocument.Version,
		Content:    params.TextDocument.Text,
	}

	s.SetDocument(doc.URI, doc)

	log.Printf("Opened document: %s (v%d, %s) - %d bytes",
		doc.URI,
		doc.Version,
		doc.LanguageID,
		len(doc.Content))
}

// handleTextDocumentDidChange handles the textDocument/didChange notification
func (s *Server) handleTextDocumentDidChange(request Request) {
	if s.GetState() != StateInitialized {
		log.Println("Received didChange but server not initialized")
		return
	}

	// Parse params
	var params DidChangeTextDocumentParams
	paramsData, err := json.Marshal(request.Params)
	if err != nil {
		log.Printf("Failed to marshal params: %v", err)
		return
	}

	if err := json.Unmarshal(paramsData, &params); err != nil {
		log.Printf("Failed to parse didChange params: %v", err)
		return
	}

	// Get existing document
	doc, ok := s.GetDocument(params.TextDocument.URI)
	if !ok {
		log.Printf("Document not found: %s", params.TextDocument.URI)
		return
	}

	// Apply changes
	// Since we're using Full sync mode, we only need to handle full document updates
	for _, change := range params.ContentChanges {
		if change.Range == nil {
			// Full document update
			doc.Content = change.Text
			doc.Version = params.TextDocument.Version
			log.Printf("Updated document: %s (v%d) - %d bytes",
				doc.URI,
				doc.Version,
				len(doc.Content))
		} else {
			// Incremental update (not implemented in Full sync mode)
			log.Printf("Warning: Received incremental change but only Full sync is supported")
		}
	}

	// Update the document
	s.SetDocument(doc.URI, doc)
}

// handleTextDocumentDidClose handles the textDocument/didClose notification
func (s *Server) handleTextDocumentDidClose(request Request) {
	if s.GetState() != StateInitialized {
		log.Println("Received didClose but server not initialized")
		return
	}

	// Parse params
	var params DidCloseTextDocumentParams
	paramsData, err := json.Marshal(request.Params)
	if err != nil {
		log.Printf("Failed to marshal params: %v", err)
		return
	}

	if err := json.Unmarshal(paramsData, &params); err != nil {
		log.Printf("Failed to parse didClose params: %v", err)
		return
	}

	// Remove the document
	s.DeleteDocument(params.TextDocument.URI)

	log.Printf("Closed document: %s", params.TextDocument.URI)
}
